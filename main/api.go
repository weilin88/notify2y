package main

import (
	"encoding/json"
	"fmt"

	"github.com/weilin88/notify2y/task"
)

type API struct {
	Context *WebContext
}

func (api *API) DelTask(req *JsonRequest, result *JsonResponse) {
	fmt.Println(req.Argvs)
	context := api.Context
	if checkParamVaild(req, 1, result) {
		taskData := req.Argvs[0]
		fmt.Println("data = ", taskData)
		ot := new(task.Task)
		err := json.Unmarshal([]byte(taskData), ot)
		if err != nil {
			result.Error = true
			result.Message = err.Error()
		} else {
			err := context.TaskAPI.DelTask(ot.ID)
			if err != nil {
				result.Error = true
				result.Message = err.Error()
			}
		}
	}
}
func (api *API) CreateTask(req *JsonRequest, result *JsonResponse) {
	fmt.Println(req.Argvs)
	context := api.Context
	if checkParamVaild(req, 1, result) {
		orderData := req.Argvs[0]
		fmt.Println("data = ", orderData)
		ot := new(task.Task)
		err := json.Unmarshal([]byte(orderData), ot)
		if err != nil {
			result.Error = true
			result.Message = err.Error()
		} else {
			if err != nil {
				result.Error = true
				result.Message = err.Error()
				return
			}
			if ot.Version < 1 {
				err = context.TaskAPI.AddTask(ot)
			} else {
				err = context.TaskAPI.UpdateTask(ot)
			}
			//set order object
			result.Data = ot
			if err != nil {
				result.Error = true
				result.Message = err.Error()
			}
		}
	}
}

func (api *API) SearchTask(req *JsonRequest, result *JsonResponse) {
	fmt.Println(req.Argvs)
	context := api.Context
	if checkParamVaild(req, 0, result) {
		//consumer := req.Argvs[0]
		//dates := req.Argvs[1]
		li := []*task.Task{}
		var err error
		li, err = context.TaskAPI.ListTask()
		if err != nil {
			result.Error = true
			result.Message = err.Error()
		} else {
			result.Data = li
		}
	}
}

func (api *API) TaskNotify2You(req *JsonRequest, result *JsonResponse) {
	fmt.Println(req.Argvs)
	context := api.Context
	if checkParamVaild(req, 1, result) {
		id := req.Argvs[0]
		var err error
		err = context.TaskAPI.TaskNotify2You(context.Cli, APP_CONFIG.Email, id)
		if err != nil {
			result.Error = true
			result.Message = err.Error()
		}
	}
}
