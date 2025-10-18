package handler

type ApiResponse struct {
	Ok      bool        `json:"ok"`
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Time    string      `json:"time"`
}
type ApiPayload struct {
	Path  string            `json:"path"`
	Param map[string]string `json:"param,omitempty"`
	Body  any               `json:"body,omitempty"`
}
