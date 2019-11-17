package presentors

// Ping defines json struct for ping response
type Ping struct {
	Status      string `json:"status,omitempty"`
	Environment string `json:"environment"`
}
