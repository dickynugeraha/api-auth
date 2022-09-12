package domains

type UserId struct {
	ID     string `json:"id"`
}

type Login struct {
	Name     string
	Email    string
	Password string
}

type Register struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
}

type ChangePassword struct {
	Email           string
	NewPassword     string `json:"newPassword"`
	PasswordConfirm string `json:"passwordConfirm"`
}