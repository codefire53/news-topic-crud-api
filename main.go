package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"news-topic-api/infrastructures"
	"os"
	"os/signal"
	"syscall"
	"time"
	"news-topic-api/routes"
)

func main() {
	godotenv.Load()

	var rt routes.Route

	r := rt.Init()
	infrastructures.InitDB()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("Wait for 2 second to finish processing")
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()
	headersOK := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOK := handlers.AllowedOrigins([]string{"*"})
	methodsOK := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"})

	port := "8080"
	fmt.Println("Server served at port " + port)

	if err := http.ListenAndServe(":"+port, handlers.CORS(originsOK, headersOK, methodsOK)(r)); err != nil {
		log.Fatal("Unable to start service: " + err.Error())
	}

}
