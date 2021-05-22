package routes

import (
	"github.com/gorilla/mux"
	"news-topic-api/controllers"
	"news-topic-api/repositories"
	"news-topic-api/services"
)

// Route ...
type Route struct{}

// Init is the initiator for the route location
func (r *Route) Init() *mux.Router {
	// init repositories
	newsRepository := new(repositories.NewsRepository)
	tagRepository := new(repositories.TagRepository)

	// init services
	newsService := services.InitNewsService(newsRepository)
	tagService := services.InitTagService(tagRepository)

	// init Controllers
	newsController := controllers.InitNewsController(newsService)
	tagController := controllers.InitTagController(tagService)

	// init routes
	router := mux.NewRouter().StrictSlash(false)
	news := router.PathPrefix("/news").Subrouter()
	tag := router.PathPrefix("/tag").Subrouter()

	//news endpoint
	news.HandleFunc("/", newsController.Create).Methods("POST")
	news.HandleFunc("/{id}", newsController.Update).Methods("PUT")
	news.HandleFunc("/{id}", newsController.Delete).Methods("DELETE")
	news.HandleFunc("/{id}", newsController.GetDetail).Methods("GET")
	news.HandleFunc("", newsController.List).Methods("GET")

	//tag endpoint
	tag.HandleFunc("/", tagController.Create).Methods("POST")
	tag.HandleFunc("/{id}", tagController.Update).Methods("PUT")
	tag.HandleFunc("/{id}", tagController.Delete).Methods("DELETE")
	tag.HandleFunc("", tagController.List).Methods("GET")

	return router
}
