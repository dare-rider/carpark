package presentors

type DefaultResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
