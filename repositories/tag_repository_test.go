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
	insertQueryTags     = "^INSERT INTO \"tags\".+$"
	updateQueryTags      = `^UPDATE "tags".*WHERE "id" = .*$`
	getQueryTags         = "^SELECT (.+) FROM \"tags\".+$"
	deleteQueryTags = `^UPDATE "news".*WHERE "tags"."id" = .*$`

)
func getMockTag() models.Tag {
	entity := models.Tag{
		Name: "crypto",
	}
	return entity
}

func getMockTagList() models.TagsList {
	tagsList := []models.Tag{getMockTag()}
	entities := models.TagsList{Data: tagsList}
	return entities
}

func TestTagCreateSuccess(t *testing.T) {
	testMock, assertion := setUpTag(t)
	testMock.ExpectQuery(insertQueryTags).WillReturnRows(sqlmock.NewRows([]string{"1", "1"})).WillReturnError(nil)
	mockTag := getMockTag()
	tagRepo := new(TagRepository)
	response, err := tagRepo.Create(mockTag)
	assertion.Nil(err, "Should be no error")
	assertion.NotNil(response, "Response should be not nil")
}

func TestTagCreateFailed(t *testing.T) {
	testMock, assertion := setUpTag(t)
	testMock.ExpectQuery(insertQueryTags).WillReturnRows(sqlmock.NewRows([]string{"0", "1"})).WillReturnError(errors.New("Insertion error"))
	mockTag := getMockTag()
	tagRepo := new(TagRepository)
	_, err := tagRepo.Create(mockTag)
	assertion.NotNil(err, "Should be an error")
	testMock.ExpectationsWereMet()
}

func TestTagUpdateSuccess(t *testing.T) {
	testMock, assertion := setUpTag(t)
	returnRow := mockRowTag()
	mockUpdateData := getMockTag()
	testMock.ExpectQuery(getQueryTags).WillReturnRows(returnRow)
	testMock.ExpectExec(updateQueryTags).WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	tagRepo := new(TagRepository)
	_, err := tagRepo.Update(uint(1), mockUpdateData)
	assertion.Nil(err, "Should be no error")
	testMock.ExpectationsWereMet()
}

func TestTagUpdateIDNotFoundReturnError(t *testing.T) {
	testMock, assertion := setUpTag(t)
	returnRow := mockRowTag()
	mockUpdateData := getMockTag()
	testMock.ExpectQuery(getQueryTags).WillReturnRows(returnRow).WillReturnError(fmt.Errorf("record not found"))
	tagRepo := new(TagRepository)
	_, err := tagRepo.Update(uint(1), mockUpdateData)
	assertion.NotNil(err, "Should be an error")
	testMock.ExpectationsWereMet()
}


func TestTagUpdateFailureReturnError(t *testing.T) {
	testMock, assertion := setUpTag(t)
	returnRow := mockRowTag()
	mockUpdateData := getMockTag()
	testMock.ExpectQuery(getQueryTags).WillReturnRows(returnRow)
	testMock.ExpectExec(updateQueryTags).WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg()).WillReturnError(fmt.Errorf("update error"))
	tagRepo := new(TagRepository)
	_, err := tagRepo.Update(uint(1), mockUpdateData)
	assertion.NotNil(err, "Should be an error")
	testMock.ExpectationsWereMet()
}

func TestTagDeleteSuccess(t *testing.T) {
	testMock, assertion := setUpTag(t)
	returnRow := mockRowTag()
	testMock.ExpectQuery(getQueryTags).WillReturnRows(returnRow)
	testMock.ExpectExec(deleteQueryTags).WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	tagRepo := new(TagRepository)
	err := tagRepo.Delete(uint(1))
	assertion.Nil(err, "Should be no error")
	testMock.ExpectationsWereMet()
}

func TestTagDeleteIDNotFoundReturnError(t *testing.T) {
	testMock, assertion := setUpTag(t)
	returnRow := mockRowTag()
	testMock.ExpectQuery(getQueryTags).WillReturnRows(returnRow).WillReturnError(fmt.Errorf("record not found"))
	tagRepo := new(TagRepository)
	err := tagRepo.Delete(uint(1))
	assertion.NotNil(err, "Should be an error")
	testMock.ExpectationsWereMet()
}


func TestTagDeleteFailureReturnError(t *testing.T) {
	testMock, assertion := setUpTag(t)
	returnRow := mockRowTag()
	testMock.ExpectQuery(getQueryTags).WillReturnRows(returnRow)
	testMock.ExpectExec(deleteQueryTags).WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg()).WillReturnError(fmt.Errorf("delete error"))
	tagRepo := new(TagRepository)
	err := tagRepo.Delete(uint(1))
	assertion.NotNil(err, "Should be an error")
	testMock.ExpectationsWereMet()
}

func TestTagListSuccess (t *testing.T) {
	testMock, assertion := setUpTag(t)
	returnRow := mockRowTag()
	testMock.ExpectQuery(getQueryTags).WillReturnRows(returnRow)
	tagRepo := new(TagRepository)
	tags, err := tagRepo.List()
	assertion.Equal(len(tags), 1)
	assertion.NotNil(tags, "Entities are returned")
	assertion.Nil(err, "Should be no error")
}

func TestTagListFailureReturnError (t *testing.T) {
	testMock, assertion := setUpTag(t)
	_ = mockRowsTag()
	testMock.ExpectQuery(getQueryTags).WillReturnError(fmt.Errorf("rows not found"))
	tagRepo := new(TagRepository)
	_, err := tagRepo.List()
	assertion.NotNil(err, "There should be an error")
}



func setUpTag(t *testing.T) (sqlmock.Sqlmock, *assert.Assertions) {
	mock := setUpMockTagDB()
	assertions := assert.New(t)
	return mock, assertions
}

func setUpMockTagDB() sqlmock.Sqlmock {
	mockDB, mock, _ := sqlmock.New()
	gormMockDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: mockDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	infrastructures.SetDB(gormMockDB)
	return mock
}


func mockRowTag() *sqlmock.Rows {
	tagFieldColumns := []string{"id","name"}
	rows := sqlmock.NewRows(tagFieldColumns)
	rows.AddRow("1", "cryptocurrency")
	return rows
}