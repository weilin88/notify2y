package one

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	chttp "github.com/weilin88/notify2y/http"

	"github.com/weilin88/notify2y/core"
)

var CLIENT_ID string = "a30ee80e-529e-4da2-9bd0-91a48c8f0e2a"

var CLIENT_SECRET string = "K2c8Q~GkX43l-E_wYXI2tybUombi3MnL5iQOocjn"

var SCOPE string = "Mail.ReadWrite Mail.Send offline_access"

var CALLBACK_URL = "http://localhost:1234/checkLogin"

//OneClient is context object
type OneClient struct {
	HTTPClient *chttp.HttpClient

	SSOHost string
	APIHost string

	UserName   string
	ConfigFile string
	Token      *AuthToken
	CurDriveID string
}

//NewDefaultCli new a default oneshow client
func NewDefaultCli() (*OneClient, error) {
	return NewOneClient()
}

// GetAuthCodeURL gen an URL for get auto code from API
func (cli *OneClient) GetAuthCodeURL() string {
	parms := cli.GetOneDriveAppInfo()
	parms["access_type"] = "offline"
	parms["response_type"] = "code"
	parms["state"] = "uid-notify2y"

	uri := "/common/oauth2/v2.0/authorize"
	parmsVal := url.Values{}
	for k, v := range parms {
		parmsVal.Add(k, v)
	}
	URL := cli.SSOHost + uri + "?" + parmsVal.Encode()
	core.Println(URL)
	return URL
}

//GetFirstToken get token and reflesh token by code
func (cli *OneClient) GetFirstToken(code string) error {
	parms := cli.GetOneDriveAppInfoWithSecret()
	parms["grant_type"] = "authorization_code"
	parms["code"] = code

	uri := "/common/oauth2/v2.0/token"
	URL := cli.SSOHost + uri

	token, err := HandleResponForParseToken(cli.HTTPClient.HttpFormPost(URL, nil, parms))
	if err != nil {
		core.Println("err=", err)
		return err
	}
	core.Println("first refresh token:", token.RefreshToken)
	cli.Token = token
	cli.SaveToken2HomeDefault(token)
	return nil
}

//UpdateToken update expried token
func (cli *OneClient) UpdateToken() (*AuthToken, error) {
	parms := cli.GetOneDriveAppInfoWithSecret()
	parms["grant_type"] = "refresh_token"
	parms["refresh_token"] = cli.Token.RefreshToken

	uri := "/common/oauth2/v2.0/token"
	URL := cli.SSOHost + uri
	token, err := HandleResponForParseToken(cli.HTTPClient.HttpFormPost(URL, nil, parms))
	if err != nil {
		core.Println("err=", err)
		return nil, err
	}
	//using old drive ID
	token.DriveID = cli.CurDriveID
	cli.Token = token
	err = cli.SaveToken2UserConfig(token)
	if err != nil {
		return nil, err
	}
	core.Println(token.AccessToken)
	return token, nil
}

//SaveToken2Home save token to local
func (cli *OneClient) SaveToken2UserConfig(token *AuthToken) error {
	exTime := time.Now().Add(time.Second * time.Duration(token.ExpiresIn-60))
	token.ExpiresTime = Timestamp(exTime)
	cli.SaveToken2Home(token)
	return nil
}

//SaveToken2Home save token to local
func (cli *OneClient) SaveToken2HomeDefault(token *AuthToken) {
	exTime := time.Now().Add(time.Second * time.Duration(token.ExpiresIn-60))
	token.ExpiresTime = Timestamp(exTime)
	SaveToken2DefaultPath(token)
}

//GetOneDriveAppInfo setup application information
func (cli *OneClient) GetOneDriveAppInfo() map[string]string {
	parms := map[string]string{}
	parms["client_id"] = CLIENT_ID
	parms["scope"] = SCOPE
	parms["redirect_uri"] = CALLBACK_URL
	return parms
}

//GetOneDriveAppInfoWithSecret with secret info
func (cli *OneClient) GetOneDriveAppInfoWithSecret() map[string]string {
	p := cli.GetOneDriveAppInfo()
	p["client_secret"] = CLIENT_SECRET
	return p
}

