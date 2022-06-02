package restaurant

import "context"

type Filter struct {
	Day            string
	StartTime      string
	EndTime        string
	NDishes        int
	DishPriceStart float64
	DishPriceEnd   float64
}

type Repository interface {
	GetRestaurants(context.Context, Filter) (RestaurantSchemas, error)
	GetRestaurantOpeningHours(ctx context.Context, restaurantID string) (OpeningHourSchemas, error)
	GetRestaurantMenus(ctx context.Context, restaurantID string) (MenuSchemas, error)
}
