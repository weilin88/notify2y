package main

import (
	"embed"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "embed"

	"github.com/weilin88/notify2y/cmd"
	"github.com/weilin88/notify2y/core"
	ncron "github.com/weilin88/notify2y/cron"
	"github.com/weilin88/notify2y/notify"
	"github.com/weilin88/notify2y/one"
	"github.com/weilin88/notify2y/task"
)

//go:embed static/*
var staticSource embed.FS

func setFuns(ct *cmd.Context) {
	ct.CmdMap = map[string]*cmd.Program{}

	//#list
	pro := new(cmd.Program)
	pro.Name = "list"
	pro.Desc = "list email"
	pro.Usage = "usage: " + pro.Name + " [OPTION] path"
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		Name:      "h",
		LongName:  "help",
		NeedValue: false,
		Desc:      "print help"}

	pro.ParamDefMap["l"] = &cmd.ParamDef{
		Name:      "l",
		LongName:  "list",
		NeedValue: false,
		Desc:      "list files detail"}
	pro.ParamDefMap["d"] = &cmd.ParamDef{
		Name:      "d",
		LongName:  "direct_url",
		NeedValue: false,
		Desc:      "list files direct url"}

	ct.CmdMap[pro.Name] = pro
	pro.Cmd = func(pro *cmd.Program) {
		if ct.ParamGroupMap["h"] != nil {
			cmd.PrintCmdHelp(pro)
			return
		}
		dirPath := pro.Target
		if dirPath == "" {
			dirPath = "/"
		}
		strLen := len(dirPath)
		if strLen > 1 && dirPath[strLen-1] == '/' {
			dirPath = dirPath[:strLen-1]
		}
		cli, err := one.NewOneClient()
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
		ret, err := cli.APIListFilesByPath(cli.CurDriveID, dirPath)
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
		if ct.ParamGroupMap["l"] != nil {
			for _, v := range ret.Value {
				//name size owner
				mdTime := time.Time(v.LastModifiedDateTime)
				dsTime := mdTime.Local().Format(time.RFC3339)
				Name := v.Name
				if v.Folder != nil {
					Name = v.Name + "/"
				}
				fmt.Printf("%-10s%-16s%-28s%-100s\n", one.ViewHumanShow(v.Size), v.CreatedBy.User.DisplayName, dsTime, Name)
			}
		} else if ct.ParamGroupMap["d"] != nil {
			for _, v := range ret.Value {
				fmt.Printf("[%s]%-200s\n\n", v.Name, v.DownloadURL)
			}

		} else {
			for _, v := range ret.Value {
				Name := v.Name
				if v.Folder != nil {
					Name = v.Name + "/"
				}
				fmt.Printf("%s\n", Name)
			}

		}
	}

	//next remove command
	//#send
	pro = new(cmd.Program)
	pro.Name = "send"
	pro.Desc = "send mail"
	pro.Usage = "usage: " + pro.Name + " [OPTION]  [content]"
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		Name:      "h",
		LongName:  "help",
		NeedValue: false,
		Desc:      "print help"}

	pro.ParamDefMap["r"] = &cmd.ParamDef{
		Name:      "r",
		LongName:  "recipient",
		NeedValue: true,
		Desc:      "recipient email"}
	ct.CmdMap[pro.Name] = pro
	pro.Cmd = func(pro *cmd.Program) {
		if ct.ParamGroupMap["h"] != nil {
			cmd.PrintCmdHelp(pro)
			return
		}
		cli, err := one.NewOneClient()
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
		person := ""
		if ct.ParamGroupMap["r"] != nil {
			person = ct.ParamGroupMap["r"].Value
		}
		if person == "" {
			fmt.Printf("pls enter recipient email\n")
			return
		}
		if pro.Target == "" {
			fmt.Printf("pls enter email content\n")
			return
		}
		err = cli.APISendMail(person, "api test on notify2y", pro.Target, "text")
		if err != nil {
			fmt.Printf("err = %s\n", err.Error())
		} else {
			fmt.Printf("sended\n")
		}
	}
	//#notify
	pro = new(cmd.Program)
	pro.Name = "notify"
	pro.Desc = "notify to you by cron"
	pro.Usage = "usage: " + pro.Name + " [OPTION]"
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		Name:      "h",
		LongName:  "help",
		NeedValue: false,
		Desc:      "print help"}

	pro.ParamDefMap["r"] = &cmd.ParamDef{
		Name:      "r",
		LongName:  "recipient",
		NeedValue: true,
		Desc:      "recipient email"}
	pro.ParamDefMap["t"] = &cmd.ParamDef{
		Name:      "t",
		LongName:  "test",
		NeedValue: false,
		Desc:      "only for test,send notify without waiting"}
	ct.CmdMap[pro.Name] = pro
	pro.Cmd = func(pro *cmd.Program) {
		if ct.ParamGroupMap["h"] != nil {
			cmd.PrintCmdHelp(pro)
			return
		}
		person := ""
		if ct.ParamGroupMap["r"] != nil {
			person = ct.ParamGroupMap["r"].Value
		}
		if person == "" {
			fmt.Printf("person can not be empty.")
			return
		}
		if ct.ParamGroupMap["t"] != nil {
			backend(person)
		} else {
			cron(person)
		}
	}

	//next add new user
	//#auth
	pro = new(cmd.Program)
	pro.Name = "auth"
	pro.Desc = "get a auth for new user"
	pro.Usage = "usage: " + pro.Name + " [OPTION]"
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		Name:      "h",
		LongName:  "help",
		NeedValue: false,
		Desc:      "print help"}

	ct.CmdMap[pro.Name] = pro
	pro.Cmd = func(pro *cmd.Program) {
		if ct.ParamGroupMap["h"] != nil {
			cmd.PrintCmdHelp(pro)
			return
		}
		cli := one.NewBaseOneClient()
		cli.DoAutoForNewUser()
	}

	//#web
	pro = new(cmd.Program)
	pro.Name = "web"
	pro.Desc = "run this http super serivce (beta version)"
	pro.Usage = "usage: " + pro.Name + " [OPTION]"
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		Name:      "h",
		LongName:  "help",
		NeedValue: false,
		Desc:      "print help"}
	pro.ParamDefMap["s"] = &cmd.ParamDef{
		Name:      "s",
		LongName:  "https",
		NeedValue: false,
		Desc:      "enable https service ,need pri.key ,cert.key on current dir"}
	pro.ParamDefMap["e"] = &cmd.ParamDef{
		Name:      "e",
		LongName:  "embed",
		NeedValue: false,
		Desc:      "enable embed static source,default using static dir."}
	pro.ParamDefMap["u"] = &cmd.ParamDef{
		Name:      "u",
		LongName:  "url",
		NeedValue: true,
		Desc:      "setup service address for this service,as -u :5555"}

	ct.CmdMap[pro.Name] = pro
	pro.Cmd = func(pro *cmd.Program) {

		if ct.ParamGroupMap["h"] != nil {
			cmd.PrintCmdHelp(pro)
			return
		}
		ctx := new(WebContext)
		ctx.Address = ":8080"
		upp := ct.ParamGroupMap["u"]
		if upp != nil {
			ctx.Address = upp.Value
		}
		if ct.ParamGroupMap["s"] != nil {
			ctx.EnableTLS = true
		} else {
			ctx.EnableTLS = false
		}

		if ct.ParamGroupMap["e"] != nil {
			ctx.EnableEmbed = true
		} else {
			ctx.EnableEmbed = false
		}
		Serivce(ctx)
	}

	//#users
	pro = new(cmd.Program)
	pro.Name = "users"
	pro.Desc = "list login users"
	pro.Usage = "usage: " + pro.Name + " [OPTION]"
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		Name:      "h",
		LongName:  "help",
		NeedValue: false,
		Desc:      "print help"}
	ct.CmdMap[pro.Name] = pro
	pro.Cmd = func(pro *cmd.Program) {
		if ct.ParamGroupMap["h"] != nil {
			cmd.PrintCmdHelp(pro)
			return
		}
		op := new(UserOP)
		li, err := op.ListUsers()
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
		if len(li) == 0 {
			fmt.Println("pls call saveUser command for save a session")
			return
		}
		for _, user := range li {
			fmt.Println(user)
		}
	}
	//swich to other session
	//#su
	pro = new(cmd.Program)
	pro.Name = "su"
	pro.Desc = "swich to other logined user"
	pro.Usage = "usage: " + pro.Name + " [OPTION]... [UserName]"
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		Name:      "h",
		LongName:  "help",
		NeedValue: false,
		Desc:      "print help"}

	ct.CmdMap[pro.Name] = pro
	pro.Cmd = func(pro *cmd.Program) {
		if ct.ParamGroupMap["h"] != nil {
			cmd.PrintCmdHelp(pro)
			return
		}
		user := pro.Target
		if user == "" {
			fmt.Println("user name cannot be empty")
			return
		}
		op := new(UserOP)
		err := op.SwitchUser(user)
		if err != nil {
			fmt.Println("err = ", err)
		} else {
			fmt.Println("switch to ", user)
		}
	}

	//next program
	//#saveUser
	pro = new(cmd.Program)
	pro.Name = "saveUser"
	pro.Desc = "save current user to name"
	pro.Usage = "usage: " + pro.Name + " [OPTION]... [UserName]"
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		Name:      "h",
		LongName:  "help",
		NeedValue: false,
		Desc:      "print help"}

	ct.CmdMap[pro.Name] = pro
	pro.Cmd = func(pro *cmd.Program) {
		if ct.ParamGroupMap["h"] != nil {
			cmd.PrintCmdHelp(pro)
			return
		}
		user := pro.Target
		if user == "" {
			fmt.Println("user name cannot be empty")
			return
		}

		op := new(UserOP)
		err := op.SaveUser(user)
		if err != nil {
			fmt.Println("err = ", err)
		} else {
			fmt.Println("save to ", user)
		}
	}
	//next program
	//#who
	pro = new(cmd.Program)
	pro.Name = "who"
	pro.Desc = "show current user name"
	pro.Usage = "usage: " + pro.Name
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		Name:      "h",
		LongName:  "help",
		NeedValue: false,
		Desc:      "print help"}

	ct.CmdMap[pro.Name] = pro
	pro.Cmd = func(pro *cmd.Program) {
		if ct.ParamGroupMap["h"] != nil {
			cmd.PrintCmdHelp(pro)
			return
		}
		op := new(UserOP)
		userName, err := op.Who()
		if err != nil {
			fmt.Println("who command call failed, err = ", err)
		} else {
			fmt.Println("current user:", userName)
		}
	}
	//#ntest
	pro = new(cmd.Program)
	pro.Name = "ntest"
	pro.Desc = "ntest"
	pro.Usage = "usage: " + pro.Name + " [OPTION]  [content]"
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		Name:      "h",
		LongName:  "help",
		NeedValue: false,
		Desc:      "print help"}

	pro.ParamDefMap["r"] = &cmd.ParamDef{
		Name:      "r",
		LongName:  "recipient",
		NeedValue: true,
		Desc:      "recipient email"}
	ct.CmdMap[pro.Name] = pro
	pro.Cmd = func(pro *cmd.Program) {
		if ct.ParamGroupMap["h"] != nil {
			cmd.PrintCmdHelp(pro)
			return
		}
		person := ""
		if ct.ParamGroupMap["r"] != nil {
			person = ct.ParamGroupMap["r"].Value
		}
		if person == "" {
			fmt.Printf("pls enter recipient email\n")
			return
		}
		engine := notify.NewEngine()
		ll := notify.NewURLListener()
		engine.RegisterListener(ll)
		engine.StartAll()
		select {}
	}

}

