package userapi

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/timur-raja/order-tracking-rest-go/app"
	"github.com/timur-raja/order-tracking-rest-go/app/user"
	"github.com/timur-raja/order-tracking-rest-go/app/user/usersql"
	"golang.org/x/crypto/bcrypt"
)

type userSigninHandler struct {
	db     *pgxpool.Pool
	params *user.UserSigninParams
}

func UserSigninHandler(db *pgxpool.Pool) gin.HandlerFunc {
	// Initialize the handler struct with the db connection
	return func(c *gin.Context) {
		h := &userSigninHandler{
			db:     db,
			params: new(user.UserSigninParams),
		}
		h.exec(c)
	}
}

func (h *userSigninHandler) exec(c *gin.Context) {
	// load the user params from the request body
	h.params = new(user.UserSigninParams)
	if err := c.ShouldBindJSON(h.params); err != nil {
		app.AbortWithErrorResponse(c, app.ErrFailedToLoadParams, err)
		return
	}

	// sanitize the params against xss attacks
	h.params.Sanitize()

	password := h.params.Password
	email := strings.ToLower(strings.TrimSpace(h.params.Email))

	// fetch the user associated to the email in the request
	query := usersql.NewSelectUserByEmailQuery(h.db)
	query.Where.Email = email
	if err := query.Run(c); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			app.AbortWithErrorResponse(c, user.ErrUserNotFound, err)
			return
		}
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}

	userInfo := query.User

	// compare the passwords
	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(password)); err != nil {
		app.AbortWithErrorResponse(c, user.ErrInvalidCredentials, err)
		return
	}

	// create session token
	token, err := user.GenerateSessionToken(32)
	if err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}

	// insert user session
	query2 := usersql.NewInsertSessionQuery(h.db)
	query2.Values.Session.Token = token
	query2.Values.UserID = userInfo.ID
	query2.Values.Session.CreatedAt = time.Now()
	if err := query2.Run(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}

	// set cookie in response
	c.SetCookie(
		"session_token",
		token,
		int(90*24*time.Hour.Seconds()),
		"/",
		"",
		true,
		true,
	)

	query3 := usersql.NewSelectUserViewByIDQuery(h.db)
	query3.Where.ID = userInfo.ID
	if err := query3.Run(c); err != nil {
		app.AbortWithErrorResponse(c, app.ErrServerError, err)
		return
	}

	c.JSON(201, query3.UserView)
}
