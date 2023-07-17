package models

type Message struct {
	Hostname  string    `json:"hostname,omitempty"`
	Time      string    `json:"time,omitempty"`
	Status    string    `json:"status,omitempty"`
	Count     string    `json:"count,omitempty"`
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
