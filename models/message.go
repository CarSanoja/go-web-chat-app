package models

type Message struct {
	Username string `json:"username"`
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	Data     string `json:"data,omitempty"`
}
