package types

type JsonResponse struct {
	Error   bool   `json:"error"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   *AuthPayload `json:"auth,omitempty"`
	Log    *LogPayload  `json:"log,omitempty"`
	Mail   *MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type AuthResponse struct {
	Valid bool   `json:"valid"`
	User  *User  `json:"user,omitempty"`
	Token string `json:"token,omitempty"`
}