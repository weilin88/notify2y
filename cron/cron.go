package cron

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"github.com/weilin88/notify2y/task"
)

type CronManager struct {
	TaskService *task.TaskService
	core        *cron.Cron
	curCronTask map[string]*TaskAndCron
}

type TaskAndCron struct {
	TaskID      string
	TaskVersion int
	CronID      cron.EntryID
}

func (c *CronManager) Start() error {
	li, err := c.TaskService.ListTask()
	if err != nil {
		return err
	}
	core := cron.New()
	c.curCronTask = map[string]*TaskAndCron{}
	c.core = core

	//
	err = c.addCron(li)
	if err != nil {
		return err
	}

	core.Start()
	c.addTaskUpdateTask2Cron()
	return nil
}
func (c *CronManager) addCron(li []*task.Task) error {
	if li == nil {
		return nil
	}
	for _, t := range li {
		if t.Cron != "" {
			cid, err := c.core.AddFunc(t.Cron, func() {
				fmt.Printf("by %s rute \n", t.Cron)
				c.TaskService.SendMail(t)
			})
			if err != nil {
				fmt.Println("err = ", err)
			} else {
				fmt.Printf("%s on %s \n", t.Subject, t.Cron)
				tac := new(TaskAndCron)
				tac.CronID = cid
				tac.TaskID = t.ID
				tac.TaskVersion = t.Version
				c.curCronTask[t.ID] = tac
			}
		} else {
			fmt.Printf("empty cron value for  %s \n", t.Subject)
		}
	}
	return nil
}
func (c *CronManager) addTaskUpdateTask2Cron() {
	fmt.Println("executing addTaskUpdateTask2Cron")
	//id, err := c.core.AddFunc("*/5 * * * *", func() {
	id, err := c.core.AddFunc("*/5 * * * *", func() {
		err := c.Update()
		if err != nil {
			fmt.Println("execute addTaskUpdateTask2Cron err = ", err)
		}
	})
	if err != nil {
		fmt.Println("add addTaskUpdateTask2Cron task err = ", err)
	} else {
		fmt.Println("add addTaskUpdateTask2Cron id =  ", id)
	}
}

func (c *CronManager) Update() error {
	fmt.Println("task cron update...")
	err := c.TaskService.Init()
	if err != nil {
		return err
	}
	taskList, err := c.TaskService.ListTask()
	if err != nil {
		return err
	}
	addList := []*task.Task{}
	stillList := []*task.Task{}
	//removeList := []*task.Task{}

	for _, ct := range taskList {

		if c.curCronTask[ct.ID] != nil {
			oldTask := c.curCronTask[ct.ID]
			if oldTask.TaskVersion != ct.Version {
				//update task content.
				stillList = append(stillList, ct)
				// cancel old task
				c.core.Remove(oldTask.CronID)
				fmt.Printf("remove cron# cronid=%d taskid=%s taskVersion=%d \n", oldTask.CronID, oldTask.TaskID, oldTask.TaskVersion)
			} else {
				//nothing
				fmt.Printf("don't touch cronid=%d taskid=%s \n", oldTask.CronID, oldTask.TaskID)
			}
		} else {
			// new add
			addList = append(addList, ct)
		}
	}
	c.addCron(stillList)
	c.addCron(addList)
	return nil
}
