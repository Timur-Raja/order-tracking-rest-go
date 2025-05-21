package user

import "text/template"

type UserCreateParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=3"`
}

func (p *UserCreateParams) Sanitize() {
	p.Email = template.HTMLEscapeString(p.Email)
	p.Name = template.HTMLEscapeString(p.Name)
	p.Password = template.HTMLEscapeString(p.Password)
}

type UserSigninParams struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (p *UserSigninParams) Sanitize() {
	p.Email = template.HTMLEscapeString(p.Email)
	p.Password = template.HTMLEscapeString(p.Password)
}
