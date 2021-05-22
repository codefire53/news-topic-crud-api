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

//TagController ...
type TagController struct {
	tagService services.ITagService
}

//InitTagController initializes tag controller given the tag service
func InitTagController(tagService services.ITagService) TagController {
	tagController := new(TagController)
	tagController.tagService = tagService
	return *tagController
}

//Create controller that handles create tag request
func (t *TagController) Create(res http.ResponseWriter, req *http.Request) {
	reqBody, err := t.decodeRequest(req)
	if err != nil {
		fmt.Println(err)
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	resultData, err := t.tagService.Create(reqBody)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	helpers.Response(res, http.StatusCreated, resultData)
}

//Update controller that handles update tag request
func (t *TagController) Update(res http.ResponseWriter, req *http.Request) {
	reqBody, err := t.decodeRequest(req)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	tagID, err := t.parseID(req)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	resultData, err := t.tagService.Update(tagID, reqBody)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	helpers.Response(res, http.StatusOK, resultData)
}

//Delete controller that handles delete tag request
func (t *TagController) Delete(res http.ResponseWriter, req *http.Request) {
	tagID, err := t.parseID(req)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	err = t.tagService.Delete(tagID)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	helpers.Response(res, http.StatusOK, nil)
}

//List controller that handles list tag request
func (t *TagController) List(res http.ResponseWriter, req *http.Request) {
	var resultData models.TagsList
	resultData, err := t.tagService.List()
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	helpers.Response(res, http.StatusOK, resultData)
}

func (t *TagController) decodeRequest(req *http.Request) (models.Tag, error) {
	reqContent := models.Tag{}
	if err := json.NewDecoder(req.Body).Decode(&reqContent); err != nil {
		return models.Tag{}, err
	}
	return reqContent, nil
}

func (t *TagController) parseID(req *http.Request) (uint, error) {
	id, err := strconv.ParseUint(mux.Vars(req)["id"], 10, 64)
	if err != nil {
		return uint(0), fmt.Errorf("invalid format for id")
	}
	return uint(id), nil
}