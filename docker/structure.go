package main

type NintendoChannelWebPayload struct {
	Arguments    string `json:"arguments"`
	SocketSecret string `json:"secret"`
}

type SocketFailResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type Config struct {
	SocketSecret string `xml:"socketSecret"`
}
