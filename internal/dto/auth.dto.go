package dto

type AuthResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Token string `json:"token"`
}

type AuthRequest struct {
	Email    string `json:"email" example:"test@example.com" binding:"required"`
	Password string `json:"password" example:"Rahasia1!" binding:"required"`
}
