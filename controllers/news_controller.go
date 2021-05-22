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
	resultData, err := n.newsService.Update(newsID, reqBody)
	if err != nil {
		helpers.ResponseError(res, http.StatusBadRequest, err)
		return
	}
	helpers.Response(res, http.StatusOK, resultData)
}

//List controller that handles list permohonan penyitaan request
func (p *PermohonanPenyitaanController) List(res http.ResponseWriter, req *http.Request) {
	role, pengajuID := p.getRoleAndUserIDFromHeader(req)
	if !helpers.IsAdminPidana(role) && !helpers.IsPimpinanKajari(role) && !helpers.IsJaksa(role) {
		helpers.ResponseBadRequest(res, http.StatusUnauthorized, constants.UnauthorizedRoleError())
		return
	}
	page, limit, err := p.parsePageAndLimit(req)
	if err != nil {
		helpers.ResponseBadRequest(res, http.StatusBadRequest, err)
		return
	}
	keyword := req.URL.Query().Get("keyword")
	searchParams := map[string]string{"keyword": keyword}
	var resultData models.DaftarPermohonanPenyitaan
	if helpers.IsJaksa(role) {
		searchParams["pengajuID"] = strconv.Itoa(pengajuID)
	}
	resultData, err = p.permohonanPenyitaanService.List(searchParams, page, limit)
	if err != nil {
		helpers.ResponseBadRequest(res, http.StatusBadRequest, err)
		return
	}
	helpers.Response(res, http.StatusOK, resultData)
}


//GetDetail controller that handles get permohonan penyitaan by id request
func(p *PermohonanPenyitaanController) GetDetail(res http.ResponseWriter, req *http.Request) {
	role, pengajuID := p.getRoleAndUserIDFromHeader(req)
	if !helpers.IsJaksa(role) && !helpers.IsPimpinanKajari(role) && !helpers.IsAdminPidana(role) {
		helpers.ResponseBadRequest(res, http.StatusUnauthorized, constants.UnauthorizedRoleError())
		return
	}
	if !helpers.IsJaksa(role) {
		pengajuID = -1
	}
	penyitaanID, err := p.parseID(req)
	if err != nil {
		helpers.ResponseBadRequest(res, http.StatusBadRequest, err)
		return
	}
	resultData, err := p.permohonanPenyitaanService.GetDetail(penyitaanID, pengajuID)
	if err != nil {
		helpers.ResponseBadRequest(res, http.StatusBadRequest, err)
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


