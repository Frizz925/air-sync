package handlers

import "github.com/gorilla/mux"

type RouteHandler interface {
	RegisterRoutes(r *mux.Router)
}
