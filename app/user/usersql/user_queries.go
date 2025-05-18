package usersql

import (
	"context"
	"time"

	"github.com/timur-raja/order-tracking-rest-go/app/user"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/timur-raja/order-tracking-rest-go/db"
)

type insertUserQuery struct {
	db.BaseQuery
	Values struct {
		user.User
	}
	Returning struct {
		ID int `db:"id"`
	}
}

func NewInsertUserQuery(conn db.PGExecer) *insertUserQuery {
	return &insertUserQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
	}
}

func (q *insertUserQuery) Run(ctx context.Context) error {
	query := `
        INSERT INTO users (email, password, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
    `
	// Execute the query and scan the returned ID into q.Returned.ID
	err := pgxscan.Get(ctx, q.DBConn, &q.Returning, query,
		q.Values.User.Email,
		q.Values.User.Password,
		q.Values.User.Name,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

type SelectUserByEmailQuery struct {
	db.BaseQuery
	Where struct {
		Email string
	}
	*user.User
}

func NewSelectUserByEmailQuery(conn db.PGExecer) *SelectUserByEmailQuery {
	return &SelectUserByEmailQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
		User:      &user.User{},
	}
}

func (q *SelectUserByEmailQuery) Run(ctx context.Context) error {
	query := `
	SELECT *
	FROM users
	WHERE email = $1`

	err := pgxscan.Get(ctx, q.DBConn, q.User, query,
		q.Where.Email,
	)
	if err != nil {
		return err
	}
	return nil
}

// the user view will be the user object that we return to the FE, including eventual necessary extra fields from other tables
type selectUserViewByIDQuery struct {
	db.BaseQuery
	Where struct {
		ID int
	}
	*user.UserView
}

func NewSelectUserViewByIDQuery(Conn db.PGExecer) *selectUserViewByIDQuery {
	return &selectUserViewByIDQuery{
		BaseQuery: db.BaseQuery{DBConn: Conn},
		UserView:  &user.UserView{},
	}
}

func (q *selectUserViewByIDQuery) Run(ctx context.Context) error {
	q.UserView = new(user.UserView)
	query := `
        SELECT * FROM users_view WHERE id = $1;
    `
	// Execute the query and scan the returned ID into q.Returned.ID
	err := pgxscan.Get(ctx, q.DBConn, q.UserView, query,
		q.Where.ID,
	)

	if err != nil {
		return err
	}
	return nil
}
