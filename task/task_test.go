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
	task := Task{Content: "test task content", Subject: "pls"}
	s.AddTask(&task)
}
