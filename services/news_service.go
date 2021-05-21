
package services

import (
	"news-topic-api/models"
	"news-topic-api/repositories"
)


//INewsService interface for news service
type INewsService interface {
	Create(news models.News) (models.News, error)
	Update(newsID uint,  news models.News) (models.News, error)
	Delete(newsID uint) (error)
	List(queryParams map[string]string) (models.NewsList, error)
	GetDetail(newsID uint) (models.News, error)
}

//NewsService ...
type NewsService struct {
	newsRepository repositories.INewsRepository
}

//InitNewsService initialize a permohonan penyitaan service instance with specific news repository
func InitNewsService(newsRepository repositories.INewsRepository) INewsService {
	newsService := new(NewsService)
	newsService.newsRepository = newsRepository
	return newsService
}

//Create ...
func (n NewsService) Create(news models.News) (models.News, error) {
	instance, err := n.newsRepository.Create(news)
	if err != nil {
		return models.News{}, err
	}
	return instance, nil
}

//Update ...
func (n NewsService) Update(newsID uint,  news models.News) (models.News, error) {
	instance, err := n.newsRepository.Update(newsID, news)
	if err != nil {
		return models.News{}, err
	}
	return instance, nil
}

//Delete ...
func (n NewsService) Delete(newsID uint) (error) {
	err := n.newsRepository.Delete(newsID)
	return err
}

//List list all news that are matched with provided filters
func (n NewsService) List(queryParams map[string]string) (models.NewsList, error) {
	response, err := n.newsRepository.List(queryParams)
	if err != nil {
		return models.NewsList{}, err
	}
	return models.NewsList{Data: response}, nil
}

//GetDetail ...
func (n NewsService) GetDetail(newsID uint) (models.News, error) {
	response, err := n.newsRepository.GetByID(newsID)
	if err != nil {
		return models.News{}, err
	}
	return response, nil
}