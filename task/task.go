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
	Subject    string `json:"subject"`
	Type       string `json:"type"`
	Content    string `json:"Content"`
	Importance string `json:"importance"`

	Deadline one.Timestamp `json:"sentDateTime"`

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
	return n
}

const data_json = "task_data.json"

type TaskService struct {
	taskData []*Task
	index    map[string]*Task
}

func (s *TaskService) AddTask(t *Task) error {
	if t.ID == "" {
		t.ID = uuid.NewV1().String()
	}
	if s.index[t.ID] != nil {
		return fmt.Errorf("exist task, id = %s", t.ID)
	}
	s.taskData = append(s.taskData, t)
	s.index[t.ID] = t
	s.saveAll()
	return nil
}

func (s *TaskService) UpdateTask(t *Task) error {
	if s.index[t.ID] != nil {
		old := s.index[t.ID]
		old.Content = t.Content
		old.Subject = t.Subject
		old.Deadline = t.Deadline
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
