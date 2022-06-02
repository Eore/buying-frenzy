package api

import (
	"encoding/json"
	"fmt"
	"glints/internal"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type HTTP struct {
	usecase internal.Usecase
}

func NewHTTPAPI(usecase internal.Usecase) HTTP {
	return HTTP{usecase: usecase}
}

func (h HTTP) StartServer(port int) {
	router := httprouter.New()

	router.GET("/restaurants", h.RestaurantHandler)
	router.GET("/process/:user_id", h.ProcessHandler)

	portStr := fmt.Sprintf(":%d", port)
	log.Printf("server start on port %d\n", port)
	if err := http.ListenAndServe(portStr, router); err != nil {
		panic(err)
	}
}

func (h HTTP) RestaurantHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	urlQuery := r.URL.Query()
	day := urlQuery.Get("day")
	startTime := urlQuery.Get("start_time")
	endTime := urlQuery.Get("end_time")
	nDishes, _ := strconv.Atoi(urlQuery.Get("ndishes"))
	dishPriceStart, _ := strconv.ParseFloat(urlQuery.Get("dish_price_start"), 64)
	dishPriceEnd, _ := strconv.ParseFloat(urlQuery.Get("dish_price_end"), 64)

	filter := internal.RestaurantFilter{
		Day:            day,
		StartTime:      startTime,
		EndTime:        endTime,
		NDishes:        nDishes,
		DishPriceStart: dishPriceStart,
		DishPriceEnd:   dishPriceEnd,
	}

	restaurants := h.usecase.RestaurantList(filter)
	b, _ := json.Marshal(restaurants)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (h HTTP) ProcessHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	userIDStr := p.ByName("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if userIDStr == "" || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.usecase.Process(userID); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
