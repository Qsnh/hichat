package main

type MessageUser struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type Message struct {
	User MessageUser `json:"user"`
	T    string      `json:"t"`
	C    string      `json:"c"`
}
