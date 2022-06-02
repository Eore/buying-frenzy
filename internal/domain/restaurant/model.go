package restaurant

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Restaurant struct {
	Name         string             `json:"restaurantName,omitempty"`
	OpeningHours OpeningHoursString `json:"openingHours,omitempty"`
	CashBalance  float64            `json:"cashBalance,omitempty"`
	Menu         Menus              `json:"menu,omitempty"`
}

type Menu struct {
	DishName string  `json:"dishName,omitempty"`
	Price    float64 `json:"price,omitempty"`
}

type Menus []Menu

type OpeningHour struct {
	Day   Day       `json:"day"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type OpeningHours []*OpeningHour

type RestaurantSchema struct {
	ID          string
	Name        string
	CashBalance float64
}

type RestaurantSchemas []RestaurantSchema

type MenuSchema struct {
	ID           string
	RestaurantID string
	DishName     string
	Price        float64
}

type MenuSchemas []MenuSchema

type OpeningHourSchema struct {
	ID           string
	RestaurantID string
	Day          Day
	Start        time.Time
	End          time.Time
}

type OpeningHourSchemas []OpeningHourSchema

type Day time.Weekday

func (r Restaurant) ToSchema() (RestaurantSchema, []OpeningHourSchema, []MenuSchema) {
	idRes, _ := gonanoid.New()

	res := RestaurantSchema{
		ID:          idRes,
		Name:        r.Name,
		CashBalance: r.CashBalance,
	}

	openingHours := []OpeningHourSchema{}
	opns := r.OpeningHours.Parse()
	for _, opn := range opns {
		idOpn, _ := gonanoid.New()
		openingHours = append(openingHours, OpeningHourSchema{
			ID:           idOpn,
			RestaurantID: idRes,
			Day:          opn.Day,
			Start:        opn.Start,
			End:          opn.End,
		})
	}

	menus := []MenuSchema{}

	for _, menu := range r.Menu {
		idMenu, _ := gonanoid.New()
		menus = append(menus, MenuSchema{
			ID:           idMenu,
			RestaurantID: idRes,
			DishName:     menu.DishName,
			Price:        menu.Price,
		})
	}

	return res, openingHours, menus
}

func (r RestaurantSchema) ToCSVRow() []string {
	// column sequence :
	// id, name, cash_balance
	return []string{r.ID, r.Name, fmt.Sprintf("%0.2f", r.CashBalance)}
}

func (o OpeningHourSchema) ToCSVRow() []string {
	// column sequence :
	// id, restaurant_id, day, start, end
	return []string{o.ID, o.RestaurantID, o.Day.String(), o.Start.Format("15:04"), o.End.Format("15:04")}
}

func (m MenuSchema) ToCSVRow() []string {
	// column sequence :
	// id, restaurant_id, dish_name, price
	return []string{m.ID, m.RestaurantID, m.DishName, fmt.Sprintf("%0.2f", m.Price)}
}

type OpeningHoursString string

func (o OpeningHoursString) Parse() OpeningHours {
	scheds := strings.Split(string(o), " / ")

	opnHours := OpeningHours{}

	for _, sched := range scheds {
		days := o.getDays(sched)
		start, end := o.getHours(sched)

		for _, day := range days {
			opnHours.PushSchedule(day, start, end)
		}
	}

	return opnHours
}

func (o OpeningHoursString) getDays(str string) []Day {
	days := "(Mon|Tues|Weds|Thurs|Fri|Sat|Sun)"
	dayTokenMap := map[string]time.Weekday{
		"Mon":   time.Monday,
		"Tues":  time.Tuesday,
		"Weds":  time.Wednesday,
		"Thurs": time.Thursday,
		"Fri":   time.Friday,
		"Sat":   time.Saturday,
		"Sun":   time.Sunday,
	}

	{
		ls := []Day{}
		result := regexp.MustCompile(fmt.Sprintf(`%s - %s`, days, days)).FindString(str)
		if result != "" {
			s := strings.Split(result, " - ")
			start := dayTokenMap[s[0]]
			end := dayTokenMap[s[1]]

			for i := start; i <= end; i++ {
				ls = append(ls, Day(i))
			}

			return ls
		}
	}

	{
		result := regexp.MustCompile(fmt.Sprintf(`%s, %s`, days, days)).FindString(str)
		if result != "" {
			s := strings.Split(result, ", ")
			day1 := dayTokenMap[s[0]]
			day2 := dayTokenMap[s[1]]

			return []Day{Day(day1), Day(day2)}
		}
	}

	{
		result := regexp.MustCompile(fmt.Sprint(days)).FindString(str)
		day := dayTokenMap[result]

		return []Day{Day(day)}
	}
}

func (o OpeningHoursString) getHours(str string) (time.Time, time.Time) {
	parse := func(strTime string) time.Time {
		var t time.Time
		var err error

		t, err = time.Parse("3:04 pm", strTime)
		if err != nil {
			t, _ = time.Parse("3 pm", strTime)
		}

		return t
	}

	hours := regexp.MustCompile(`((\d{1,2}:\d{1,2}|\d{1,2}) (am|pm))`).FindAllString(str, -1)

	return parse(hours[0]), parse(hours[1])
}

func (d Day) String() string {
	return time.Weekday(d).String()
}

func (d Day) MarshalJSON() ([]byte, error) {
	str := time.Weekday(d).String()

	return json.Marshal(str)
}

func (d *Day) Scan(value interface{}) error {
	mapDay := map[string]Day{
		time.Sunday.String():    Day(time.Sunday),
		time.Monday.String():    Day(time.Monday),
		time.Tuesday.String():   Day(time.Tuesday),
		time.Wednesday.String(): Day(time.Wednesday),
		time.Thursday.String():  Day(time.Thursday),
		time.Friday.String():    Day(time.Friday),
		time.Saturday.String():  Day(time.Saturday),
	}

	str := value.(string)

	day, ok := mapDay[str]
	if !ok {
		return fmt.Errorf("%s not found", str)
	}

	*d = day

	return nil
}

func (o *OpeningHours) PushSchedule(day Day, start, end time.Time) {
	*o = append(*o, &OpeningHour{
		Day:   day,
		Start: start,
		End:   end,
	})
}

func (o *OpeningHours) Scan(value interface{}) error {
	str := value.(string)
	parsed := OpeningHoursString(str).Parse()
	o = &parsed

	return nil
}
