package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/multapply/multapply/util/constants"
	"github.com/multapply/multapply/util/userRoles"
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
			return
		}

		// otherwise, valid and we set the context
		ctx := context.WithValue(r.Context(), constants.ContextKeyUserID, claims["uid"])
		ctx = context.WithValue(ctx, constants.ContextKeyRoles, claims["roles"])
		next(w, r.WithContext(ctx), ps)
	})
}

// Authorize - Middleware for making sure a user has access to this route
func Authorize(roles ...string) func(next httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			// Check context for uid (user_id) and roles
			ctx := r.Context()
			userID := ctx.Value(constants.ContextKeyUserID)
			if userID == nil {
				http.Error(w, "Missing claims", 401)
				return
			}
			usrRoles := ctx.Value(constants.ContextKeyRoles)
			if usrRoles == nil {
				http.Error(w, "Missing claims", 401)
				return
			}

			// If user is admin, just continue
			// TODO: Check against DB because being admin is very big
			if userAllowed(usrRoles.(string), userRoles.Admin) {
				next(w, r, ps)
			}

			// Otherwise, we check the roles and user is forbidden
			if !userAllowed(usrRoles.(string), roles...) {
				http.Error(w, "Unauthorized", 403)
				return
			}

			// Otherwise, user is permitted to continue
			next(w, r, ps)
		})
	}
}

// userAllowed - Checks whether a given user's roles string passes the given requirements
// TODO: Test this func
func userAllowed(rolestring string, requirements ...string) bool {
	roles := strings.Split(rolestring, " ")
	reqs := make(map[string]string)

	for _, req := range requirements {
		reqs[req] = ""
	}
	for _, role := range roles {
		if _, ok := reqs[role]; ok {
			delete(reqs, role)
		}
	}

	return len(reqs) == 0
}

// ValidateRefreshToken - Middleware for checking if refresh token in cookie is valid
func ValidateRefreshToken(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Attempt to get refresh token from cookies in request
		mRT, err := r.Cookie("mRT")
		if err != nil {
			http.Error(w, "Refresh token missing", 400)
			return
		}
		refreshToken := mRT.Value

		// Parse refresh token
		parsedToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
			// Check to make sure token was signed with right method
			// TODO: Hide signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Token signature invalid")
			}
			// TODO: put token secrets somewhere else so we can return it here non-hardcoded
			return []byte("lmaoayy"), nil
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
			return
		}

		// Otherwise, set the context and continue
		ctx := context.WithValue(r.Context(), constants.ContextKeyTokenID, claims["tid"])

		next(w, r.WithContext(ctx), ps)
	})
}
