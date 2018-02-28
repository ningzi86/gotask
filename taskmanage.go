package gotask

import (
	"sync"
	"time"
)

// TaskManage 任务管理类
type TaskManage struct {
	locker     *sync.Mutex
	capnum     uint
	mc         chan uint
	Datas      []interface{}
	Processing func(data interface{})

	isRunning bool
	Ext       interface{}
	Pull      func(manage *TaskManage) interface{}
	Init      func(manage *TaskManage)
	Process   func(manage *TaskManage, data interface{})

	Completed func(manage *TaskManage)
}

// NewTaskManage 初始化
func NewTaskManage(num uint) *TaskManage {
	mc := make(chan uint, num)
	return &TaskManage{capnum: num, mc: mc, locker: new(sync.Mutex)}
}

// getOne 获取一个协程
func (this *TaskManage) getOne() {
	this.mc <- 1
}

// freeOne 释放一个协程
func (this *TaskManage) freeOne() {
	<-this.mc
}

// has 待处理协程数
func (this *TaskManage) has() uint {
	return uint(len(this.mc))
}

// left 剩余数
func (this *TaskManage) left() uint {
	return this.capnum - uint(len(this.mc))
}

// Start 开始执行多协程操作
func (this *TaskManage) Start() {

	if this.Init == nil {
		panic("未设置初始化方法")
	}

	//置为启动
	if this.isRunning {
		return
	}

	//初始化
	this.Init(this)
	this.isRunning = true

	for {

		//没有待处理数据
		if len(this.Datas) == 0 {
			break
		}

		for {

			//获取需要处理的数据
			this.locker.Lock()
			data := this.Pull(this)
			this.locker.Unlock()

			if data == nil {
				if this.has() == 0 {
					break
				} else {
					time.Sleep(500)
					continue
				}
			}

			this.getOne()
			go func(data interface{}) {

				defer this.freeOne()

				//处理数据
				this.Process(this, data)

			}(data)
		}

	}

	//置为停止
	this.isRunning = false

	//完成时方法
	if this.Completed != nil {
		this.Completed(this)
	}

}
