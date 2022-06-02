package main

import (
	"context"
	"glints/internal"
	"glints/internal/api"
	restaurant "glints/internal/domain/restaurant/repository_impl"
	user "glints/internal/domain/user/repository_impl"
	"glints/pkg"
	"os"
)

var (
	DBHost = os.Getenv("DBHOST")
	DBPort = os.Getenv("DBPORT")
	DBUser = os.Getenv("DBUSER")
	DBPass = os.Getenv("DBPASS")
	DBName = os.Getenv("DBNAME")
)

func main() {
	ctx := context.Background()
	dbc := pkg.PostgresConnect(ctx, DBHost, DBPort, DBName, DBUser, DBPass)
	defer dbc.Close()

	httpAPI := api.NewHTTPAPI(internal.NewUsecase(
		restaurant.NewPostgresRepo(dbc),
		user.NewPostgresRepo(dbc),
	))

	httpAPI.StartServer(8000)
}
