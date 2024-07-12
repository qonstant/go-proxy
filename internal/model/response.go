package model

// ResponseData represents the JSON structure of the server's response
type ResponseData struct {
	ID      string            `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Length  int               `json:"length"`
	Body    string            `json:"body"`
}
