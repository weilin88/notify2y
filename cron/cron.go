package cron

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/weilin88/notify2y/task"
)

type CronManager struct {
	TaskService *task.TaskService
}

func (c *CronManager) Start() error {
	li, err := c.TaskService.ListTask()
	if err != nil {
		return err
	}
	ct := cron.New()
	for _, t := range li {
		if t.Cron != "" {
			_, err = ct.AddFunc(t.Cron, func() {
				c.TaskService.SendMail(t)
			})
			if err != nil {
				fmt.Printf("\n")
			} else {
				fmt.Printf("%s on %s \n", t.Subject, t.Cron)
			}
		} else {
			fmt.Printf("empty cron value for  %s \n", t.Subject)
		}
	}
	ct.Start()
	return nil
}
