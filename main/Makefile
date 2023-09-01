buildTime=$(shell date +"%F %T")
out=notify2y
objs=notify.go web.go user.go
ver=v1.0.0
$(out) : $(objs)
	echo $(buildTime)
	go build  -ldflags "-s -w -X 'github.com/weilin88/notify2y/cmd.VERSION=$(ver) build time:$(buildTime)' -X 'github.com/weilin88/notify2y/cmd.CmdName=notify2y' "  -o $(out)
.PHONY : win
.PHONY : static
clean :
	rm -rf $(out)
	rm -rf $(out).exe
install :
	upx $(out)
win :
	GOOS=windows GOARCH=386 go build  -ldflags "-s -w -X 'github.com/weilin88/notify2y/cmd.VERSION=$(ver) build time:$(buildTime)' -X 'github.com/weilin88/oneshow/cmd.CmdName=notify2y' "  -o $(out).exe
arm64 :
	GOOS=linux GOARCH=arm64 go build  -ldflags "-s -w -X 'github.com/weilin88/notify2y/cmd.VERSION=$(ver) build time:$(buildTime)' -X 'github.com/weilin88/oneshow/cmd.CmdName=notify2y' "  -o $(out)-arm
static :
	CGO_ENABLED=0 go build -ldflags "-s -w" -o $(out)
.PHONY : pkg
pkg : $(out) win
	zip -r -0 $(out)-linux-amd64.zip $(out)
	zip -r -0 $(out)-windows-386.zip $(out).exe
