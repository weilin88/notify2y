package one

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/weilin88/notify2y/core"
)

//CurUser who is the current user
const CurUser string = ".od_cur_user.id"

//ConfigFileDefault default user when login
const ConfigFileDefault string = ".id.json"

const notify2y_app_config_file string = ".notify2y.json"
const app_config_dir = ".config/notify2y"

var NOTIFY2Y_CONFIG *Notify2YConfig

func GetConfigDir() string {
	home, _ := os.UserHomeDir()
	configDir := app_config_dir

	ret := filepath.Join(home, configDir)
	if core.ExistFile(ret) {
		return ret
	}
	os.MkdirAll(ret, os.ModePerm)
	return ret
}

func getCurUser() string {
	envUser := os.Getenv("oneuser")
	envUser = strings.TrimSpace(envUser)
	home := GetConfigDir()
	user := ""
	if envUser != "" {
		user = envUser
		fmt.Println("user envUser :", user)
	} else {
		buff, err := ioutil.ReadFile(filepath.Join(home, CurUser))
		if err != nil {
			user = ""
		} else {
			userName := string(buff)
			userName = strings.TrimSpace(userName)
			user = userName
		}
	}
	core.Println("using config = ", user)
	return user
}
func (u *OneClient) setUserInfo(name string) {
	u.UserName = name
	if name == "" {
		u.ConfigFile = ConfigFileDefault
	} else {
		u.ConfigFile = ConfigFileDefault + "." + name
	}
}
func (u *OneClient) findConfigFile() (string, error) {
	home := GetConfigDir()
	buff, err := ioutil.ReadFile(filepath.Join(home, u.ConfigFile))
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

//InitOneShowConfig load oneshow config information
func InitOneShowConfig() {
	//HOME USER PWD SHELL
	NOTIFY2Y_CONFIG = new(Notify2YConfig)
	home := GetConfigDir()
	if home != "" {
		fullPath := filepath.Join(home, notify2y_app_config_file)
		buff, err := ioutil.ReadFile(fullPath)
		if err != nil {
			return
		}
		err = json.Unmarshal(buff, NOTIFY2Y_CONFIG)
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
		//set application config
		setupNotify2yConfig()
	}
}

func setupNotify2yConfig() {
	cfg := NOTIFY2Y_CONFIG
	if cfg.Client_ID != "" && cfg.ClientSecret != "" {
		//fmt.Println("using a third-party client :", cfg.Client_ID)
		CLIENT_ID = cfg.Client_ID
		CLIENT_SECRET = cfg.ClientSecret
		if cfg.Scope != "" {
			SCOPE = cfg.Scope
		}
		if cfg.RedirectURL != "" {
			CALLBACK_URL = cfg.RedirectURL
		}
	}
}
func (u *OneClient) getConfigAuthToken() *AuthToken {
	//HOME USER PWD SHELL
	cfg := new(AuthToken)
	content, err := u.findConfigFile()
	//fmt.Println(content)
	if content == "" {
		fmt.Println("can not find config file")
		return nil
	}
	err = json.Unmarshal([]byte(content), cfg)
	if err != nil {
		fmt.Println("err = ", err)
		return nil
	}
	return cfg
}

//SaveToken2Home home
func (u *OneClient) SaveToken2Home(token *AuthToken) error {
	home := GetConfigDir()
	pcfg := ""
	if home != "" {
		pcfg = filepath.Join(home, u.ConfigFile)
	} else {
		return errors.New("can not found home dir")
	}
	return SaveToken2Config(token, pcfg)
}

//SaveToken2DefaultPath save config when first login
func SaveToken2DefaultPath(token *AuthToken) error {
	home := GetConfigDir()
	pcfg := ""
	if home != "" {
		pcfg = filepath.Join(home, ConfigFileDefault)
	} else {
		return errors.New("can not found home dir")
	}
	return SaveToken2Config(token, pcfg)
}

//SaveToken2Config save to configure file
func SaveToken2Config(token *AuthToken, configFile string) error {
	buff, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, buff, 0660)
}
