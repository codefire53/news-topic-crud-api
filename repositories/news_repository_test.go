package repositories

import (
	"errors"
	"fmt"
	"news-topic-api/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"news-topic-api/infrastructures"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	insertQueryNews     = "^INSERT INTO \"news\".+$"
	updateQueryNews      = `^UPDATE "news".*WHERE "id" = .*$`
	getQueryNews         = "^SELECT (.+) FROM \"news\".+$"
	deleteQueryNews = `^UPDATE "news".*WHERE "news"."id" = .*$`
	insertQueryTag = "^INSERT INTO \"tags\".+$"
	insertQueryTagNews = "^INSERT INTO \"news_tag\".+$"
)

func getMockNews() models.News {
	entity := models.News{
		Title: "Harga bitcoin anjlok",
		Thumbnail: "google.com/thumbnail.png",
		Summary: "Harga bitcoin sempat menurun namun dogecoin justru naik",
		Content: "Dikarenakan cuitan Elon Musk, nilai bitcoin sempat mengalami penurunan",
		Tags: []models.Tag{
			models.Tag{Name:"cryptocurrency"},
		},
		Topic: "bitcoin",
		Status: "draft",
	}
	return entity
}
func getMockNewsList() models.NewsList {
	newsList := []models.News{getMockNews(), getMockNews()}
	entities := models.NewsList{Data: newsList}
	return entities
}
func getMockListParamsNews() map[string]string {
	params := map[string]string {
		"topic": "bitcoin",
		"status": "draft",
		"tag": "1",
	}
	return params
}

func TestNewsCreateSuccess(t *testing.T) {
	testMock, assertion := setUpNews(t)
	testMock.ExpectQuery(insertQueryNews).WillReturnRows(sqlmock.NewRows([]string{"1", "1"})).WillReturnError(nil)
	testMock.ExpectQuery(insertQueryTag).WillReturnRows(sqlmock.NewRows([]string{"1", "1"})).WillReturnError(nil)
	testMock.ExpectExec(insertQueryTagNews).WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mockNews := getMockNews()
	newsRepo := new(NewsRepository)
	response, err := newsRepo.Create(mockNews)
	assertion.Nil(err, "Should be no error")
	assertion.NotNil(response, "Response should be not nil")
}

func TestNewsCreateFailed(t *testing.T) {
	testMock, assertion := setUpNews(t)
	testMock.ExpectQuery(insertQueryNews).WillReturnRows(sqlmock.NewRows([]string{"0", "1"})).WillReturnError(errors.New("Insertion error"))
	mockNews := getMockNews()
	newsRepo := new(NewsRepository)
	_, err := newsRepo.Create(mockNews)
	assertion.NotNil(err, "Should be an error")
	testMock.ExpectationsWereMet()
}

func TestNewsUpdateSuccess(t *testing.T) {
	testMock, assertion := setUpNews(t)
	returnRow := mockRowNews()
	mockUpdateData := getMockNews()
	testMock.ExpectQuery(getQueryNews).WillReturnRows(returnRow)
	testMock.ExpectExec(updateQueryNews).WithArgs(
		sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),
		sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	newsRepo := new(NewsRepository)
	_, err := newsRepo.Update(uint(1), mockUpdateData)
	assertion.Nil(err, "Should be no error")
	testMock.ExpectationsWereMet()
}

func TestNewsUpdateIDNotFoundReturnError(t *testing.T) {
	testMock, assertion := setUpNews(t)
	returnRow := mockRowNews()
	mockUpdateData := getMockNews()
	testMock.ExpectQuery(getQueryNews).WillReturnRows(returnRow).WillReturnError(fmt.Errorf("record not found"))
	newsRepo := new(NewsRepository)
	_, err := newsRepo.Update(uint(1), mockUpdateData)
	assertion.NotNil(err, "Should be an error")
	testMock.ExpectationsWereMet()
}


func TestNewsUpdateFailureReturnError(t *testing.T) {
	testMock, assertion := setUpNews(t)
	returnRow := mockRowNews()
	mockUpdateData := getMockNews()
	testMock.ExpectQuery(getQueryNews).WillReturnRows(returnRow)
	testMock.ExpectExec(updateQueryNews).WithArgs(
		sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),
		sqlmock.AnyArg(),sqlmock.AnyArg()).WillReturnError(fmt.Errorf("update error"))
	newsRepo := new(NewsRepository)
	_, err := newsRepo.Update(uint(1), mockUpdateData)
	assertion.NotNil(err, "Should be an error")
	testMock.ExpectationsWereMet()
}

func TestNewsDeleteSuccess(t *testing.T) {
	testMock, assertion := setUpNews(t)
	returnRow := mockRowNews()
	testMock.ExpectQuery(getQueryNews).WillReturnRows(returnRow)
	testMock.ExpectExec(deleteQueryNews).WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	newsRepo := new(NewsRepository)
	err := newsRepo.Delete(uint(1))
	assertion.Nil(err, "Should be no error")
	testMock.ExpectationsWereMet()
}