func waitTime(person string) {
	now := time.Now()

	rt := 13
	runStr := os.Getenv("runTime")
	i, err := strconv.ParseInt(runStr, 10, 64)
	if err == nil {
		rt = int(i)
	}
	fmt.Printf("run on %d hour \n", rt)
	nextTen := time.Date(now.Year(), now.Month(), now.Day(), rt, 20, 0, 0, now.Location())
	if now.After(nextTen) {
		nextTen = nextTen.Add(24 * time.Hour)
	}

	waitTime := nextTen.Sub(now)

	timer := time.NewTimer(waitTime)
	<-timer.C

	//fmt.Println("It's 10:00 AM. Running your function now.")
	backend(person)
}
func backend(person string) {
	taskS := task.NewTaskService()
	cli, err := one.NewOneClient()
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	taskS.Notify2You(cli, person)
}
func cron(person string) {
	taskS := task.NewTaskService()
	cli, err := one.NewOneClient()
	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	taskS.NotifySender = cli
	taskS.Person = person

	crr := new(ncron.CronManager)
	crr.TaskService = taskS
	crr.Start()
	for {
		fmt.Println("next 1 min")
		time.Sleep(time.Minute)
	}
}
func main() {
	core.Debug = true
	one.InitOneShowConfig()
	ct := cmd.NewContext()
	setFuns(ct)
	ct.Run()
}
