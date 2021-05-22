package repositories

import (
	"news-topic-api/infrastructures"
	"news-topic-api/models"
	"strconv"
)

//INewsRepository interface for news repository
type INewsRepository interface {
	Create(news models.News) (models.News, error)
	Update(newsID uint, news models.News) (models.News, error)
	Delete(newsID uint) (error)
	GetByID(penyitaanID uint) (models.News, error)
	List(queryParams map[string]string) ([]models.News, error)
}

//NewsRepository ...
type NewsRepository struct{
}

//Create ...
func (n NewsRepository) Create(news models.News) (models.News, error) {
	db := infrastructures.GetDB()
	err := db.Create(&news).Error
	db.Model(&news).Association("Tags").Find(&news.Tags)
	return news, err
}

//Update ...
func (n NewsRepository) Update(newsID uint, news models.News) (models.News, error) {
	var targetNews models.News
	db := infrastructures.GetDB()
	err := db.Where("id = ?", newsID).First(&targetNews).Error
	if err != nil {
		return models.News{}, err
	}
	updateData := map[string]interface{} {
		"title": news.Title,
		"thumbnail": news.Thumbnail,
		"summary": news.Summary,
		"content": news.Content,
		"topic": news.Topic,
		"status": news.Status,
	}
	err = db.Model(&targetNews).Omit("created_at").Updates(updateData).Error
	if err != nil {
		return models.News{}, err
	}
	db.Model(&targetNews).Association("Tags").Replace(news.Tags)
	db.Model(&targetNews).Association("Tags").Find(&targetNews.Tags)
	return targetNews, nil
}

//Delete ...
func (n NewsRepository) Delete(newsID uint) (error) {
	var targetNews models.News
	db := infrastructures.GetDB()
	err := db.Where("id = ?", newsID).First(&targetNews).Error
	if err != nil {
		return err
	}
	err = db.Delete(&targetNews).Error
	if err != nil {
		return err
	}
	return nil
}

//GetByID ...
func (n NewsRepository) GetByID(newsID uint) (models.News, error) {
	var targetNews models.News
	db := infrastructures.GetDB()
	queryByID := db.Where("id = ?", newsID)
	err := queryByID.First(&targetNews).Error
	if err != nil {
		return models.News{}, err
	}
	db.Model(&targetNews).Association("Tags").Find(&targetNews.Tags)
	return targetNews, nil
}

//List retrieve list of news given filters (topic, status, etc)
func (n NewsRepository) List(queryParams map[string]string) ([]models.News, error) {
	var newsList []models.News
	db := infrastructures.GetDB()
	querySearch := db.Table("news")
	status := queryParams["status"]
	topic := queryParams["topic"]
	tag := queryParams["tag"]
	if tag != "" {
		tagID, err := strconv.Atoi(tag)
		if err != nil {
			return []models.News{}, err
		}
		querySearch = querySearch.Joins("JOIN news_tag ON news_tag.news_id = news.id AND news.id = ?", uint(tagID))
	}
	if topic != "" {
		querySearch = querySearch.Where("topic = ?",topic)
	}
	if status != "" {
		querySearch = querySearch.Where("status = ?",status)
	}
	err := querySearch.Order("id DESC").Find(&newsList).Error
	if err != nil {
		return []models.News{}, err
	}
	for i := range newsList {
		db.Model(&newsList[i]).Association("Tags").Find(&newsList[i].Tags)
	}
	return newsList, nil
}