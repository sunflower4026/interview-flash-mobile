package auth

import (
	"github.com/google/uuid"
	userDomain "gitlab.com/sunflower4026/interview-flash-mobile/model/domain/user"
)

type RegisterResponse struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	CreatedAt   string    `json:"created_at"`
}

func ToRegisterResponse(user userDomain.User) RegisterResponse {
	return RegisterResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
		CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
