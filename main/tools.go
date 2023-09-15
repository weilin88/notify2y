package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Tools struct{}

func (t *Tools) GetFileContent(context *WebContext, root http.FileSystem, fileName string) ([]byte, error) {
	//fix not found bug
	if context.EnableEmbed {
		fileName = "/static/" + fileName
	}
	f, err := root.Open(fileName)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}
func (t *Tools) RenderHtmlObject(html []byte, obj interface{}) []byte {
	buff := bytes.Buffer{}
	buff.Grow(1024 * 5)
	jsonObj, _ := json.Marshal(obj)
	key := "////PAGE_DATA="
	varObj := "PAGE_DATA="
	buff.WriteString(varObj)
	buff.Write(jsonObj)
	buff.WriteString(";")
	return bytes.Replace(html, []byte(key), buff.Bytes(), 1)
}
func (t *Tools) ParseLocalTime(timeStr string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
}
func (t *Tools) ParseExcelLocalTime2(timeStr string) (time.Time, error) {
	return time.ParseInLocation("01/02/2006 15:04:05", timeStr, time.Local)
}
func (t *Tools) ParseExcelLocalTime(timeStr string) (time.Time, error) {
	return time.ParseInLocation("01-02-2006 15:04:05", timeStr, time.Local)
}

func (t *Tools) ParseCSTTime(timeStr string) (time.Time, error) {
	d, _ := time.LoadLocation("Asia/Shanghai")
	return time.ParseInLocation("2006-01-02 15:04:05", timeStr, d)
}

func (tools *Tools) GetUTCStr(t time.Time) string {
	return t.Format(time.RFC3339)
}
