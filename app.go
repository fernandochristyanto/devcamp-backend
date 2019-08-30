package main

import (
	"github.com/fernandochristyanto/devcamp-backend/internal"
	// kCache "github.com/koding/cache"
	"database/sql"
	"flag"
	"fmt"
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

func initHandler(handler *internal.Handler) error {

	// Initialize SQL DB
	db, err := sql.Open("postgres", "postgres://devcamp:devcamp@157.230.245.59:5432/?sslmode=disable")
	if err != nil {
		return err
	}
	handler.DB = db

	return nil
}

func initRouter(router *httprouter.Router, handler *internal.Handler) {

	router.GET("/", handler.Index)

	// Single user API
	router.POST("/users", handler.GetUserByEmailAndPassword)

	// `httprouter` library uses `ServeHTTP` method for it's 404 pages
	router.NotFound = handler
}

func main() {
	args := new(internal.Args)
	initFlags(args)

	handler := new(internal.Handler)
	if err := initHandler(handler); err != nil {
		panic(err)
	}

	router := httprouter.New()
	initRouter(router, handler)

	fmt.Printf("Apps served on :%d\n", args.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", args.Port), router))
}
