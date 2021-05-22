package repositories

import (
	"news-topic-api/infrastructures"
	"news-topic-api/models"
)

//ITagRepository interface for tag repository
type ITagRepository interface {
	Create(tag models.Tag) (models.Tag, error)
	Update(tagID uint, tag models.Tag) (models.Tag, error)
	Delete(tagID uint) (error)
	List() ([]models.Tag, error)
}

//TagRepository ...
type TagRepository struct{
}

//Create ...
func (t TagRepository) Create(tag models.Tag) (models.Tag, error) {
	db := infrastructures.GetDB()
	err := db.Create(&tag).Error
	return tag, err
}

//Update ...
func (t TagRepository) Update(tagID uint, tag models.Tag) (models.Tag, error) {
	var targetTag models.Tag
	db := infrastructures.GetDB()
	err := db.Where("id = ?", tagID).First(&targetTag).Error
	if err != nil {
		return models.Tag{}, err
	}
	updateData := map[string]interface{} {
		"name": tag.Name,
	}
	err = db.Model(&targetTag).Omit("created_at").Updates(updateData).Error
	if err != nil {
		return models.Tag{}, err
	}
	return targetTag, nil
}

//Delete ...
func (t TagRepository) Delete(tagID uint) (error) {
	var targetTag models.Tag
	db := infrastructures.GetDB()
	err := db.Where("id = ?", tagID).First(&targetTag).Error
	if err != nil {
		return err
	}
	err = db.Delete(&targetTag).Error
	if err != nil {
		return err
	}
	return nil
}


//List ...
func (t TagRepository) List() ([]models.Tag, error) {
	var tagsList []models.Tag
	db := infrastructures.GetDB()
	querySearch := db.Table("tags")
	err := querySearch.Order("id DESC").Find(&tagsList).Error
	if err != nil {
		return []models.Tag{}, err
	}
	return tagsList, nil
}