//SetOneDriveAPIToken for http request setup a token
func (cli *OneClient) SetOneDriveAPIToken() map[string]string {
	header := map[string]string{}
	header["Content-Type"] = "application/json"
	header["Authorization"] = "Bearer " + cli.Token.AccessToken
	return header
}

//APIGetUserInfo get user info from api
func (cli *OneClient) APIGetUserInfo() {
	uri := "/me"
	URL := cli.APIHost + uri
	header := cli.SetOneDriveAPIToken()
	json, err := chttp.HandleRespon2String(cli.HTTPClient.HttpGet(URL, header, nil))
	if err != nil {
		fmt.Println("err=", err)
		return
	}
	fmt.Println(json)
}

//APIGetMeDrive get onedrive infomation
func (cli *OneClient) APIGetMeDrive() (*Drive, error) {
	uri := "/me/drive"
	URL := cli.APIHost + uri
	header := cli.SetOneDriveAPIToken()
	dri := new(Drive)
	resp, err := cli.HTTPClient.HttpGet(URL, header, nil)
	err = HandleResponForParseAPI(resp, err, dri)
	if err != nil {
		fmt.Println("err=", err)
		return nil, err
	}
	core.Println("id=", dri.ID)
	return dri, nil
}

//APIListFilesByPath get files by path
func (cli *OneClient) APIListFilesByPath(driveID string, path string) (*ListChildrenResponse, error) {
	uri := "/drives/%s/root:%s:/children"
	URL := cli.APIHost + fmt.Sprintf(uri, driveID, path)
	if path == "/" {
		uri := "/drives/%s/root/children"
		URL = cli.APIHost + fmt.Sprintf(uri, driveID)
	}
	core.Println("APIListFilesByPath request url = ", URL)
	header := cli.SetOneDriveAPIToken()
	objs := new(ListChildrenResponse)
	resp, err := cli.HTTPClient.HttpGet(URL, header, nil)
	err = HandleResponForParseAPI(resp, err, objs)
	if err != nil {
		return nil, err
	}
	return objs, nil
}

//APISearchByKey search files by Key
func (cli *OneClient) APISearchByKey(driveID string, key string) (*ListChildrenResponse, error) {
	uri := "/drives/%s/root/search(q='%s')"
	URL := cli.APIHost + fmt.Sprintf(uri, driveID, key)
	core.Println("APISearchByKey request url = ", URL)
	header := cli.SetOneDriveAPIToken()
	objs := new(ListChildrenResponse)
	resp, err := cli.HTTPClient.HttpGet(URL, header, nil)
	err = HandleResponForParseAPI(resp, err, objs)
	if err != nil {
		return nil, err
	}
	return objs, nil
}

//APIGetFile get a file by file path
func (cli *OneClient) APIGetFile(driveID string, path string) (*Item, error) {
	URL := ""
	if path == "/" {
		uri := "/drives/%s/root"
		URL = cli.APIHost + fmt.Sprintf(uri, driveID)
	} else {
		uri := "/drives/%s/root:%s"
		URL = cli.APIHost + fmt.Sprintf(uri, driveID, path)
	}
	core.Println("URI = ", URL)
	header := cli.SetOneDriveAPIToken()
	objs := new(Item)
	resp, err := cli.HTTPClient.HttpGet(URL, header, nil)
	err = HandleResponForParseAPI(resp, err, objs)
	if err != nil {
		return nil, err
	}
	return objs, nil
}
func (cli *OneClient) APIGetFileByID(driveID string, ID string) (*Item, error) {
	uri := "/drives/%s/items/%s"
	URL := cli.APIHost + fmt.Sprintf(uri, driveID, ID)
	header := cli.SetOneDriveAPIToken()
	objs := new(Item)
	resp, err := cli.HTTPClient.HttpGet(URL, header, nil)
	err = HandleResponForParseAPI(resp, err, objs)
	if err != nil {
		return nil, err
	}
	return objs, nil
}

//APIDelFileByItemID delete file by item ID
func (cli *OneClient) APIUpdateFileByItemID(driveID string, itemID string, newName string, newPathID string) (bool, error) {
	uri := "/drives/%s/items/%s"
	URL := cli.APIHost + fmt.Sprintf(uri, driveID, itemID)
	header := cli.SetOneDriveAPIToken()

	bodyTmp := `{
  "parentReference": {
    "id": "%s"
  },
  "name": "%s"
}`
	bodyTmp = fmt.Sprintf(bodyTmp, newPathID, newName)
	core.Println("body =", bodyTmp)
	resp, err := cli.HTTPClient.HttpRequest("PATCH", URL, header, bodyTmp)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == 200 {
		return true, nil
	} else {
		return false, HandleResponForParseAPI(resp, nil, nil)
	}
}

