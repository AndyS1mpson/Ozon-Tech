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
	service := domain.New()
	stocksHandler := stocks.New(service)
	http.Handle("/stocks", srvwrapper.New(stocksHandler.Handle))

	createGoodHandler := createorder.New(service)
	http.Handle("/createOrder", srvwrapper.New(createGoodHandler.Handle))

	listOrderHandler := listorder.New(service)
	http.Handle("/listOrder", srvwrapper.New(listOrderHandler.Handle))

	orderPayed := orderpayed.New(service)
	http.Handle("/orderPayed", srvwrapper.New(orderPayed.Handle))

	cancelOrder := cancelorder.New(service)
	http.Handle("/cancelOrder", srvwrapper.New(cancelOrder.Handle))

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("ERR: ", err)
	}
}
