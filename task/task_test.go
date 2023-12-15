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
	l, err := s.ListTask()
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	for _, v := range l {
		fmt.Printf("subject=%s ; ID=%s \n", v.Subject, v.ID)
	}
}
