package web

type AuthRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
