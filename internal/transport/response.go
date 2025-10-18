package transport

type ApiResponse struct {
	Ok      bool        `json:"ok"`
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Time    string      `json:"time"`
}
