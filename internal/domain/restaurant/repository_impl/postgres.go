package repository_impl

import (
	"context"
	"fmt"
	"glints/internal/domain/restaurant"

	"github.com/jackc/pgx/v4"
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

func (p Postgres) GetRestaurants(ctx context.Context, filter restaurant.Filter) (restaurant.RestaurantSchemas, error) {
	queryDefault := `
		SELECT restaurants.id, restaurants.name, restaurants.cash_balance 
		FROM restaurants
	`

	queryOpening := `
		SELECT restaurants.id, restaurants.name, restaurants.cash_balance 
		FROM restaurants
		JOIN opening_hours on opening_hours.restaurant_id = restaurants.id
		WHERE TRUE
			AND (opening_hours.day = $1 OR $1 = '')
			AND (opening_hours.start_time >= $2 OR $2 = '')
			AND (opening_hours.end_time <= $3 OR $3 = '')
	`

	queryMenus := `
		SELECT DISTINCT ON (rst.id) rst.id, rst.name, rst.cash_balance 
		FROM (
			SELECT 
				restaurants.id, 
				restaurants.name, 
				restaurants.cash_balance,
				(
					SELECT COUNT(menus.id)
					FROM menus
					WHERE menus.restaurant_id = restaurants.id
				) AS ndishes
			FROM restaurants
		) AS rst
		JOIN menus on menus.restaurant_id = rst.id
		WHERE TRUE
			AND (rst.ndishes = $1 OR $1 = 0)
			AND (menus.price >= $2 OR $2 = 0)
			AND (menus.price <= $3 OR $3 = 0)
	`

	var (
		rows pgx.Rows
		err  error
	)

	switch {
	case filter.Day != "" || filter.StartTime != "" || filter.EndTime != "":
		rows, err = p.dbc.Query(ctx, queryOpening, filter.Day, filter.DishPriceStart, filter.DishPriceEnd)
	case filter.NDishes > 0 || filter.DishPriceStart > 0 || filter.DishPriceEnd > 0:
		rows, err = p.dbc.Query(ctx, queryMenus, filter.NDishes, filter.DishPriceStart, filter.DishPriceEnd)
	default:
		rows, err = p.dbc.Query(ctx, queryDefault)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	restaurants := restaurant.RestaurantSchemas{}
	for rows.Next() {
		r := restaurant.RestaurantSchema{}
		if err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.CashBalance,
		); err != nil {
			return nil, err
		}

		restaurants = append(restaurants, r)
	}

	return restaurants, nil
}

func (p Postgres) GetRestaurantOpeningHours(ctx context.Context, restaurantID string) (restaurant.OpeningHourSchemas, error) {
	query := `
		SELECT id, restaurant_id, day, start_time, end_time
		FROM opening_hours
		WHERE restaurant_id = $1
	`

	rows, err := p.dbc.Query(ctx, query, restaurantID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	openings := restaurant.OpeningHourSchemas{}
	for rows.Next() {
		o := restaurant.OpeningHourSchema{}
		if err := rows.Scan(
			&o.ID,
			&o.RestaurantID,
			&o.Day,
			&o.Start,
			&o.End,
		); err != nil {

			fmt.Println(err)
			return nil, err
		}

		openings = append(openings, o)
	}

	return openings, nil
}

func (p Postgres) GetRestaurantMenus(ctx context.Context, restaurantID string) (restaurant.MenuSchemas, error) {
	query := `
		SELECT id, restaurant_id, dish_name, price
		FROM menus
		WHERE restaurant_id = $1
	`

	rows, err := p.dbc.Query(ctx, query, restaurantID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	menus := restaurant.MenuSchemas{}
	for rows.Next() {
		m := restaurant.MenuSchema{}
		if err := rows.Scan(
			&m.ID,
			&m.RestaurantID,
			&m.DishName,
			&m.Price,
		); err != nil {
			return nil, err
		}

		menus = append(menus, m)
	}

	return menus, nil
}
