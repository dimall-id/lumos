package message

type Message struct {
	SendTo string `json:"send_to"`
	SendType string `json:"send_type"`
	Subject string `json:"subject"`
	Template string `json:"template"`
	Data map[string]interface{} `json:"data"`
}
