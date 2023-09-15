package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	"github.com/weilin88/notify2y/one"
	"github.com/weilin88/notify2y/task"
)

var tools = new(Tools)

type WebContext struct {
	Address     string
	EnableTLS   bool
	EnableEmbed bool
	Cli         *one.OneClient
	TaskAPI     *task.TaskService
}

type JsonRequest struct {
	Argvs []string
}
type JsonResponse struct {
	Error   bool
	Message string
	Data    interface{}
}

const sys_err = `{
	"Error":true
	"Message":"system error"
	"Data" : null
}
`

func parseRequst(w http.ResponseWriter, r *http.Request) (*JsonRequest, error) {
	defer r.Body.Close()
	buff := &bytes.Buffer{}
	_, err := io.Copy(buff, r.Body)
	if err != nil {
		return nil, err
	}
	jsonObj := new(JsonRequest)
	err = json.Unmarshal(buff.Bytes(), jsonObj)
	if err != nil {
		return nil, err
	}
	return jsonObj, nil
}

func responseResult(w http.ResponseWriter, r *http.Request, result *JsonResponse) {
	w.Header().Add("Content-Type", "application/json")
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		fmt.Println("err = ", err)
		w.Write([]byte(sys_err))
		return
	}
	w.Write(jsonBytes)
}

func checkParamVaild(req *JsonRequest, argsLen int, result *JsonResponse) bool {
	if req != nil && len(req.Argvs) == argsLen {
		return true
	}
	result.Error = false
	result.Message = "invalid argvs"
	return false
}

func AutoUpdateToken(cli *one.OneClient) {
	for {
		CheckToken(cli)
		time.Sleep(time.Minute)
	}
}
func CheckToken(cli *one.OneClient) error {
	expires := time.Time(cli.Token.ExpiresTime)

	expires = expires.Truncate(time.Minute)
	if time.Now().After(expires) {
		//fmt.Println("to expries time, update token")
		newToken, err := cli.UpdateToken()
		if err != nil {
			return err
		}
		cli.Token = newToken
	}
	return nil

}

func GetQueryParamByKey(r *http.Request, key string) string {
	keys, ok := r.URL.Query()[key]
	if !ok || len(keys[0]) < 1 {
		return ""
	}
	return keys[0]
}

func Serivce(ctx *WebContext) {
	cli, err := one.NewOneClient()
	if err != nil {
		panic(err.Error())
	}
	s := new(task.TaskService)
	err = s.Init()
	if err != nil {
		fmt.Println("err = ", err)
		panic(err.Error())
	}
	var httpRoot http.FileSystem
	if ctx.EnableEmbed {
		httpRoot = http.FS(staticSource)
		http.Handle("/static/", http.FileServer(httpRoot))
	} else {
		sourceDir := "./static"
		httpRoot = http.Dir(sourceDir)
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(httpRoot)))
	}
	ctx.Cli = cli
	ctx.TaskAPI = s

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"error\":\"cannot access /\"}"))
	})

	http.HandleFunc("/html/", func(w http.ResponseWriter, r *http.Request) {
		notFoundHtml := `<html>
			<head>
				<title>Not Found</title>
			</head>
			<body>
			404 Not Found
			</body>
		</html>`
		sourceHtml := filepath.Base(r.URL.Path)
		w.Header().Add("Content-Type", "text/html; charset=utf-8")

		buff, err := tools.GetFileContent(ctx, httpRoot, sourceHtml)
		if err != nil {
			fmt.Printf("load source file to fail,file = %s,err = %s\n", sourceHtml, err.Error())
			w.Write([]byte(notFoundHtml))
		} else {
			if sourceHtml == "task-create.html" {
				taskCreateHtml(w, r, ctx, buff)
			} else if sourceHtml == "task-detail.html" {
				//billCreateHtml(w, r, context, buff)
			} else {
				w.Write([]byte(notFoundHtml))
			}
		}
	})

	api := new(API)
	api.Context = ctx
	http.HandleFunc("/call", func(w http.ResponseWriter, r *http.Request) {
		result := new(JsonResponse)
		method := GetQueryParamByKey(r, "method")
		fmt.Println("call method:", method)
		req, err := parseRequst(w, r)
		if err != nil {
			result.Error = true
			result.Message = err.Error()
			responseResult(w, r, result)
			return
		}
		switch method {
		case "createTask":
			api.CreateTask(req, result)
		case "delTask":
			api.DelTask(req, result)
		case "sales":
			api.SearchTask(req, result)
		case "notify2you":
			api.TaskNotify2You(req, result)
		default:
			result.Error = true
			result.Message = "can not find called method :" + method
		}
		responseResult(w, r, result)
	})

	if ctx.EnableTLS {
		fmt.Println("https server on ", ctx.Address)
		err = http.ListenAndServeTLS(ctx.Address, "cer.pem", "key.pem", nil)
	} else {
		fmt.Println("http server on ", ctx.Address)
		err = http.ListenAndServe(ctx.Address, nil)
	}
	if err != nil {
		fmt.Println("run thie service to failed on error = ", err)
	}
}

func taskCreateHtml(w http.ResponseWriter, r *http.Request, context *WebContext, buff []byte) {
	ID := GetQueryParamByKey(r, "id")
	var ot *task.Task
	var err error
	if ID == "" {
		ot, err = nil, nil
	} else {
		ot, err = context.TaskAPI.GetTask(ID)
	}
	result := new(JsonResponse)
	if err != nil {
		result.Error = true
		result.Message = err.Error()
	}
	result.Data = ot

	html := tools.RenderHtmlObject(buff, result)
	w.Write(html)
}
