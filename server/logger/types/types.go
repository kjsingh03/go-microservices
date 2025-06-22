package types

type JsonResponse struct {
	Success bool 	`json:"success"`
	Message string	`json:"message"`
	Data any 		`json:"data,omitempty"`
	Error string 	`json:"error,omitempty"`
}