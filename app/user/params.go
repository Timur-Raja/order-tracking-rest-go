package user

type UserCreateParams struct {
	Email    *string `json:"email" binding:"required,email"`
	Password *string `json:"password" binding:"required,min=8"`
	Name     *string `json:"name" binding:"required,min=3"`
}
