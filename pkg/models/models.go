package models

type Message struct {
	Destination    string
	Email          string
	Username       string
	MessageSubject string
	Message        string
}

type MessagePush struct {
	MessageSubject string
	Message        string
}
