package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/multapply/multapply/handlers"
	"github.com/multapply/multapply/middleware"
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
	r.GET("/", middleware.Authenticate(env.GetUser))
	r.POST("/user/register", env.CreateUser)

	log.Print("Running server on " + PORT)
	log.Fatal(http.ListenAndServe(PORT, r))
}

// chain - chains middleware functions together
// compatible with both net/http's http.Handler and httprouter's httprouter.Handle
func chain(middlewares ...interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		for _, f := range middlewares {
			switch f := f.(type) {
			case httprouter.Handle:
				f(w, r, ps)
			case http.Handler:
				f.ServeHTTP(w, r)
			case func(http.ResponseWriter, *http.Request):
				f(w, r)
			default:
				http.Error(w, "Error", 500)
			}
		}
	}
}
