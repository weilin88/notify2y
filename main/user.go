package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/weilin88/notify2y/one"
)

const config string = one.ConfigFileDefault + "."

type UserOP struct {
}

func (op *UserOP) ListUsers() ([]string, error) {
	home := one.GetConfigDir()
	return op.loopDir(home)
}
func (op *UserOP) SaveUser(user string) error {
	home := one.GetConfigDir()
	userDec := filepath.Join(home, config+user)
	userSrc := filepath.Join(home, one.ConfigFileDefault)
	return op.copyUser(userSrc, userDec)
}
func (op *UserOP) SwitchUser(user string) error {
	home := one.GetConfigDir()
	decFile := filepath.Join(home, one.CurUser)
	return ioutil.WriteFile(decFile, []byte(user), 0660)
}
func (op *UserOP) copyUser(userSrc string, userDec string) error {
	src, err := os.Open(userSrc)
	if err != nil {
		return err
	}
	defer src.Close()
	curFile, err := os.OpenFile(userDec, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}
	defer curFile.Close()
	_, err = io.Copy(curFile, src)
	return err
}
func (op *UserOP) Who() (string, error) {
	home := one.GetConfigDir()
	decFile := filepath.Join(home, one.CurUser)
	buff, err := ioutil.ReadFile(decFile)
	return string(buff), err
}
func (op *UserOP) loopDir(dirName string) ([]string, error) {
	li := []string{}
	fileList, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, err
	}
	for _, f := range fileList {
		info := f
		if f.IsDir() {
			continue
		}
		path := filepath.Join(dirName, info.Name())
		lname := strings.ToLower(info.Name())
		if strings.HasPrefix(lname, config) {
			li = append(li, path)
		}
	}
	return li, nil
}
