package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/multapply/multapply/handlers"
	mw "github.com/multapply/multapply/middleware"
	"github.com/multapply/multapply/util/userRoles"
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
	// TODO: Rename env to API or something, 'env' is weird
	env := &handlers.Env{DB: db}

	PORT := ":8080"
	r := httprouter.New()

	// Routes TODO: Logging middleware?
	// TODO: Api prefix versioning, i.e. "api/v1/..."
	r.GET("/", chain(env.GetUser,
		mw.Authenticate,
		mw.Authorize(userRoles.BasicUser)))
	r.POST("/user/register", env.CreateUser)
	r.POST("/user/login", env.LoginUser)

	r.GET("/auth/token", chain(env.RequestNewToken,
		mw.ValidateRefreshToken))

	log.Print("Running server on " + PORT)
	log.Fatal(http.ListenAndServe(PORT, r))
}

// chain - chains middleware functions together
// compatible with both net/http's http.Handler and httprouter's httprouter.Handle
func chain(endpoint func(http.ResponseWriter, *http.Request, httprouter.Params),
	middleware ...func(httprouter.Handle) httprouter.Handle) httprouter.Handle {

	for i := len(middleware) - 1; i >= 0; i-- {
		mw := middleware[i]
		endpoint = mw(endpoint)
	}
	return httprouter.Handle(endpoint)
}
