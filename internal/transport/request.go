package transport

type ApiPayload struct {
	Path  string            `json:"path"`
	Param map[string]string `json:"param,omitempty"`
	Body  any               `json:"body,omitempty"`
}
