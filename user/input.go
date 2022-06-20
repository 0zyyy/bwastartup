package user

type RegisterUserInput struct {
	Nama       string `json:"nama" binding:"required"`
	Occupation string `json:"occupation" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type EmailUserInput struct {
	Email string `json:"email" binding:"required,email"`
}

type AvatarUserInput struct {
}
