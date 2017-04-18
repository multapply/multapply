package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

// Authenticate - Middleware for requiring jwt token auth for a route
func Authenticate(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var tokenString string

		// Retrieve token from request header
		// Format: 'Authorization: Bearer <tokenstring>'
		tokens, ok := r.Header["Authorization"]
		if ok && len(tokens) >= 1 {
			tokenString = strings.TrimPrefix(tokens[0], "Bearer ")
		} else {
			http.Error(w, "No Authorization header set", 400)
			return
		}
		// if we have the header but no token
		if tokenString == "" {
			http.Error(w, "Auth token missing", 400)
			return
		}

		// otherwise we parse the token
		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check to make sure token was signed with right method
			// TODO: Hide signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Token signature invalid")
			}
			// TODO: put token secrets somewhere else so we can return it here non-hardcoded
			return []byte("ayylmao"), nil
		})
		if err != nil {
			http.Error(w, "Error parsing token", 401)
			return
		}

		// token is invalid
		if parsedToken == nil || !parsedToken.Valid {
			http.Error(w, "Token is invalid", 401)
			return
		}

		// Extract and verify claims
		err = parsedToken.Claims.Valid()
		if err != nil {
			http.Error(w, "Token is invalid", 401)
			return
		}

		// otherwise, valid and we set the context
		// ctx := context.WithValue(r.Context())
		next(w, r, ps)
	})
}

// Authorize - Middleware for making sure a user has access to this route
func Authorize(roles ...string) func(next httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		})
	}
}