//APIDelFileByItemID delete file by item ID
func (cli *OneClient) APIDelFileByItemID(driveID string, itemID string) (bool, error) {
	uri := "/drives/%s/items/%s"
	URL := cli.APIHost + fmt.Sprintf(uri, driveID, itemID)
	header := cli.SetOneDriveAPIToken()

	resp, err := cli.HTTPClient.HttpRequest("DELETE", URL, header, "")
	if err != nil {
		return false, err
	}
	if resp.StatusCode == 204 {
		return true, nil
	} else {
		return false, HandleResponForParseAPI(resp, nil, nil)
	}
}

//APIDelFile delete file by file path
func (cli *OneClient) APIDelFile(driveID string, filePath string) (bool, error) {
	uri := "/drives/%s/root:%s"
	URL := cli.APIHost + fmt.Sprintf(uri, driveID, filePath)
	header := cli.SetOneDriveAPIToken()

	resp, err := cli.HTTPClient.HttpRequest("DELETE", URL, header, "")
	if err != nil {
		return false, err
	}
	if resp.StatusCode == 204 {
		return true, nil
	} else {
		return false, HandleResponForParseAPI(resp, nil, nil)
	}
}

//APImkdir create a dir
func (cli *OneClient) APImkdir(driveID string, path string, dirName string) (*Item, error) {
	uri := "/drives/%s/root:%s:/children"
	URL := cli.APIHost + fmt.Sprintf(uri, driveID, path)
	if path == "/" {
		uri := "/drives/%s/root/children"
		URL = cli.APIHost + fmt.Sprintf(uri, driveID)
	}
	core.Println(URL)
	header := cli.SetOneDriveAPIToken()
	objs := new(Item)

	bodyTpl := `{
  "name": "%s",
  "folder": { },
  "@microsoft.graph.conflictBehavior": "rename"
}`
	body := fmt.Sprintf(bodyTpl, dirName)
	resp, err := cli.HTTPClient.HttpPost(URL, header, body)
	err = HandleResponForParseAPI(resp, err, objs)
	if err != nil {
		return nil, err
	}
	return objs, nil
}

//APIUploadText upload a text
func (cli *OneClient) APIUploadText(driveID string, path string, content string) (*Item, error) {
	uri := "/drives/%s/root:%s:/content"
	URL := cli.APIHost + fmt.Sprintf(uri, driveID, path)
	core.Println(URL)
	header := cli.SetOneDriveAPIToken()
	header["Content-Type"] = "text/plain"
	objs := new(Item)
	resp, err := cli.HTTPClient.HttpRequest("PUT", URL, header, content)
	err = HandleResponForParseAPI(resp, err, objs)
	if err != nil {
		return nil, err
	}
	return objs, nil
}

//HandleResponForParseToken parse token
func HandleResponForParseToken(resp *http.Response, err error) (*AuthToken, error) {
	if resp == nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	core.Println("token = ", string(buff))
	core.Println("statuscode = ", resp.StatusCode)
	if resp.StatusCode == 200 {
		token := new(AuthToken)
		perr := json.Unmarshal(buff, token)
		if perr != nil {
			return nil, perr
		}
		return token, nil
	} else {
		apiErr := new(AuthError)
		perr := json.Unmarshal(buff, apiErr)
		if perr != nil {
			return nil, perr
		}
		return nil, errors.New(apiErr.ErrorDescription)
	}
}

// HandleResponForParseAPI parse api
func HandleResponForParseAPI(resp *http.Response, err error, objs interface{}) error {
	if resp == nil {
		return err
	}
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	core.Println(string(buff))
	core.Println("code,", resp.StatusCode)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		perr := json.Unmarshal(buff, objs)
		if perr != nil {
			return perr
		}
	} else {
		apiErr := new(Error)
		perr := json.Unmarshal(buff, apiErr)
		if perr != nil {
			return perr
		}
		return errors.New(apiErr.ErrorInfo.Message)
	}
	return nil
}