func TestNewsDeleteIDNotFoundReturnError(t *testing.T) {
	testMock, assertion := setUpNews(t)
	returnRow := mockRowNews()
	testMock.ExpectQuery(getQueryNews).WillReturnRows(returnRow).WillReturnError(fmt.Errorf("record not found"))
	newsRepo := new(NewsRepository)
	err := newsRepo.Delete(uint(1))
	assertion.NotNil(err, "Should be an error")
	testMock.ExpectationsWereMet()
}


func TestNewsDeleteFailureReturnError(t *testing.T) {
	testMock, assertion := setUpNews(t)
	returnRow := mockRowNews()
	testMock.ExpectQuery(getQueryNews).WillReturnRows(returnRow)
	testMock.ExpectExec(deleteQueryNews).WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg()).WillReturnError(fmt.Errorf("delete error"))
	newsRepo := new(NewsRepository)
	err := newsRepo.Delete(uint(1))
	assertion.NotNil(err, "Should be an error")
	testMock.ExpectationsWereMet()
}


func TestNewsGetByIDSuccess(t *testing.T) {
	testMock, assertion := setUpNews(t)
	returnRow := mockRowNews()
	testMock.ExpectQuery(getQueryNews).WillReturnRows(returnRow)
	newsRepo := new(NewsRepository)
	news, err := newsRepo.GetByID(uint(1))
	assertion.NotNil(news, "Entity is returned")
	assertion.Nil(err, "Should be no error")
}

func TestNewsGetByIDFailed(t *testing.T) {
	testMock, assertion := setUpNews(t)
	_ = mockRowNews()
	testMock.ExpectQuery(getQueryNews).WillReturnError(fmt.Errorf("record not found"))
	newsRepo := new(NewsRepository)
	_, err := newsRepo.GetByID(uint(1))
	assertion.NotNil(err, "There should be error")
}

func TestNewsListSuccess (t *testing.T) {
	testMock, assertion := setUpNews(t)
	returnRows := mockRowsNews()
	testMock.ExpectQuery(getQueryNews).WillReturnRows(returnRows)
	newsRepo := new(NewsRepository)
	searchParams := getMockListParamsNews()
	news, err := newsRepo.List(searchParams)
	assertion.Equal(len(news), 2)
	assertion.NotNil(news, "Entities are returned")
	assertion.Nil(err, "Should be no error")
}

func TestNewsListTagNotIntReturnError (t *testing.T) {
	testMock, assertion := setUpNews(t)
	returnRows := mockRowsNews()
	testMock.ExpectQuery(getQueryNews).WillReturnRows(returnRows)
	newsRepo := new(NewsRepository)
	searchParams := getMockListParamsNews()
	searchParams["tag"] = "asdadadasdadadasdasd"
	_, err := newsRepo.List(searchParams)
	assertion.NotNil(err, "Should be no error")
}

func TestNewsListFailureReturnError (t *testing.T) {
	testMock, assertion := setUpNews(t)
	_ = mockRowsNews()
	testMock.ExpectQuery(getQueryNews).WillReturnError(fmt.Errorf("rows not found"))
	newsRepo := new(NewsRepository)
	searchParams := getMockListParamsNews()
	_, err := newsRepo.List(searchParams)
	assertion.NotNil(err, "There should be an error")
}



func setUpNews(t *testing.T) (sqlmock.Sqlmock, *assert.Assertions) {
	mock := setUpMockNewsDB()
	assertions := assert.New(t)
	return mock, assertions
}

func setUpMockNewsDB() sqlmock.Sqlmock {
	mockDB, mock, _ := sqlmock.New()
	gormMockDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	infrastructures.SetDB(gormMockDB)
	return mock
}

func mockRowsNews() *sqlmock.Rows {
	newsFieldColumns := []string{"id","title","thumbnail","summary","content", "topic", "status"}
	rows := sqlmock.NewRows(newsFieldColumns)
	rows.AddRow("1", "Harga bitcoin anjlok", "google.com/thumbnail.png", "Harga bitcoin sempat menurun namun dogecoin justru naik", "Dikarenakan cuitan Elon Musk, nilai bitcoin sempat mengalami penurunan", "bitcoin", "draft")
	rows.AddRow("2", "Harga bitcoin anjlok", "google.com/thumbnail.png", "Harga bitcoin sempat menurun namun dogecoin justru naik", "Dikarenakan cuitan Elon Musk, nilai bitcoin sempat mengalami penurunan", "bitcoin", "draft")
	return rows
}

func mockRowNews() *sqlmock.Rows {
	newsFieldColumns := []string{"id","title","thumbnail","summary","content", "topic", "status"}
	rows := sqlmock.NewRows(newsFieldColumns)
	rows.AddRow("1", "Harga bitcoin anjlok", "google.com/thumbnail.png", "Harga bitcoin sempat menurun namun dogecoin justru naik", "Dikarenakan cuitan Elon Musk, nilai bitcoin sempat mengalami penurunan", "bitcoin", "draft")
	return rows
}
