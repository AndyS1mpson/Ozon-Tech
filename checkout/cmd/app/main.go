// Project launch
package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/products"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/checkout/internal/handlers/listcart"
	"route256/checkout/internal/handlers/purchase"
	"route256/libs/srvwrapper"
)

const port = ":8080"

// Service start point
func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatalln("ERR: ", err)
	}

	service := domain.New(loms.New(cfg.Services.Loms), products.New(cfg.Services.Products, cfg.Token))

	addToCartHandler := &addtocart.Handler{
		Service: service,
	}
	http.Handle("/addToCart", srvwrapper.New(addToCartHandler.Handle))

	deleteFromCartHandler := &deletefromcart.Handler{
		Service: service,
	}
	http.Handle("/deleteFromCart", srvwrapper.New(deleteFromCartHandler.Handle))

	listCartHandler := &listcart.Handler{Service: service}
	http.Handle("/listCart", srvwrapper.New(listCartHandler.Handle))

	purchaseHandler := &purchase.Handler{
		Service: service,
	}
	http.Handle("/purchase", srvwrapper.New(purchaseHandler.Handle))

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("Can not run server: ", err)
	}

}
