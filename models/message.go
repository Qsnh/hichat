package models

type MessageUser struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type Message struct {
	User MessageUser `json:"user"`
	C    string      `json:"c"`
}

type ReceiveMessage struct {
	Type  string `json:"t"`
	Value string `json:"v"`
}
