package service

import "github.com/z7zmey/golang-microservice-boilerplate/internal/repository"

type BoilerplateService interface {
	GetCar(id int) ([]repository.Car, error)
}

type boilerplateService struct {
	carRepo repository.CarRepository
}

func NewSearchService(elasticsearch repository.CarRepository) *boilerplateService {
	return &boilerplateService{
		carRepo: elasticsearch,
	}
}

func (s boilerplateService) GetCar(id int) ([]repository.Car, error) {
	return s.carRepo.Find(id)
}
