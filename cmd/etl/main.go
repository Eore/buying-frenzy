package main

import (
	"encoding/json"
	"glints/internal/domain/restaurant"
	"glints/internal/domain/user"
	"glints/internal/util"
	"os"
)

func OpenFile(filePath string, model interface{}) error {
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		return err
	}

	return json.NewDecoder(file).Decode(&model)
}

func main() {
	restaurants := []restaurant.Restaurant{}
	OpenFile("restaurant_with_menu.json", &restaurants)

	resCSV := util.New("restaurants.csv")
	defer resCSV.Done()

	opnCSV := util.New("opening_hours.csv")
	defer opnCSV.Done()

	menuCSV := util.New("menus.csv")
	defer menuCSV.Done()

	for _, restaurant := range restaurants {
		restaurant, openingHours, menus := restaurant.ToSchema()
		resCSV.Write(restaurant)
		for _, openingHour := range openingHours {
			opnCSV.Write(openingHour)
		}

		for _, menu := range menus {
			menuCSV.Write(menu)
		}
	}

	users := []user.User{}
	if err := OpenFile("users_with_purchase_history.json", &users); err != nil {
		panic(err)
	}

	usrCSV := util.New("users.csv")
	defer usrCSV.Done()

	hstCSV := util.New("purchase_histories.csv")
	defer hstCSV.Done()

	for _, user := range users {
		usr, histories := user.ToSchema()
		usrCSV.Write(usr)
		for _, history := range histories {
			hstCSV.Write(history)
		}
	}
}
