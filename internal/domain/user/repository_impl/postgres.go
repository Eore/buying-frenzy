package repository_impl

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	dbc *pgxpool.Pool
}

func NewPostgresRepo(dbc *pgxpool.Pool) Postgres {
	return Postgres{
		dbc: dbc,
	}
}

func (p Postgres) ProcessPurchased(ctx context.Context, userID int) error {
	querySelectUserBalance := `
		SELECT cash_balance FROM users
		WHERE id = $1
		FOR UPDATE
	`
	queryTotalPurchase := `
		SELECT SUM(transaction_amount) FROM purchase_histories
		WHERE user_id = $1
	`

	querySelectRestaurantTrx := `
		SELECT 
			restaurants.id AS restaurant_id,
			(restaurants.cash_balance + purchase_histories.transaction_amount) AS cash_balance_after
		FROM purchase_histories
		JOIN restaurants ON restaurants.name = purchase_histories.restaurant_name
		WHERE purchase_histories.user_id = $1
		FOR UPDATE
	`

	queryUpdateUserBalance := `
		UPDATE users SET cash_balance = $1 WHERE id = $2
	`

	queryUpdateRestaurantBalance := `
		UPDATE restaurants SET cash_balance = $1 WHERE id = $2
	`

	tx, err := p.dbc.Begin(ctx)
	if err != nil {
		return err
	}

	var cashBalance, totalPurchase float64

	tx.QueryRow(ctx, querySelectUserBalance, userID).Scan(&cashBalance)
	tx.QueryRow(ctx, queryTotalPurchase, userID).Scan(&totalPurchase)

	if totalPurchase > cashBalance {
		return fmt.Errorf("balance not enough")
	}

	type restTrx struct {
		restaurantID     string
		cashBalanceAfter float64
	}

	rows, _ := tx.Query(ctx, querySelectRestaurantTrx, userID)
	defer rows.Close()

	trxs := []restTrx{}
	for rows.Next() {
		var s restTrx
		if err := rows.Scan(
			&s.restaurantID,
			&s.cashBalanceAfter,
		); err != nil {
			tx.Rollback(ctx)
			return err
		}

		trxs = append(trxs, s)
	}

	for _, trx := range trxs {
		if _, err := tx.Exec(ctx, queryUpdateRestaurantBalance, trx.cashBalanceAfter, trx.restaurantID); err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	if _, err := tx.Exec(ctx, queryUpdateUserBalance, cashBalance-totalPurchase, userID); err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
