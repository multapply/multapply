package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

// Context keys
type contextKey string

func (c contextKey) String() string {
	return "pkg/middleware context key " + string(c)
}

const (
	contextKeyUserID = contextKey("uid")
	contextKeyRoles  = contextKey("roles")
)

// Authenticate - Middleware for requiring jwt token auth for a route
func Authenticate(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var tokenString string

		// Retrieve token from request header
		// Format: 'Authorization: Bearer <tokenstring>'
		// TODO: Just get token from cookie since they are sent automatically
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
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Token is invalid", 401)
		}

		// otherwise, valid and we set the context
		ctx := context.WithValue(r.Context(), contextKeyUserID, claims["uid"])
		ctx = context.WithValue(ctx, contextKeyRoles, claims["roles"])
		next(w, r.WithContext(ctx), ps)
	})
}

// Authorize - Middleware for making sure a user has access to this route
func Authorize(roles ...string) func(next httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			// TODO: Actually check permissions vs user roles from context
			log.Println(r.Context().Value(contextKeyUserID))
			next(w, r, ps)
		})
	}
}
