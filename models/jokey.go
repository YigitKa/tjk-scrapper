package jokey

import (
	"encoding/json"
	"fmt"
)

type Jokey struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	SuspendDate string `json:"suspendDate"`
	DueDate     string `json:"dueDate"`   // TODO: change type time.Time?
	BanReason   string `json:"banReason"` // TODO: change type time.Time?
}

func ToJson(jokeys []Jokey) (s string) {
	e, err := json.Marshal(jokeys)
	if err != nil {
		fmt.Println(err)
		return
	}

	return string(e)
}
