package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/multapply/multapply/handlers"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// TEMP
	//fmt.Println(time.Now().UTC().Format(time.RFC3339))

	// Connect to DB
	db, err := handlers.InitDB("postgres://jangerino:fkmcrul3z@localhost/multapply?sslmode=disable")
	if err != nil { // TODO: Use helper function for logging errors
		log.Panic(err)
	}
	env := &handlers.Env{DB: db}

	PORT := ":8080"
	r := httprouter.New()

	// Routes TODO: Logging middleware?
	r.GET("/user/:id", env.GetUser)
	r.POST("/user/register", env.CreateUser)

	log.Print("Running server on " + PORT)
	log.Fatal(http.ListenAndServe(PORT, r))
}
