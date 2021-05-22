package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"news-topic-api/helpers"
	"news-topic-api/models"
	"news-topic-api/services"
	"net/http"
	"strconv"
)

//NewsController ...
type NewsController struct {
	newsService services.INewsService
}

//InitNewsController initializes news controller given the news service
func InitNewsController(newsService services.INewsService) NewsController {
	newsController := new(NewsController)
	newsController.newsService = newsService
	return *newsController
}

//Create controller that handles create news request
func (n *NewsController) Create(res http.ResponseWriter, req *http.Request) {
	reqBody, err := n.decodeRequest(req)
	if err != nil {
		fmt.Println(err)
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	resultData, err := n.newsService.Create(reqBody)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	helpers.Response(res, http.StatusCreated, resultData)
}

//Update controller that handles update news request
func (n *NewsController) Update(res http.ResponseWriter, req *http.Request) {
	reqBody, err := n.decodeRequest(req)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	newsID, err := n.parseID(req)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	resultData, err := n.newsService.Update(newsID, reqBody)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	helpers.Response(res, http.StatusOK, resultData)
}

//Delete controller that handles delete news request
func (n *NewsController) Delete(res http.ResponseWriter, req *http.Request) {
	newsID, err := n.parseID(req)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	err = n.newsService.Delete(newsID)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	helpers.Response(res, http.StatusOK, nil)
}

//List controller that handles list news request
func (n *NewsController) List(res http.ResponseWriter, req *http.Request) {
	searchParams := n.parseParams(req)
	var resultData models.NewsList
	resultData, err := n.newsService.List(searchParams)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	helpers.Response(res, http.StatusOK, resultData)
}


//GetDetail controller that handles get news by id request
func(n *NewsController) GetDetail(res http.ResponseWriter, req *http.Request) {
	newsID, err := n.parseID(req)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	resultData, err := n.newsService.GetDetail(newsID)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	helpers.Response(res, http.StatusOK, resultData)
}


func (n *NewsController) decodeRequest(req *http.Request) (models.News, error) {
	reqContent := models.News{}
	if err := json.NewDecoder(req.Body).Decode(&reqContent); err != nil {
		return models.News{}, err
	}
	return reqContent, nil
}

func (n *NewsController) parseID(req *http.Request) (uint, error) {
	id, err := strconv.ParseUint(mux.Vars(req)["id"], 10, 64)
	if err != nil {
		return uint(0), fmt.Errorf("invalid format for id")
	}
	return uint(id), nil
}

func (n *NewsController) parseParams(req *http.Request) map[string]string {
	topic := req.URL.Query().Get("topic")
	tag := req.URL.Query().Get("tag")
	status := req.URL.Query().Get("status")
	searchParams := make(map[string]string)
	if topic != "" {
		searchParams["topic"] = topic
	}
	if tag != "" {
		searchParams["tag"] = tag
	}
	if status != "" {
		searchParams["status"] = status
	}
	return searchParams
}