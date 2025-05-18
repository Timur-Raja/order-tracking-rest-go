package usersql

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/timur-raja/order-tracking-rest-go/app/user"
	"github.com/timur-raja/order-tracking-rest-go/db"
)

type insertSessionQuery struct {
	db.BaseQuery
	Values struct {
		user.Session
	}
	Returning struct {
		Token string `db:"token"`
	}
}

func NewInsertSessionQuery(conn db.PGExecer) *insertSessionQuery {
	return &insertSessionQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
		Values:    struct{ user.Session }{Session: user.Session{}},
	}
}

func (q *insertSessionQuery) Run(ctx context.Context) error {
	query := `INSERT INTO user_sessions (token, user_id, created_at) 
	VALUES($1, $2, $3)
	RETURNING token;`

	if err := pgxscan.Get(ctx, q.DBConn, &q.Returning.Token, query,
		q.Values.Session.Token,
		q.Values.Session.UserID,
		q.Values.Session.CreatedAt); err != nil {
		return err
	}
	return nil
}

type selectSessionByTokenQuery struct {
	db.BaseQuery
	*user.Session
	Where struct {
		Token string `db:"token"`
	}
}

func NewSelectSessionByTokenQuery(conn db.PGExecer) *selectSessionByTokenQuery {
	return &selectSessionByTokenQuery{
		BaseQuery: db.BaseQuery{DBConn: conn},
		Session:   &user.Session{},
	}
}

func (q *selectSessionByTokenQuery) Run(ctx context.Context) error {
	query := `SELECT *
	FROM user_sessions
	WHERE token = $1;`

	if err := pgxscan.Get(ctx, q.DBConn, q.Session, query,
		q.Where.Token); err != nil {
		return err
	}
	return nil
}
