package one

import (
	"fmt"

	chttp "github.com/weilin88/notify2y/http"
)

//APISendMail get user info from api
func (cli *OneClient) APISendMail() error {
	uri := "/me/sendMail"
	URL := cli.APIHost + uri

	header := cli.SetOneDriveAPIToken()
	bodyTpl := `
{
    "message": {
        "subject": "test to send mail",
        "body": {
            "contentType": "Text",
            "content": "mail content for you"
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
	body := fmt.Sprintf(bodyTpl, "canyelengxin@163.com")
	fmt.Println(body)
	resp, err := cli.HTTPClient.HttpPost(URL, header, body)
	//err = HandleResponForParseAPI(resp, err, objs)
	json, err := chttp.HandleRespon2String(resp, err)
	fmt.Println("json data = ", json)
	if err != nil {
		fmt.Println("err=", err)
		return err
	}
	return err
}
