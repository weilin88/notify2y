package task

import (
	"github.com/weilin88/notify2y/one"
)

type Task struct {
	ID         string `json:"id"`
	Subject    string `json:"subject"`
	Type       string `json:"type"`
	Content    string `json:"Content"`
	Importance string `json:"importance"`

	Deadline one.Timestamp `json:"sentDateTime"`

	CreatedDateTime      one.Timestamp `json:"createdDateTime"`
	LastModifiedDateTime one.Timestamp `json:"lastModifiedDateTime"`
}
