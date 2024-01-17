package models

type Error struct {
	Code int    `json:"code",yaml:"code"`
	Text string `json:"text",yaml:"text"`
}
