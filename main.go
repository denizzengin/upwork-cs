package main

import (
	b64 "encoding/base64"

	"context"
	"net/http"
	"strings"

	"github.com/denizzengin/upwork-cs/internal/api"
	"github.com/denizzengin/upwork-cs/internal/bll"
	"github.com/denizzengin/upwork-cs/internal/storage"
	"github.com/denizzengin/upwork-cs/model"
	"github.com/denizzengin/upwork-cs/pkg/config"
	"github.com/denizzengin/upwork-cs/pkg/logger"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	db, err := storage.GetDBConnection(config.Config.Database)
	if err != nil {
		logger.Log.Error("connection error")
	}

	s := storage.New(db)

	bll.Migrate(s)

	router := mux.NewRouter()
	apiRoutes := mux.NewRouter().PathPrefix("/api").Subrouter()
	managementRoutes := mux.NewRouter().PathPrefix("/management").Subrouter()

	//apiRoutes.HandleFunc("/login", api)
	managementRoutes.HandleFunc("/register", api.CreateUser(s)).Methods(http.MethodPost)
	managementRoutes.HandleFunc("/login", api.Login(s)).Methods(http.MethodPost)
	apiRoutes.HandleFunc("/users", api.FindAll(s)).Methods(http.MethodGet)
	apiRoutes.HandleFunc("/users/me", api.FindUserByName(s)).Methods(http.MethodGet)
	apiRoutes.HandleFunc("/users/{id:[0-9]+}", api.FindUserByID(s)).Methods(http.MethodGet)
	apiRoutes.HandleFunc("/users/admin/{id:[0-9]+}", api.UpdateUser(s)).Methods(http.MethodPut)

	common := negroni.New()
	router.PathPrefix("/api").Handler(common.With(negroni.HandlerFunc(checkToken), negroni.HandlerFunc(checkRoleRight), negroni.Wrap(apiRoutes)))
	router.PathPrefix("/management").Handler(managementRoutes)
	common.UseHandler(router)
	logger.Log.Info("running server...")
	if err := http.ListenAndServe(":8080", common); err != nil {
		logger.Log.Error(err.Error())
	}

}

var jwtKey = []byte("my_secret_key")

func checkRoleRight(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := r.Context()
	user, ok := ctx.Value("user").(*model.Claims)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	db, err := storage.GetDBConnection(config.Config.Database)
	if err != nil {
		logger.Log.Error("connection error")
	}

	store := storage.New(db)

	dbUser, err := store.Users().FindByName(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dbUserRoles, err := store.Users().FindUserRoles(dbUser.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if dbUserRoles == nil && len(*dbUserRoles) <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	hasRoleRight := false
	for _, role := range *dbUserRoles {
		// admin role
		if role.RoleId == 1 {
			hasRoleRight = true
			break
		}

		if store.Users().CheckRoleRight(role.RoleId, r.URL.Path) {
			hasRoleRight = true
			break
		}
	}

	if !hasRoleRight {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	next(w, r)
}

func checkToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//todo check jwt base authentication
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	bearerToken = strings.ReplaceAll(bearerToken, "Bearer", "")
	bearerToken = strings.TrimPrefix(bearerToken, " ")
	decodedToken, _ := b64.StdEncoding.DecodeString(bearerToken)

	decodedTokenStr := string(decodedToken)
	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(decodedTokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, "user", claims)
	next(w, r.WithContext(ctx))
}
