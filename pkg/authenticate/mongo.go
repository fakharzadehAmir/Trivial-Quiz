package authenticate

import "context"

type MongoDB interface {
	GetAccountByUsername(cts *context.Context, username string) (*Account, error)
}
