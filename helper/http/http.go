package httputils

// StandardEnvelope is standard JSON HTTP response envelope.
type StandardEnvelope struct {
	Header *StandardHeader `json:"header,omitempty"`
	Status *StandardStatus `json:"status,omitempty"`
	Data   interface{}     `json:"data,omitempty"`
	Errors []StandardError `json:"errors,omitempty"`
}

// StandardHeader is standard JSON header HTTP response.
type StandardHeader struct {
	TotalData   int                    `json:"total_data"`
	ProcessTime float64                `json:"process_time"`
	Meta        map[string]interface{} `json:"meta"`
}

// StandardStatus is standard JSON status http response.
type StandardStatus struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

// StandardError is standard JSON HTTP Error.
type StandardError struct {
	Code   string      `json:"code"`
	Title  string      `json:"title"`
	Detail string      `json:"detail"`
	Object ErrorObject `json:"object"`
}

// ErrorObject holds any additional details of an error.
type ErrorObject struct {
	Text []string `json:"text"`
	Type int64    `json:"type"`
}
