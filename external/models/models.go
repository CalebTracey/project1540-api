package models

type Request struct {
	ID string `json:"id,omitempty"`
}

type Message struct {
	Hostname  string
	Time      string
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
