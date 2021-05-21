package services

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	mockRepositories "news-topic-api/mocks/repositories"
	"news-topic-api/models"
	"reflect"
	"testing"
)

func getMockTag() models.Tag {
	entity := models.Tag {
		Name: "cryptocurrency",
	}
	return entity
}

func getMockTags() []models.Tag {
	entities := []models.Tag{getMockTag(), getMockTag()}
	return entities
}

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

func getMockNewsList() []models.News {
	newsList := []models.News{getMockNews(), getMockNews()}
	return newsList
}


func getExpectedNewsListOutput() models.NewsList {
	newsList := getMockNewsList()
	response := models.NewsList{Data: newsList,}
	return response
}
func getMockSearchParams() map[string]string {
	params := map[string]string {
		"status": "draft",
		"topic": "bitcoin",
		"tag": "1",
	}
	return params
}
func TestCreateNewsSuccessReturnCreatedEntity(t *testing.T) {
	mockedNewsRepository := new(mockRepositories.INewsRepository)
	mockNewsEntity := getMockNews()
	mockedNewsRepository.On("Create", mockNewsEntity).Return(mockNewsEntity, nil)
	newsService := InitNewsService(mockedNewsRepository)
	response, err  := newsService.Create(mockNewsEntity)
	assert.Nil(t, err, "There should be no error")
	assert.True(t, reflect.DeepEqual(mockNewsEntity, response) , "Response should be same as input")
}

func TestCreateNewsFailedReturnError(t *testing.T) {
	mockedNewsRepository := new(mockRepositories.INewsRepository)
	mockNewsEntity := getMockNews()
	mockedNewsRepository.On("Create", mockNewsEntity).Return(models.News{}, fmt.Errorf("News creation failed"))
	newsService := InitNewsService(mockedNewsRepository)
	_, err  := newsService.Create(mockNewsEntity)
	assert.NotNil(t, err, "There should be an error")
}

func TestUpdateNewsSuccessReturnUpdatedEntity(t *testing.T) {
	mockedNewsRepository := new(mockRepositories.INewsRepository)
	mockNewsEntity := getMockNews()
	mockedNewsRepository.On("Update", uint(1),  mockNewsEntity).Return(mockNewsEntity,nil)
	newsService := InitNewsService(mockedNewsRepository)
	response, err  := newsService.Update(uint(1), mockNewsEntity)
	assert.Nil(t, err, "There should be no error")
	assert.True(t, reflect.DeepEqual(mockNewsEntity, response) , "Response should be same as input")
}

func TestUpdateNewsFailedReturnError(t *testing.T) {
	mockedNewsRepository := new(mockRepositories.INewsRepository)
	mockNewsEntity := getMockNews()
	mockedNewsRepository.On("Update", uint(1),  mockNewsEntity).Return(models.News{}, fmt.Errorf("News with specified id not found"))
	newsService := InitNewsService(mockedNewsRepository)
	_, err  := newsService.Update(uint(1), mockNewsEntity)
	assert.NotNil(t, err, "There should be an error")
}

func TestDeleteNewsSuccessReturnNoError(t *testing.T) {
	mockedNewsRepository := new(mockRepositories.INewsRepository)
	mockedNewsRepository.On("Delete", uint(1)).Return(nil)
	newsService := InitNewsService(mockedNewsRepository)
	err  := newsService.Delete(uint(1))
	assert.Nil(t, err, "There should be no error")

}

func TestDeleteNewsFailedReturnError(t *testing.T) {
	mockedNewsRepository := new(mockRepositories.INewsRepository)
	mockedNewsRepository.On("Delete", uint(1)).Return(fmt.Errorf("News with specified id not found"))
	newsService := InitNewsService(mockedNewsRepository)
	err  := newsService.Delete(uint(1))
	assert.NotNil(t, err, "There should be an error")
}

func TestGetNewsDetailSuccessReturnCorrespondingEntity(t *testing.T) {
	mockedNewsRepository := new(mockRepositories.INewsRepository)
	mockNewsEntity := getMockNews()
	mockedNewsRepository.On("GetByID", uint(1)).Return(mockNewsEntity, nil)
	newsService := InitNewsService(mockedNewsRepository)
	response, err  := newsService.GetDetail(uint(1))
	assert.Nil(t, err, "There should be no error")
	assert.True(t, reflect.DeepEqual(mockNewsEntity, response) , "Response should be same as input")
}

func TestGetNewsDetailFailedReturnError(t *testing.T) {
	mockedNewsRepository := new(mockRepositories.INewsRepository)
	mockedNewsRepository.On("GetByID", uint(1)).Return(models.News{}, fmt.Errorf("News not found"))
	newsService := InitNewsService(mockedNewsRepository)
	_, err  := newsService.GetDetail(uint(1))
	assert.NotNil(t, err, "There should be an error")
}

func TestListNewsSuccessReturnEntities(t *testing.T) {
	mockedNewsRepository := new(mockRepositories.INewsRepository)
	mockNewsEntities := getMockNewsList()
	expectedOutput := getExpectedNewsListOutput()
	newsService := InitNewsService(mockedNewsRepository)

	searchParams := getMockSearchParams()
	mockedNewsRepository.On("List",searchParams).Return(mockNewsEntities, nil)
	response, err  := newsService.List(searchParams)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, len(expectedOutput.Data), len(response.Data), "Should return correct length")
	assert.True(t, reflect.DeepEqual(expectedOutput.Data, response.Data), "List should be the same")
}

func TestListNewsFailedReturnError(t *testing.T) {
	mockedNewsRepository := new(mockRepositories.INewsRepository)
	newsService := InitNewsService(mockedNewsRepository)
	searchParams := getMockSearchParams()
	mockedNewsRepository.On("List",searchParams).Return([]models.News{}, fmt.Errorf("Records not available"))
	_, err  := newsService.List(searchParams)
	assert.NotNil(t, err, "There should be an error")
}
