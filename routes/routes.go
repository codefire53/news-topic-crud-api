package routes

import (
	"github.com/gorilla/mux"
)

// Route ...
type Route struct{}

// Init is the initiator for the route location
func (r *Route) Init() *mux.Router {
	// init routes
	router := mux.NewRouter().StrictSlash(false)
	return router
}
