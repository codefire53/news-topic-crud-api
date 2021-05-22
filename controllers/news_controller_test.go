package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"news-topic-api/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mockServices "news-topic-api/mocks/services"
)

//getNewsRouter is a function that prepares a router to test the http routing
func getNewsRouter(newsController NewsController, requestType string) *mux.Router {
	router := mux.NewRouter()
	var pathSuffix string
	var method string
	controllerFunc := newsController.Create
	if requestType == "Create" {
		pathSuffix = "/"
		method = "POST"
		controllerFunc = newsController.Create
	} else if requestType == "Update" {
		pathSuffix = "/{id}"
		method = "PUT"
		controllerFunc = newsController.Update
	} else if requestType == "Delete" {
		pathSuffix = "/{id}"
		method = "DELETE"
		controllerFunc = newsController.Delete
	} else if requestType == "List" {
		pathSuffix = ""
		method = "GET"
		controllerFunc = newsController.List
	} else {
		pathSuffix = "/{id}"
		method = "GET"
		controllerFunc = newsController.GetDetail
	}
	router.HandleFunc(fmt.Sprintf("/news%s", pathSuffix), controllerFunc).Methods(method)
	return router
}

func createJSONRequestNews(method string, url string, data map[string]interface{}) *http.Request {
	jsonData, _ := json.Marshal(data)
	request, _ := http.NewRequest(method, url, bytes.NewReader(jsonData))
	request.Header.Add("Content-Type", "application/json")
	return request
}

func createURLParamRequestNews(method string, url string, data map[string]string) *http.Request {
	request, _ := http.NewRequest(method, url, nil)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	urlQuery := request.URL.Query()
	for key, value := range data {
		urlQuery.Add(key, value)
	}
	request.URL.RawQuery = urlQuery.Encode()
	return request
}

func createURLStandardRequestNews(method, url string) *http.Request {
	request, _ := http.NewRequest(method, url, nil)
	return request
}
func getMockReqTag() map[string]interface{} {
	tagData := make(map[string]interface{})
	tagData["name"] = "cryptocurrency"
	return tagData
}
func getMockReqNews() map[string]interface{} {
	newsData := make(map[string]interface{})
	newsData["title"] = "Harga bitcoin anjlok"
	newsData["thumbnail"] = "google.com/thumbnail.png"
	newsData["summary"] = "Harga bitcoin sempat menurun namun dogecoin justru naik"
	newsData["content"] = "Dikarenakan cuitan Elon Musk, nilai bitcoin sempat mengalami penurunan"
	newsData["tags"] = []map[string]interface{}{getMockReqTag()}
	newsData["topic"] = "bitcoin"
	newsData["status"] = "draft"
	return newsData
}
func getMockRequestParamsNews() map[string]string {
	queryParams := make(map[string]string)
	queryParams["status"] = "draft"
	queryParams["tag"] = "1"
	queryParams["topic"] = "bitcoin"
	return queryParams
}


func getMockNews() models.News {
	entity := models.News{
		Title: "Harga bitcoin anjlok",
		Thumbnail: "google.com/thumbnail.png",
		Summary: "Harga bitcoin sempat menurun namun dogecoin justru naik",
		Content: "Dikarenakan cuitan Elon Musk, nilai bitcoin sempat mengalami penurunan",
		Tags: []models.Tag{
			models.Tag{Model:gorm.Model{ID:1}},
		},
		Topic: "bitcoin",
		Status: "draft",
	}
	return entity
}

func getMockNewsList() []models.News {
	newsList := []models.News{getMockNews(), getMockNews()}
	return newsList
}


