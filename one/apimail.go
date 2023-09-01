package one

import (
	"fmt"

	"github.com/weilin88/notify2y/core"
	chttp "github.com/weilin88/notify2y/http"
)

//APISendMail get user info from api
func (cli *OneClient) APISendMail(person, subject, content string) error {
	uri := "/me/sendMail"
	URL := cli.APIHost + uri

	header := cli.SetOneDriveAPIToken()
	bodyTpl := `
{
    "message": {
        "subject": "%s",
        "body": {
            "contentType": "Text",
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
	body := fmt.Sprintf(bodyTpl, EscapeJSONString(subject), EscapeJSONString(content), person)
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
