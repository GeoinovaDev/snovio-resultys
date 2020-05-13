package snovio

// Status ...
type Status struct {
	Description string `json:"description"`
	Identifier  string `json:"identifier"`
}

// ProtocolData ...
type ProtocolData struct {
	SMTPStatus string `json:"smtpStatus"`
}

// ProtocolStatus ...
type ProtocolStatus struct {
	Data   ProtocolData `json:"data"`
	Status Status       `json:"status"`
}
