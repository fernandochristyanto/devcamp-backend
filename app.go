package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/fernandochristyanto/devcamp-backend/internal"
	"github.com/fernandochristyanto/devcamp-backend/internal/handler"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

const (
	booksTTL = 10
)

func initFlags(args *internal.Args) {
	port := flag.Int("port", 8000, "port number for your apps")
	args.Port = *port
}

func initHandler(handler *handler.Handler) error {

	// Initialize SQL DB
	db, err := sql.Open("postgres", "postgres://devcamp:devcamp@157.230.245.59:5432/?sslmode=disable")
	if err != nil {
		return err
	}
	handler.DB = db

	return nil
}

func initRouter(router *httprouter.Router, handler *handler.Handler) {

	router.GET("/", handler.Index)

	// Single user API
	router.POST("/users/shopregistration", handler.SellerRegistration)
	router.GET("/users/:id/products", handler.GetProductsByUser)

	router.GET("/products/garagesale", handler.GetGarageSales)
	router.GET("/products/detail/:id", handler.GetProductDetail)
	router.POST("/products/garagesale", handler.NewGarageSaleProduct)

	// `httprouter` library uses `ServeHTTP` method for it's 404 pages
	router.NotFound = handler
}

func main() {
	args := new(internal.Args)
	initFlags(args)

	handler := new(handler.Handler)
	if err := initHandler(handler); err != nil {
		panic(err)
	}

	router := httprouter.New()
	initRouter(router, handler)

	fmt.Printf("Apps served on :%d\n", args.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", args.Port), router))
}
