package entity

type ClientData struct {
	Type       string   `json:"type"`
	Clientname string   `json:"name"`
	PolicyID   string   `json:"policy"`
	Data       []string `json:"ignorelist"`
}
