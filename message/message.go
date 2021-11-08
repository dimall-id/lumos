package message

import (
	"bytes"
	"encoding/json"
	"github.com/dimall-id/lumos/v2/errs"
	"github.com/dimall-id/lumos/v2/event"
	"gorm.io/gorm"
)

type SendTo struct {
	SendTo string `json:"send_to"`
	SendType string `json:"send_type"`
}

type Message struct {
	SendTo      []SendTo               `json:"send_to"`
	Subject     string                 `json:"subject"`
	Template    string                 `json:"template"`
	CountryCode string                 `json:"country_code"`
	Data        map[string]interface{} `json:"data"`
}

func SendMessageCommand (message Message, tx *gorm.DB) error {
	jsonRes, err := json.Marshal(message)
	var res bytes.Buffer
	err = json.Compact(&res, jsonRes)
	if err != nil {return errs.DecodeJsonErr(err)}
	msg := event.LumosMessage{
		Topic: "SEND_MESSAGE_COMMAND",
		Key: "DATA",
		Value: res.String(),
		Headers: map[string]string{},
	}
	err = event.GenerateOutbox(tx, msg)
	if err != nil {return errs.GenerateOutboxErr(err)}
	return nil
}
