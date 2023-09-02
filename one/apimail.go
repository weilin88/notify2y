package one

import (
	"fmt"

	"github.com/weilin88/notify2y/core"
	chttp "github.com/weilin88/notify2y/http"
)

//APISendMail get user info from api
func (cli *OneClient) APISendMail(person, subject, content string, contentType string) error {
	uri := "/me/sendMail"
	URL := cli.APIHost + uri

	header := cli.SetOneDriveAPIToken()
	bodyTpl := `
{
    "message": {
        "subject": "%s",
        "body": {
            "contentType": "%s",
            "content": "%s"
        },
        "toRecipients": [
            {
                "emailAddress": {
                    "address": "%s"
                }
            }
        ]
    },
    "saveToSentItems": "true"
}
`
	if contentType != "html" {
		contentType = "text"
	}
	body := fmt.Sprintf(bodyTpl, EscapeJSONString(subject), contentType, EscapeJSONString(content), person)
	fmt.Println(body)
	resp, err := cli.HTTPClient.HttpPost(URL, header, body)
	if resp.StatusCode == 202 && err == nil {
		return nil
	}
	json, err := chttp.HandleRespon2String(resp, err)
	core.Println("json data = ", json)
	if err != nil {
		return err
	}
	return fmt.Errorf("status code = %d", resp.StatusCode)
}

func (cli *OneClient) APIListMessages() (*ListMessagesResponse, error) {
	uri := "/me/messages"
	URL := cli.APIHost + uri

	header := cli.SetOneDriveAPIToken()
	objs := new(ListMessagesResponse)
	resp, err := cli.HTTPClient.HttpGet(URL, header, nil)
	err = HandleResponForParseAPI(resp, err, objs)
	if err != nil {
		return nil, err
	}
	return objs, nil
}

func (cli *OneClient) APIGetMessageByID(ID string) (*EMessage, error) {
	uri := "/me/messages/%s"
	uri = fmt.Sprintf(uri, ID)
	URL := cli.APIHost + uri

	header := cli.SetOneDriveAPIToken()
	objs := new(EMessage)
	resp, err := cli.HTTPClient.HttpGet(URL, header, nil)
	err = HandleResponForParseAPI(resp, err, objs)
	if err != nil {
		return nil, err
	}
	return objs, nil
}

func (cli *OneClient) APIGetMsgContentByID(ID string) (*EMessage, error) {
	uri := "/me/messages/%s/$value"
	uri = fmt.Sprintf(uri, ID)
	URL := cli.APIHost + uri

	fmt.Println("URL=", URL)
	header := cli.SetOneDriveAPIToken()

	//header["Accept"] = "text/plain"

	//objs := new(EMessage)
	resp, err := cli.HTTPClient.HttpGet(URL, header, nil)
	//err = HandleResponForParseAPI(resp, err, objs)
	json, err := chttp.HandleRespon2String(resp, err)
	fmt.Println("json = ", json)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
