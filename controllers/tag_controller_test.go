package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"news-topic-api/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mockServices "news-topic-api/mocks/services"
)

//getTagRouter is a function that prepares a router to test the http routing
func getTagRouter(tagController TagController, requestType string) *mux.Router {
	router := mux.NewRouter()
	var pathSuffix string
	var method string
	controllerFunc := tagController.Create
	if requestType == "Create" {
		pathSuffix = "/"
		method = "POST"
		controllerFunc = tagController.Create
	} else if requestType == "Update" {
		pathSuffix = "/{id}"
		method = "PUT"
		controllerFunc = tagController.Update
	} else if requestType == "Delete" {
		pathSuffix = "/{id}"
		method = "DELETE"
		controllerFunc = tagController.Delete
	} else {
		pathSuffix = ""
		method = "GET"
		controllerFunc = tagController.List
	}
	router.HandleFunc(fmt.Sprintf("/tag%s", pathSuffix), controllerFunc).Methods(method)
	return router
}

func createJSONRequestTag(method string, url string, data map[string]interface{}) *http.Request {
	jsonData, _ := json.Marshal(data)
	request, _ := http.NewRequest(method, url, bytes.NewReader(jsonData))
	request.Header.Add("Content-Type", "application/json")
	return request
}

func createURLStandardRequestTag(method, url string) *http.Request {
	request, _ := http.NewRequest(method, url, nil)
	return request
}
func getMockRequestTag() map[string]interface{} {
	tagData := make(map[string]interface{})
	tagData["name"] = "cryptocurrency"
	return tagData
}


func getMockTag() models.Tag {
	entity := models.Tag{
		Name: "cryptocurrency",
	}
	return entity
}

func getMockTagList() []models.Tag {
	tagList := []models.Tag{getMockTag()}
	return tagList
}

func TestCreateTagFailedShouldReturnBadRequest(t *testing.T) {
	mockedRequestData := getMockRequestTag()
	mockedTagService := new(mockServices.ITagService)
	mockedTagService.On("Create", mock.Anything).Return(models.Tag{}, fmt.Errorf("service can't create tag"))
	tagController := InitTagController(mockedTagService)
	request := createJSONRequestTag("POST", "/tag/", mockedRequestData)
	response := httptest.NewRecorder()
	router := getTagRouter(tagController, "Create")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestCreateTagInvalidRequestShouldReturnBadRequest(t *testing.T) {
	mockedRequestData := getMockRequestTag()
	mockedRequestData["name"] = 12345
	mockedTagService := new(mockServices.ITagService)
	tagController := InitTagController(mockedTagService)
	request := createJSONRequestTag("POST", "/tag/", mockedRequestData)
	response := httptest.NewRecorder()
	router := getTagRouter(tagController, "Create")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestCreateTagSuccessShouldReturnCreated(t *testing.T) {
	mockedRequestData := getMockRequestTag()
	mockedTagService := new(mockServices.ITagService)
	mockedTagService.On("Create", mock.Anything).Return(getMockTag(), nil)
	tagController := InitTagController(mockedTagService)
	request := createJSONRequestTag("POST", "/tag/", mockedRequestData)
	response := httptest.NewRecorder()
	router := getTagRouter(tagController, "Create")
	router.ServeHTTP(response, request)
	assert.Equal(t, 201, response.Code, "response code should be 201")
}

func TestUpdateTagFailedShouldReturnBadRequest(t *testing.T) {
	mockedRequestData := getMockRequestTag()
	mockedTagService := new(mockServices.ITagService)
	mockedTagService.On("Update", uint(1), mock.Anything).Return(models.Tag{}, errors.New("Tag failed to update"))
	tagController := InitTagController(mockedTagService)
	request := createJSONRequestTag("PUT", "/tag/1", mockedRequestData)
	response := httptest.NewRecorder()
	router := getTagRouter(tagController, "Update")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestUpdateTagInvalidRequestDataShouldReturnBadRequest(t *testing.T) {
	mockedRequestData := getMockRequestTag()
	mockedRequestData["name"] = 123
	mockedTagService := new(mockServices.ITagService)
	tagController := InitTagController(mockedTagService)
	request := createJSONRequestTag("PUT", "/tag/1", mockedRequestData)
	response := httptest.NewRecorder()
	router := getTagRouter(tagController, "Update")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestUpdateTaglnvalidIDOnURLShouldReturnBadRequest(t *testing.T) {
	mockedRequestData := getMockRequestTag()
	mockedTagService := new(mockServices.ITagService)
	tagController := InitTagController(mockedTagService)
	request := createJSONRequestTag("PUT", "/tag/asdasdasdasdasdasdasd", mockedRequestData)
	response := httptest.NewRecorder()
	router := getTagRouter(tagController, "Update")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestUpdateTagSuccessShouldReturnOk(t *testing.T) {
	mockedRequestData := getMockRequestTag()
	mockedTagService := new(mockServices.ITagService)
	mockedTagService.On("Update", uint(1), mock.Anything).Return(getMockTag(), nil)
	tagController := InitTagController(mockedTagService)
	request := createJSONRequestTag("PUT", "/tag/1", mockedRequestData)
	response := httptest.NewRecorder()
	router := getTagRouter(tagController, "Update")
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "response code should be 200")
}

func TestDeleteTagFailedShouldReturnBadRequest(t *testing.T) {
	mockedTagService := new(mockServices.ITagService)
	mockedTagService.On("Delete", uint(1)).Return(errors.New("Tag failed to delete"))
	tagController := InitTagController(mockedTagService)
	request := createURLStandardRequestTag("DELETE", "/tag/1")
	response := httptest.NewRecorder()
	router := getTagRouter(tagController, "Delete")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestDeleteTagSuccessShouldReturnOk(t *testing.T) {
	mockedTagService := new(mockServices.ITagService)
	mockedTagService.On("Delete", uint(1)).Return(nil)
	tagController := InitTagController(mockedTagService)
	request := createURLStandardRequestTag("DELETE", "/tag/1")
	response := httptest.NewRecorder()
	router := getTagRouter(tagController, "Delete")
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "response code should be 200")
}

func TestDeleteTaglnvalidIDOnURLShouldReturnBadRequest(t *testing.T) {
	mockedTagService := new(mockServices.ITagService)
	tagController := InitTagController(mockedTagService)
	request := createURLStandardRequestTag("DELETE", "/tag/asdasdasdasdasdasdasd")
	response := httptest.NewRecorder()
	router := getTagRouter(tagController, "Delete")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}

func TestListTagSuccessShouldReturnOk(t *testing.T) {
	mockedServiceDataList := getMockTagList()
	mockedTagService := new(mockServices.ITagService)
	mockedTagService.On("List").Return(models.TagsList{Data: mockedServiceDataList}, nil)
	tagController := InitTagController(mockedTagService)
	request := createURLStandardRequestTag("GET", "/tag")
	response := httptest.NewRecorder()
	router := getTagRouter(tagController, "List")
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "response code should be 200")
}

func TestListTagFailedShouldReturnBadRequest(t *testing.T) {
	mockedTagService := new(mockServices.ITagService)
	mockedTagService.On("List").Return(models.TagsList{}, fmt.Errorf("Data empty"))
	tagController := InitTagController(mockedTagService)
	request := createURLStandardRequestTag("GET", "/tag")
	response := httptest.NewRecorder()
	router := getTagRouter(tagController, "List")
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "response code should be 400")
}
