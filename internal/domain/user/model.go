package user

import (
	"fmt"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type User struct {
	ID              int64             `json:"id,omitempty"`
	Name            string            `json:"name,omitempty"`
	CashBalance     float64           `json:"cashBalance,omitempty"`
	PurchaseHistory []PurchaseHistory `json:"purchaseHistory,omitempty"`
}

type PurchaseHistory struct {
	DishName          string  `json:"dishName,omitempty"`
	RestaurantName    string  `json:"restaurantName,omitempty"`
	TransactionAmount float64 `json:"transactionAmount,omitempty"`
	TransactionDate   string  `json:"transactionDate,omitempty"`
}

type UserSchema struct {
	ID          int64
	Name        string
	CashBalance float64
}

type PurchaseHistorySchema struct {
	ID                string
	UserID            int64
	DishName          string
	RestaurantName    string
	TransactionAmount float64
	TransactionDate   time.Time
}

func (u User) ToSchema() (UserSchema, []PurchaseHistorySchema) {

	user := UserSchema{
		ID:          u.ID,
		Name:        u.Name,
		CashBalance: u.CashBalance,
	}

	histories := []PurchaseHistorySchema{}

	for _, history := range u.PurchaseHistory {
		id, _ := gonanoid.New()

		date, _ := time.Parse("01/02/2006 3:04 PM", history.TransactionDate)

		histories = append(histories, PurchaseHistorySchema{
			ID:                id,
			UserID:            u.ID,
			DishName:          history.DishName,
			RestaurantName:    history.RestaurantName,
			TransactionAmount: history.TransactionAmount,
			TransactionDate:   date,
		})
	}

	return user, histories
}

func (u UserSchema) ToCSVRow() []string {
	// column sequence :
	// id, name, cash_balance
	return []string{fmt.Sprint(u.ID), u.Name, fmt.Sprintf("%0.2f", u.CashBalance)}
}

func (p PurchaseHistorySchema) ToCSVRow() []string {
	// column sequence :
	// id, user_id, dish_name, restaurant_name, transaction_amount, transaction_date
	return []string{
		p.ID,
		fmt.Sprint(p.UserID),
		p.DishName,
		p.RestaurantName,
		fmt.Sprintf("%0.2f", p.TransactionAmount),
		p.TransactionDate.Format("01/02/2006 3:04 PM"),
	}
}
