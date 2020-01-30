package endpoint

import (
	"github.com/z7zmey/golang-microservice-boilerplate/internal/repository"
	"github.com/z7zmey/golang-microservice-boilerplate/internal/service"
)

// requests & responses

type GetCarRequest struct {
	ID int
}

type GetCarResponse struct {
	Products []repository.Car `json:"products"`
}

// endpoints

type BoilerplateEndpoint interface {
	GetCar(req GetCarRequest) (GetCarResponse, error)
}

type boilerplateEndpoint struct {
	service service.BoilerplateService
}

func NewBoilerplateEndpoint(s service.BoilerplateService) BoilerplateEndpoint {
	return &boilerplateEndpoint{
		service: s,
	}
}

func (e boilerplateEndpoint) GetCar(req GetCarRequest) (GetCarResponse, error) {
	p, err := e.service.GetCar(req.ID)

	return GetCarResponse{
		Products: p,
	}, err
}
