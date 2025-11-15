package user

type SetIsActiveRequest struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type UserResponse struct {
	User User `json:"user"`
}

type User struct {
	Id       string `json:"user_id"`
	Name     string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}
