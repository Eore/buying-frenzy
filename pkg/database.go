package pkg

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

func PostgresConnect(ctx context.Context, host, port, database, username, password string) *pgxpool.Pool {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", username, password, host, port, database)
	dbc, err := pgxpool.Connect(ctx, dsn)
	if err := dbc.Ping(ctx); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return dbc
}
