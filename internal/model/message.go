package model

type Message struct {
	Command string      `json:"command"`
	Data    interface{} `json:"data"`
}