package userapi

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/timur-raja/order-tracking-rest-go/app"
	"github.com/timur-raja/order-tracking-rest-go/app/user"
	"github.com/timur-raja/order-tracking-rest-go/app/user/usersql"
	"golang.org/x/crypto/bcrypt"
)

type userCreateHandler struct {
	connections *app.Services
	params      *user.UserCreateParams
	newUser     *user.User
}

func UserCreateHandler(services *app.Services) gin.HandlerFunc {
	// Initialize the handler struct with the db connection
	return func(c *gin.Context) {
		h := &userCreateHandler{
			connections: services,
			params:      new(user.UserCreateParams),
		}
		h.exec(c)
	}
}

func (h *userCreateHandler) exec(c *gin.Context) {
	// load the user params from the request body
	if err := c.ShouldBindJSON(h.params); err != nil {
		app.AbortWithErrorResponse(c, app.ErrFailedToLoadParams, err)
		return
	}

	// sanitize the params against xss attacks
	h.params.Sanitize()

	// validate user params and populate new user struct
	if err := h.buildUser(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}

	// execute the user insert query with values read and validated from the request body
	query := usersql.NewInsertUserQuery(h.connections.DB)
	query.Values.User = *h.newUser
	if err := query.Run(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}

	// fetch the view of the newly created user to send to the FE
	query2 := usersql.NewSelectUserViewByIDQuery(h.connections.DB)
	query2.Where.ID = query.Returning.ID
	if err := query2.Run(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}

	c.JSON(201, query2.UserView)
}

func (h *userCreateHandler) buildUser(c *gin.Context) error {
	h.newUser = new(user.User)

	h.newUser.Email = strings.ToLower(strings.TrimSpace(h.params.Email))
	h.newUser.Name = strings.TrimSpace(h.params.Name)

	// check if the email already exists
	query := usersql.NewSelectEmailExistsQuery(h.connections.DB)
	query.Where.Email = h.newUser.Email
	if err := query.Run(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return app.ErrServerError.Err
	}

	if query.Exists {
		app.AbortWithErrorResponse(c, user.ErrUserAlreadyExists, user.ErrUserAlreadyExists.Err)
		return user.ErrUserAlreadyExists.Err
	}

	// hash password
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(h.params.Password), bcrypt.DefaultCost,
	)
	if err != nil {
		app.AbortWithErrorResponse(c, user.ErrInvalidCredentials, err)
		return user.ErrInvalidCredentials.Err
	}

	h.newUser.Password = string(passwordHash)
	return nil
}