//NewBaseOneClient for new user
func NewBaseOneClient() *OneClient {
	cli := new(OneClient)
	httpCli := chttp.NewHttpClient()
	cli.HTTPClient = httpCli
	cli.APIHost = "https://graph.microsoft.com/v1.0"
	cli.SSOHost = "https://login.microsoftonline.com"
	return cli
}

//NewOneClient instance a OneClient
func NewOneClient() (*OneClient, error) {
	u := getCurUser()
	return NewOneClientUser(u)
}
func NewOneClientUser(user string) (*OneClient, error) {
	cli := NewBaseOneClient()
	cli.setUserInfo(user)
	tk := cli.getConfigAuthToken()
	if tk == nil {
		return nil, errors.New("pls config a new user")
	}
	cli.Token = tk
	expires := time.Time(tk.ExpiresTime)
	if time.Now().After(expires) {
		core.Println("to expries time, update token")
		newToken, err := cli.UpdateToken()
		if err != nil {
			return nil, err
		}
		cli.Token = newToken
	}
	if tk != nil {
		cli.CurDriveID = tk.DriveID
	}
	return cli, nil
}
func (cli *OneClient) ExpresCheck() error {
	tk := cli.Token
	expires := time.Time(tk.ExpiresTime)
	if time.Now().After(expires) {
		core.Println("to expries time, update token")
		newToken, err := cli.UpdateToken()
		if err != nil {
			return err
		}
		cli.Token = newToken
	}
	return nil
}

//GetTokenHeader for other client
func (cli *OneClient) GetTokenHeader() map[string]string {
	return cli.SetOneDriveAPIToken()
}

func callShellCB(cmd string, URL ...string) error {
	mycmd := exec.Command(cmd, URL...)
	err := mycmd.Start()
	go func() {
		err = mycmd.Wait()
		if err != nil {
			fmt.Printf("Command finished with error: %v", err)
		}
	}()
	return err
}
func getQueryParamByKey(r *http.Request, key string) string {

	keys, ok := r.URL.Query()[key]
	if !ok || len(keys[0]) < 1 {
		return ""
	}

	return keys[0]
}

//DoAutoForNewUser config a new user
func (cli *OneClient) DoAutoForNewUser() {
	//open browser
	go func() {
		time.Sleep(time.Second * 2)
		autoURL := cli.GetAuthCodeURL()
		if runtime.GOOS == "linux" {
			callShellCB("xdg-open", autoURL)
		} else {
			autoURL = strings.ReplaceAll(autoURL, "&", "^&")
			callShellCB("cmd", "/C", "start", autoURL)
		}
	}()
	respURL := cli.GetOneDriveAppInfo()["redirect_uri"]
	u, _ := url.Parse(respURL)
	sm := http.NewServeMux()
	server := http.Server{Addr: u.Host, Handler: sm}
	sm.HandleFunc(u.Path, func(w http.ResponseWriter, r *http.Request) {
		dd := getQueryParamByKey(r, "code")
		fmt.Println("code=", dd)
		if dd == "" {
			return
		}
		cli.GetFirstToken(dd)
		w.Write([]byte("save token to local"))
		go func() {
			time.Sleep(time.Second * 3)
			server.Shutdown(context.Background())
		}()
	})
	err := server.ListenAndServe()
	if err != nil && http.ErrServerClosed == err {
		fmt.Println("done all")
		return
	}
	if err != nil {
		fmt.Println("run http server to failed,server err = ", err)
	}
}
func mytest() {

	//core.Debug = false

	cli, _ := NewOneClient()

	//cli.GetAuthCode()
	//cli.GetFirstToken()

	//cli.UpdateToken()
	//API##########

	//cli.APIListMessages()

	//cli.APIGetMsgContentByID("AQMkADAwATM3ZmYAZS0zZjNlLTQzZGEtMDACLTAwCgBGAAADuzV3I4FQqkmvQunlH6Hg5gcAHVFtafF-NUS5DgiI6rlG6QAAAgEJAAAAHVFtafF-NUS5DgiI6rlG6QACDP7sxgAAAA==")

	/*
		html := `
		<html>
		<head>
			<title>info</title>
		</head>
		<body>
			information to you
		</body>
		</html>
		`
		cli.APISendMail("canyelengxin@163.com", "html email test", html, "html")
	*/
	cli.findConfigFile()
}
