package main

import (
	"fmt"
	"time"

	"github.com/weilin88/notify2y/cmd"
	"github.com/weilin88/notify2y/core"
	"github.com/weilin88/notify2y/one"
)

func setFuns(ct *cmd.Context) {
	ct.CmdMap = map[string]*cmd.Program{}

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
	pro = new(cmd.Program)
	pro.Name = "send"
	pro.Desc = "remove a file or dir to trash"
	pro.Usage = "usage: " + pro.Name + " [OPTION]  [file|dir]"
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		Name:      "h",
		LongName:  "help",
		NeedValue: false,
		Desc:      "print help"}

	pro.ParamDefMap["p"] = &cmd.ParamDef{
		Name:      "p",
		LongName:  "person",
		NeedValue: true,
		Desc:      ""}
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
		if ct.ParamGroupMap["p"] != nil {
			person = ct.ParamGroupMap["p"].Value
		}
		if person == "" {
			fmt.Printf("pls enter person\n")
			return
		}
		if pro.Target == "" {
			fmt.Printf("pls enter email content\n")
			return
		}
		err = cli.APISendMail(person, "test", pro.Target)
		if err != nil {
			fmt.Printf("err = %s\n", err.Error())
		} else {
			fmt.Printf("sended\n")
		}
	}

	//next add new user
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

	pro = new(cmd.Program)
	pro.Name = "web"
	pro.Desc = "run this http super serivce (beta version)"
	pro.Usage = "usage: " + pro.Name + " [OPTION]"
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		"h",
		"help",
		false,
		"print help"}
	pro.ParamDefMap["s"] = &cmd.ParamDef{
		"s",
		"https",
		false,
		"enable https service ,need cacert.pem ,privkey.pem on current dir"}
	pro.ParamDefMap["u"] = &cmd.ParamDef{
		"u",
		"url",
		true,
		"setup service address for this service,as -u :5555"}

	ct.CmdMap[pro.Name] = pro
	pro.Cmd = func(pro *cmd.Program) {

		if ct.ParamGroupMap["h"] != nil {
			cmd.PrintCmdHelp(pro)
			return
		}
		address := ":8080"
		upp := ct.ParamGroupMap["u"]
		if upp != nil {
			address = upp.Value
		}
		https := false
		if ct.ParamGroupMap["s"] != nil {
			https = true
		}
		Serivce(address, https)
	}
	pro = new(cmd.Program)
	pro.Name = "users"
	pro.Desc = "list login users"
	pro.Usage = "usage: " + pro.Name + " [OPTION]"
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		"h",
		"help",
		false,
		"print help"}
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
	pro = new(cmd.Program)
	pro.Name = "su"
	pro.Desc = "swich to other logined user"
	pro.Usage = "usage: " + pro.Name + " [OPTION]... [UserName]"
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		"h",
		"help",
		false,
		"print help"}

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
	pro = new(cmd.Program)
	pro.Name = "saveUser"
	pro.Desc = "save current user to name"
	pro.Usage = "usage: " + pro.Name + " [OPTION]... [UserName]"
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		"h",
		"help",
		false,
		"print help"}

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
	pro = new(cmd.Program)
	pro.Name = "who"
	pro.Desc = "show current user name"
	pro.Usage = "usage: " + pro.Name
	pro.ParamDefMap = map[string]*cmd.ParamDef{}

	pro.ParamDefMap["h"] = &cmd.ParamDef{
		"h",
		"help",
		false,
		"print help"}

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

}
func main() {
	core.Debug = true
	one.InitOneShowConfig()
	ct := cmd.NewContext()
	setFuns(ct)
	ct.Run()
}
