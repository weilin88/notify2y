package task

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
	"github.com/weilin88/notify2y/core"
	"github.com/weilin88/notify2y/one"
)

type Task struct {
	ID         string `json:"id"`
	Version    int    `json:"version"`
	Subject    string `json:"subject"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	Importance string `json:"importance"`

	Deadline one.Timestamp `json:"sentDateTime"`
	Cron     string        `json:"cron"`

	CreatedDateTime      one.Timestamp `json:"createdDateTime"`
	LastModifiedDateTime one.Timestamp `json:"lastModifiedDateTime"`
}

func (t *Task) Copy() *Task {
	n := new(Task)
	n.Content = t.Content
	n.CreatedDateTime = t.CreatedDateTime
	n.Deadline = t.Deadline
	n.ID = t.ID
	n.Importance = t.Importance
	n.Subject = t.Subject
	n.Type = t.Type
	n.Cron = t.Cron
	n.Version = t.Version
	return n
}

const data_json = "task_data.json"

type TaskService struct {
	NotifySender *one.OneClient
	Person       string
	taskData     []*Task
	index        map[string]*Task
}

func (s *TaskService) AddTask(t *Task) error {
	if t.ID == "" {
		t.ID = uuid.NewV1().String()
	}
	if s.index[t.ID] != nil {
		return fmt.Errorf("exist task, id = %s", t.ID)
	}
	t.Version = 1
	s.taskData = append(s.taskData, t)
	s.index[t.ID] = t
	s.saveAll()
	return nil
}
func (s *TaskService) GetTask(ID string) (*Task, error) {
	task := s.index[ID]
	if task != nil {
		return task, nil
	}
	return nil, fmt.Errorf("Not Found Task ,ID = %s", ID)
}

func (s *TaskService) UpdateTask(t *Task) error {
	if s.index[t.ID] != nil {
		old := s.index[t.ID]
		old.Content = t.Content
		old.Subject = t.Subject
		old.Deadline = t.Deadline
		old.Cron = t.Cron
		old.Version++
		//return to UI
		t.Version = old.Version
		s.saveAll()
		return nil
	}
	return nil
}
func (s *TaskService) ListTask() ([]*Task, error) {
	nList := []*Task{}
	for _, v := range s.taskData {
		nList = append(nList, v.Copy())
	}
	return nList, nil
}
func (s *TaskService) DelTask(ID string) error {
	if s.index[ID] == nil {
		return nil
	}
	delete(s.index, ID)
	for idx, v := range s.taskData {
		if v.ID == ID {
			newArr := append(s.taskData[:idx], s.taskData[idx+1:]...)
			s.taskData = newArr
			return nil
		}
	}
	s.saveAll()
	return nil
}
func (s *TaskService) saveAll() error {
	dir := one.GetConfigDir()
	dataFile := filepath.Join(dir, data_json)
	buff, err := json.Marshal(s.taskData)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dataFile, buff, os.ModePerm)
}

func (s *TaskService) Init() error {
	dir := one.GetConfigDir()
	dataFile := filepath.Join(dir, data_json)
	if core.ExistFile(dataFile) {
		buff, err := ioutil.ReadFile(dataFile)
		if err != nil {
			return err
		}
		list := []*Task{}
		err = json.Unmarshal(buff, &list)
		s.taskData = list
		s.index = map[string]*Task{}
		for _, v := range list {
			s.index[v.ID] = v
		}
	} else {
		s.taskData = []*Task{}
		s.index = map[string]*Task{}
	}
	return nil
}

func (s *TaskService) Notify2You(cli *one.OneClient, person string) {
	li, err := s.ListTask()
	for _, v := range li {
		if v.Type == "IM" {
			err = cli.APISendMail(person, v.Subject, v.Content, "text")
			if err != nil {
				fmt.Printf("err = %s\n", err.Error())
			} else {
				fmt.Printf("sended\n")
			}
		}
	}
}
func (s *TaskService) SendMail(v *Task) {
	err := s.NotifySender.ExpresCheck()
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	err = s.NotifySender.APISendMail(s.Person, v.Subject, v.Content, "text")
	if err != nil {
		fmt.Printf("err = %s\n", err.Error())
	} else {
		fmt.Printf("sended\n")
	}
}

func (s *TaskService) TaskNotify2You(cli *one.OneClient, person string, id string) error {
	v := s.index[id]
	if person == "" {
		return fmt.Errorf("person cannot be empty")
	}
	if v != nil && v.Type == "IM" {
		err := cli.APISendMail(person, v.Subject, v.Content, "text")
		if err != nil {
			fmt.Printf("err = %s\n", err.Error())
		} else {
			fmt.Printf("sended\n")
		}
		return err
	}
	return fmt.Errorf("not found fit task")
}
