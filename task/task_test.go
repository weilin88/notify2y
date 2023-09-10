package task

import (
	"fmt"
	"testing"
)

func Test_Common(t *testing.T) {
	s := new(TaskService)
	err := s.Init()
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	task := Task{Content: "you goal list.", Subject: "import message", Type: "IM"}
	s.AddTask(&task)
}
