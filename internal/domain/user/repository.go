package user

import "context"

type Repository interface {
	ProcessPurchased(ctx context.Context, userID int) error
}
