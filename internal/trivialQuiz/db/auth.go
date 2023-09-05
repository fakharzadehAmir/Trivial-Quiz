package db

import (
	"Trivia_Quiz/pkg/authenticate"
	"context"
)

func (mdb *MongoDB) GetAccountByUsername(ctx *context.Context, username string) (*authenticate.Account, error) {
	user, err := mdb.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return &authenticate.Account{
		Username: user.Username,
		Password: user.Password,
	}, nil
}
