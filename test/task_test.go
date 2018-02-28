package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nirry/gotask"
)

func Test_Task(t *testing.T) {

	total := 0

	thread := gotask.NewTaskManage(10)
	thread.Init = func(this *gotask.TaskManage) {

		if this.Ext == nil {
			this.Ext = 0
		}

		if this.Ext.(int) == 0 {

			fmt.Println("第一次执行,初始化数据")

			//初始化数据写处理
			for index := 0; index < 100; index++ {
				this.Datas = append(this.Datas, index)
			}

		} else {

			//用于测试
			fmt.Println("第二次执行,这里不初始化数据")
		}

	}

	thread.Pull = func(this *gotask.TaskManage) interface{} {

		if len(this.Datas) == 0 {
			return nil
		}

		data := this.Datas[len(this.Datas)-1]
		this.Datas = this.Datas[0 : len(this.Datas)-1]

		return data
	}

	thread.Process = func(this *gotask.TaskManage, data interface{}) {

		//TODO do somethings
		time.Sleep(time.Second)
		fmt.Println(data.(int), "处理完成")

		total++
	}

	thread.Completed = func(this *gotask.TaskManage) {

		ext := this.Ext.(int)
		ext++
		this.Ext = ext

		fmt.Println("所有任务执行完成", total)

	}

	thread.Start()

}
