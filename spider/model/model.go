package model

type Proxy struct {
	ID       string `json:"ID,omitempty"`
	Scheme   string `json:"Scheme,omitempty"`
	Host     string `json:"Host,omitempty"`
	Username string `json:"Username,omitempty"`
	Password string `json:"Password,omitempty"`
	Blocked  bool   `json:"Blocked,omitempty"`
}
