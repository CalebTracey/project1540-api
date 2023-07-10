package models

type APIResponse struct {
	Results Responses `json:"results,omitempty"`
	Message Message   `json:"message,omitempty"`
}

type Responses []Response
type Response struct {
	ID string `json:"id"`
}

type Message struct {
	Status    string    `json:"status,omitempty"`
	ErrorLogs ErrorLogs `json:"error_logs,omitempty"`
}
type ErrorLogs []ErrorLog
type ErrorLog struct {
	Status     string `json:"status,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
	RootCause  string `json:"root_cause,omitempty"`
	Details    string `json:"details,omitempty"`
	Trace      string `json:"trace,omitempty"`
}
