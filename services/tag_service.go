package services

import (
	"news-topic-api/models"
	"news-topic-api/repositories"
)


//ITagService interface for tag service
type ITagService interface {
	Create(tag models.Tag) (models.Tag, error)
	Update(tagID uint,  tag models.Tag) (models.Tag, error)
	Delete(tagID uint) (error)
	List() (models.TagsList, error)
}

//TagService ...
type TagService struct {
	tagRepository repositories.ITagRepository
}

//InitTagService initialize a tag service instance with specific tag repository
func InitTagService(tagRepository repositories.ITagRepository) ITagService {
	tagService := new(TagService)
	tagService.tagRepository = tagRepository
	return tagService
}

//Create ...
func (t TagService) Create(tag models.Tag) (models.Tag, error) {
	instance, err := t.tagRepository.Create(tag)
	if err != nil {
		return models.Tag{}, err
	}
	return instance, nil
}

//Update ...
func (t TagService) Update(tagID uint,  tag models.Tag) (models.Tag, error) {
	instance, err := t.tagRepository.Update(tagID, tag)
	if err != nil {
		return models.Tag{}, err
	}
	return instance, nil
}

//Delete ...
func (t TagService) Delete(tagID uint) (error) {
	err := t.tagRepository.Delete(tagID)
	return err
}

//List ...
func (t TagService) List() (models.TagsList, error) {
	response, err := t.tagRepository.List()
	if err != nil {
		return models.TagsList{}, err
	}
	return models.TagsList{Data: response}, nil
}