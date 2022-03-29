package message

import (
	"bytes"
	"encoding/json"
	"github.com/dimall-id/lumos/v2/errs"
	"github.com/dimall-id/lumos/v2/event"
	"gorm.io/gorm"
)

type Destination struct {
	Address string `json:"address"`
	Type    string `json:"type"`
}

type Message struct {
	// Origin use to state the name of the sender. this only used for email.
	Origin      string                 `json:"origin"`
	Destination Destination            `json:"destination"`
	Subject     string                 `json:"subject"`
	Template    string                 `json:"template"`
	CountryCode string                 `json:"country_code"`
	Data        map[string]interface{} `json:"data"`
}

func SendMessageCommand(message Message, tx *gorm.DB) error {
	jsonRes, err := json.Marshal(message)
	var res bytes.Buffer
	err = json.Compact(&res, jsonRes)
	if err != nil {
		return errs.DecodeJsonErr(err)
	}
	msg := event.LumosMessage{
		Topic:   "SEND_MESSAGE_COMMAND",
		Key:     "DATA",
		Value:   res.String(),
		Headers: map[string]string{},
	}
	err = event.GenerateOutbox(tx, msg)
	if err != nil {
		return errs.GenerateOutboxErr(err)
	}
	return nil
}
