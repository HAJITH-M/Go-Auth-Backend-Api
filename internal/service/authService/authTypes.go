package authservice

type RegisterInput struct {
	Email       string
	Password    string
	DisplayName string
}

type RegisterResult struct {
	Email       string
	DisplayName string
	UserId      string
}

type LoginInput struct {
	Email     string
	Password  string
	IPAddress string
	Device    string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	UserId       string
	Email        string
	DisplayName  string
}

type ChangePasswordInput struct {
	Email       string
	OldPassword string
	NewPassword string
}
