package models

type Alert struct {
	Labels map[string]string `json:"labels"`
	Status Status            `json:"status"`
}

type Status struct {
	State string `json:"state"`
}
type Email struct {
	Sender   string
	Receiver string
	Subject  string
	Body     string
}

type EmailCredential struct {
	Host     string
	Port     int
	Username string
	Password string
}

type ProxmoxCredential struct {
	API   string
	Token string
}
