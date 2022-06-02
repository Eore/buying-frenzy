package internal

import (
	"context"
	"glints/internal/domain/restaurant"
	"glints/internal/domain/user"
	"log"
)

type Usecase struct {
	restaurantRepo restaurant.Repository
	userRepo       user.Repository
}

type RestaurantFilter struct {
	Day            string
	StartTime      string
	EndTime        string
	NDishes        int
	DishPriceStart float64
	DishPriceEnd   float64
}

type Restaurant struct {
	ID           string                  `json:"id"`
	Name         string                  `json:"name"`
	OpeningHours restaurant.OpeningHours `json:"openingHours"`
	CashBalance  float64                 `json:"cashBalance"`
	Menu         []restaurant.Menu       `json:"menu"`
}

func NewUsecase(restaurantRepo restaurant.Repository, userRepo user.Repository) Usecase {
	return Usecase{restaurantRepo: restaurantRepo, userRepo: userRepo}
}

func (u Usecase) RestaurantList(filter RestaurantFilter) []Restaurant {
	restaurants, err := u.restaurantRepo.GetRestaurants(context.Background(), restaurant.Filter(filter))
	if err != nil {
		log.Println(err)
		return []Restaurant{}
	}

	res := []Restaurant{}

	for _, rst := range restaurants {
		opnHours, _ := u.restaurantRepo.GetRestaurantOpeningHours(context.Background(), rst.ID)
		menus, _ := u.restaurantRepo.GetRestaurantMenus(context.Background(), rst.ID)

		ohs := restaurant.OpeningHours{}

		for _, opnHour := range opnHours {
			ohs.PushSchedule(opnHour.Day, opnHour.Start, opnHour.End)
		}

		mns := restaurant.Menus{}

		for _, menu := range menus {
			mns = append(mns, restaurant.Menu{
				DishName: menu.DishName,
				Price:    menu.Price,
			})
		}

		res = append(res, Restaurant{
			ID:           rst.ID,
			Name:         rst.Name,
			OpeningHours: ohs,
			CashBalance:  rst.CashBalance,
			Menu:         mns,
		})
	}

	return res
}

func (u Usecase) Process(userID int) error {
	return u.userRepo.ProcessPurchased(context.Background(), userID)
}
