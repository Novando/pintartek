package structs

type StdResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
	Count   int64       `json:"count,omitempty"`
}