func TestCreateNewsFailedShouldReturnBadRequest(t *testing.T) {
	mockedRequestData := getMockReqNews()
	mockedNewsService := new(mockServices.INewsService)
	mockedNewsService.On("Create", mock.Anything).Return(models.News{}, fmt.Errorf("service can't create news"))
	newsController := InitNewsController(mockedNewsService)
	request := createJSONRequestNews("POST", "/news/", mockedRequestData)
	response := httptest.NewRecorder()
	router := getNewsRouter(newsController, "Create")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestCreateNewsInvalidRequestShouldReturnBadRequest(t *testing.T) {
	mockedRequestData := getMockReqNews()
	mockedRequestData["title"] = 12345
	mockedNewsService := new(mockServices.INewsService)
	newsController := InitNewsController(mockedNewsService)
	request := createJSONRequestNews("POST", "/news/", mockedRequestData)
	response := httptest.NewRecorder()
	router := getNewsRouter(newsController, "Create")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestCreateNewsSuccessShouldReturnCreated(t *testing.T) {
	mockedRequestData := getMockReqNews()
	mockedNewsService := new(mockServices.INewsService)
	mockedNewsService.On("Create", mock.Anything).Return(getMockNews(), nil)
	newsController := InitNewsController(mockedNewsService)
	request := createJSONRequestNews("POST", "/news/", mockedRequestData)
	response := httptest.NewRecorder()
	router := getNewsRouter(newsController, "Create")
	router.ServeHTTP(response, request)
	assert.Equal(t, 201, response.Code, "response code should be 201")
}

func TestUpdateNewsFailedShouldReturnBadRequest(t *testing.T) {
	mockedRequestData := getMockReqNews()
	mockedNewsService := new(mockServices.INewsService)
	mockedNewsService.On("Update", uint(1), mock.Anything).Return(models.News{}, errors.New("News failed to update"))
	newsController := InitNewsController(mockedNewsService)
	request := createJSONRequestNews("PUT", "/news/1", mockedRequestData)
	response := httptest.NewRecorder()
	router := getNewsRouter(newsController, "Update")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestUpdateNewsInvalidRequestDataShouldReturnBadRequest(t *testing.T) {
	mockedRequestData := getMockReqNews()
	mockedRequestData["title"] = 123
	mockedNewsService := new(mockServices.INewsService)
	newsController := InitNewsController(mockedNewsService)
	request := createJSONRequestNews("PUT", "/news/1", mockedRequestData)
	response := httptest.NewRecorder()
	router := getNewsRouter(newsController, "Update")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestUpdateNewslnvalidIDOnURLShouldReturnBadRequest(t *testing.T) {
	mockedRequestData := getMockReqNews()
	mockedNewsService := new(mockServices.INewsService)
	newsController := InitNewsController(mockedNewsService)
	request := createJSONRequestNews("PUT", "/news/asdasdasdasdasdasdasd", mockedRequestData)
	response := httptest.NewRecorder()
	router := getNewsRouter(newsController, "Update")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestUpdateNewsSuccessShouldReturnOk(t *testing.T) {
	mockedRequestData := getMockReqNews()
	mockedNewsService := new(mockServices.INewsService)
	mockedNewsService.On("Update", uint(1), mock.Anything).Return(getMockNews(), nil)
	newsController := InitNewsController(mockedNewsService)
	request := createJSONRequestNews("PUT", "/news/1", mockedRequestData)
	response := httptest.NewRecorder()
	router := getNewsRouter(newsController, "Update")
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "response code should be 200")
}

func TestDeleteNewsFailedShouldReturnBadRequest(t *testing.T) {
	mockedNewsService := new(mockServices.INewsService)
	mockedNewsService.On("Delete", uint(1)).Return(errors.New("News failed to delete"))
	newsController := InitNewsController(mockedNewsService)
	request := createURLStandardRequestNews("DELETE", "/news/1")
	response := httptest.NewRecorder()
	router := getNewsRouter(newsController, "Delete")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestDeleteNewsSuccessShouldReturnOk(t *testing.T) {
	mockedNewsService := new(mockServices.INewsService)
	mockedNewsService.On("Delete", uint(1)).Return(nil)
	newsController := InitNewsController(mockedNewsService)
	request := createURLStandardRequestNews("DELETE", "/news/1")
	response := httptest.NewRecorder()
	router := getNewsRouter(newsController, "Delete")
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "response code should be 200")
}

func TestDeleteNewslnvalidIDOnURLShouldReturnBadRequest(t *testing.T) {
	mockedNewsService := new(mockServices.INewsService)
	newsController := InitNewsController(mockedNewsService)
	request := createURLStandardRequestNews("DELETE", "/news/asdasdasdasdasdasdasd")
	response := httptest.NewRecorder()
	router := getNewsRouter(newsController, "Delete")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestListNewsSuccessShouldReturnOk(t *testing.T) {
	mockedServiceDataList := getMockNewsList()
	mockedNewsService := new(mockServices.INewsService)
	searchParams := getMockRequestParamsNews()
	mockedNewsService.On("List", searchParams).Return(models.NewsList{Data: mockedServiceDataList}, nil)
	newsController := InitNewsController(mockedNewsService)
	request := createURLParamRequestNews("GET", "/news", searchParams)
	response := httptest.NewRecorder()
	router := getNewsRouter(newsController, "List")
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "response code should be 200")
}

func TestListNewsFailedShouldReturnBadRequest(t *testing.T) {
	mockedNewsService := new(mockServices.INewsService)
	searchParams := getMockRequestParamsNews()
	mockedNewsService.On("List", searchParams).Return(models.NewsList{}, fmt.Errorf("Data empty"))
	newsController := InitNewsController(mockedNewsService)
	request := createURLParamRequestNews("GET", "/news", searchParams)
	response := httptest.NewRecorder()
	router := getNewsRouter(newsController, "List")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestGetDetailNewsSuccesShouldReturnOk(t *testing.T) {
	mockedNewsEntity := getMockNews()
	mockedNewsService := new(mockServices.INewsService)
	mockedNewsService.On("GetDetail", uint(1)).Return(mockedNewsEntity, nil)
	newsController := InitNewsController(mockedNewsService)
	request := createURLStandardRequestNews("GET", "/news/1")
	response := httptest.NewRecorder()
	router := getNewsRouter(newsController, "Detail")
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "response code should be 200")
}

func TestGetDetailNewsFailedShouldReturnBadRequest(t *testing.T) {
	mockedNewsService := new(mockServices.INewsService)
	mockedNewsService.On("GetDetail", uint(1)).Return(models.News{}, fmt.Errorf("Data not exists"))
	newsController := InitNewsController(mockedNewsService)
	request := createURLStandardRequestNews("GET", "/news/1")
	response := httptest.NewRecorder()
	router := getNewsRouter(newsController, "Detail")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestGetDetailNewslnvalidIDOnURLShouldReturnBadRequest(t *testing.T) {
	mockedNewsService := new(mockServices.INewsService)
	newsController := InitNewsController(mockedNewsService)
	request := createURLStandardRequestNews("GET", "/news/asdadasdasdasdasdadas")
	response := httptest.NewRecorder()
	router := getNewsRouter(newsController, "Detail")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}
