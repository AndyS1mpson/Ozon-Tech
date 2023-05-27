// Project launch
package main

import (
	"log"
	"net/http"
	"route256/libs/srvwrapper"
	"route256/loms/internal/domain"
	"route256/loms/internal/handlers/cancelorder"
	"route256/loms/internal/handlers/createorder"
	"route256/loms/internal/handlers/listorder"
	"route256/loms/internal/handlers/orderpayed"
	"route256/loms/internal/handlers/stocks"
)

const port = ":8081"

// Service start point
func main() {
	stocksHandler := &stocks.Handler{Service: domain.New()}
	http.Handle("/stocks", srvwrapper.New(stocksHandler.Handle))

	createGoodHandler := &createorder.Handler{Service: domain.New()}
	http.Handle("/createOrder", srvwrapper.New(createGoodHandler.Handle))

	listOrderHandler := &listorder.Handler{Service: domain.New()}
	http.Handle("/listOrder", srvwrapper.New(listOrderHandler.Handle))

	orderPayed := &orderpayed.Handler{Service: domain.New()}
	http.Handle("/orderPayed", srvwrapper.New(orderPayed.Handle))

	cancelOrder := &cancelorder.Handler{Service: domain.New()}
	http.Handle("/cancelOrder", srvwrapper.New(cancelOrder.Handle))

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("ERR: ", err)
	}
}
