package web

import "github.com/letenk/use_deal_user/models/domain"

type UserCreateRequest struct {
	Fullname string `json:"fullname" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role,omitempty"`
}

type UserUpdateRequest struct {
	Fullname string `json:"fullname,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

type UserGetIDUri struct {
	ID string `uri:"id" binding:"required"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// Format for handle single response user
func FormatUserResponse(user domain.User) UserResponse {
	formatter := UserResponse{
		ID:       user.ID,
		Fullname: user.Fullname,
		Username: user.Username,
		Role:     user.Role,
	}
	return formatter
}

// Format for handle multiples response users
func FormatUsersResponse(user []domain.User) []UserResponse {
	if len(user) == 0 {
		return []UserResponse{}
	}

	var formatters []UserResponse

	for _, data := range user {
		formatter := FormatUserResponse(data)
		formatters = append(formatters, formatter)
	}

	return formatters
}
