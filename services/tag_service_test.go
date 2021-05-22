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
	entity := models.Tag{
		Name: "crypto",
	}
	return entity
}
func getMockTagList() []models.Tag {
	tagList := []models.Tag{getMockTag()}
	return tagList
}
func getExpectedTagListOutput() models.TagsList {
	tagList := getMockTagList()
	response := models.TagsList{Data: tagList,}
	return response
}

func TestCreateTagSuccessReturnCreatedEntity(t *testing.T) {
	mockedTagRepository := new(mockRepositories.ITagRepository)
	mockTagEntity := getMockTag()
	mockedTagRepository.On("Create", mockTagEntity).Return(mockTagEntity, nil)
	tagService := InitTagService(mockedTagRepository)
	response, err  := tagService.Create(mockTagEntity)
	assert.Nil(t, err, "There should be no error")
	assert.True(t, reflect.DeepEqual(mockTagEntity, response) , "Response should be same as input")
}

func TestCreateTagFailedReturnError(t *testing.T) {
	mockedTagRepository := new(mockRepositories.ITagRepository)
	mockTagEntity := getMockTag()
	mockedTagRepository.On("Create", mockTagEntity).Return(models.Tag{}, fmt.Errorf("Tag creation failed"))
	tagService := InitTagService(mockedTagRepository)
	_, err  := tagService.Create(mockTagEntity)
	assert.NotNil(t, err, "There should be an error")
}

func TestUpdateTagSuccessReturnUpdatedEntity(t *testing.T) {
	mockedTagRepository := new(mockRepositories.ITagRepository)
	mockTagEntity := getMockTag()
	mockedTagRepository.On("Update", uint(1),  mockTagEntity).Return(mockTagEntity,nil)
	tagService := InitTagService(mockedTagRepository)
	response, err  := tagService.Update(uint(1), mockTagEntity)
	assert.Nil(t, err, "There should be no error")
	assert.True(t, reflect.DeepEqual(mockTagEntity, response) , "Response should be same as input")
}

func TestUpdateTagFailedReturnError(t *testing.T) {
	mockedTagRepository := new(mockRepositories.ITagRepository)
	mockTagEntity := getMockTag()
	mockedTagRepository.On("Update", uint(1),  mockTagEntity).Return(models.Tag{}, fmt.Errorf("Tag with specified id not found"))
	tagService := InitTagService(mockedTagRepository)
	_, err  := tagService.Update(uint(1), mockTagEntity)
	assert.NotNil(t, err, "There should be an error")
}

func TestDeleteTagSuccessReturnNoError(t *testing.T) {
	mockedTagRepository := new(mockRepositories.ITagRepository)
	mockedTagRepository.On("Delete", uint(1)).Return(nil)
	tagService := InitTagService(mockedTagRepository)
	err  := tagService.Delete(uint(1))
	assert.Nil(t, err, "There should be no error")

}

func TestDeleteTagFailedReturnError(t *testing.T) {
	mockedTagRepository := new(mockRepositories.ITagRepository)
	mockedTagRepository.On("Delete", uint(1)).Return(fmt.Errorf("Tag with specified id not found"))
	tagService := InitTagService(mockedTagRepository)
	err  := tagService.Delete(uint(1))
	assert.NotNil(t, err, "There should be an error")
}

func TestListTagSuccessReturnEntities(t *testing.T) {
	mockedTagRepository := new(mockRepositories.ITagRepository)
	mockTagEntities := getMockTagList()
	expectedOutput := getExpectedTagListOutput()
	tagService := InitTagService(mockedTagRepository)
	mockedTagRepository.On("List").Return(mockTagEntities, nil)
	response, err  := tagService.List()
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, len(expectedOutput.Data), len(response.Data), "Should return correct length")
	assert.True(t, reflect.DeepEqual(expectedOutput.Data, response.Data), "List should be the same")
}

func TestListTagFailedReturnError(t *testing.T) {
	mockedTagRepository := new(mockRepositories.ITagRepository)
	tagService := InitTagService(mockedTagRepository)
	mockedTagRepository.On("List").Return([]models.Tag{}, fmt.Errorf("Records not available"))
	_, err  := tagService.List()
	assert.NotNil(t, err, "There should be an error")
}
