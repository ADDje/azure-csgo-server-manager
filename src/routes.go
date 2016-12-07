package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	// API subrouter
	// Serves all JSON REST handlers prefixed with /api
	s := r.PathPrefix("/api").Subrouter()
	for _, route := range apiRoutes {
		s.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(AuthorizeHandler(route.HandlerFunc))
		//Handler(route.HandlerFunc)
	}

	// The login handler does not check for authentication.
	s.Path("/login").
		Methods("POST").
		Name("LoginUser").
		HandlerFunc(LoginUser)

	// Serves the frontend application from the app directory
	// Uses basic file server to serve index.html and Javascript application
	// Routes match the ones defined in React application
	r.Path("/login").
		Methods("GET").
		Name("Login").
		Handler(http.StripPrefix("/login", http.FileServer(http.Dir("./app/"))))
	r.Path("/settings").
		Methods("GET").
		Name("Settings").
		Handler(AuthorizeHandler(http.StripPrefix("/settings", http.FileServer(http.Dir("./app/")))))
	r.Path("/logs").
		Methods("GET").
		Name("Logs").
		Handler(AuthorizeHandler(http.StripPrefix("/logs", http.FileServer(http.Dir("./app/")))))
	r.Path("/configs").
		Methods("GET").
		Name("Configs").
		Handler(AuthorizeHandler(http.StripPrefix("/configs", http.FileServer(http.Dir("./app/")))))
	r.Path("/templates").
		Methods("GET").
		Name("Templates").
		Handler(AuthorizeHandler(http.StripPrefix("/templates", http.FileServer(http.Dir("./app/")))))
	r.Path("/config").
		Methods("GET").
		Name("Config").
		Handler(AuthorizeHandler(http.StripPrefix("/config", http.FileServer(http.Dir("./app/")))))
	r.Path("/server").
		Methods("GET").
		Name("Server").
		Handler(AuthorizeHandler(http.StripPrefix("/server", http.FileServer(http.Dir("./app/")))))
	r.PathPrefix("/").
		Methods("GET").
		Name("Index").
		Handler(http.FileServer(http.Dir("./app/")))

	return r
}

// Middleware returns a http.HandlerFunc which authenticates the users request
// Redirects user to login page if no session is found
func AuthorizeHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := Auth.aaa.Authorize(w, r, true); err != nil {
			log.Printf("Unauthenticated request %s %s %s", r.Method, r.Host, r.RequestURI)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// Defines all API REST endpoints
// All routes are prefixed with /api
var apiRoutes = Routes{
	Route{
		"GetServers",
		"GET",
		"/azure/servers/getall",
		GetAllServers,
	}, {
		"GetDefaultConfig",
		"GET",
		"/server/defaultconfig",
		GetDefaultServerConfig,
	}, {
		"LogoutUser",
		"GET",
		"/logout",
		LogoutUser,
	}, {
		"StatusUser",
		"GET",
		"/user/status",
		GetCurrentLogin,
	}, {
		"ListUsers",
		"GET",
		"/user/list",
		ListUsers,
	}, {
		"AddUser",
		"POST",
		"/user/add",
		AddUser,
	}, {
		"RemoveUser",
		"POST",
		"/user/remove",
		RemoveUser,
	}, {
		"GetServerSettings",
		"GET",
		"/settings",
		GetServerSettings,
	}, {
		"UpdateServerSettings",
		"POST",
		"/settings/update",
		UpdateServerSettings,
	},

	// Configs
	{
		"GetConfigs",
		"GET",
		"/configs/list",
		GetServerConfigs,
	},
	{
		"GetConfig",
		"GET",
		"/configs/get/{configName}",
		GetServerConfigByName,
	},
	{
		"GetConfigText",
		"GET",
		"/configs/gettext/{configName}",
		GetServerConfigTextByName,
	},
	{
		"UpdateConfigText",
		"POST",
		"/configs/gettext/{configName}",
		UpdateServerConfigText,
	},

	// Templates
	{
		"GetTemplates",
		"GET",
		"/templates/list",
		GetDeploymentTemplates,
	},
	{
		"UpdateParameters",
		"POST",
		"/templates/{templateName}parameters",
		UpdateTemplateParameters,
	},
	{
		"UpdateTemplate",
		"POST",
		"/templates/{templateName}/update",
		UpdateTemplateText,
	},
}
