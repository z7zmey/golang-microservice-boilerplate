package handler

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/z7zmey/golang-microservice-boilerplate/internal/endpoint"
	"github.com/z7zmey/golang-microservice-boilerplate/internal/transport/rest/server/models"
	"github.com/z7zmey/golang-microservice-boilerplate/internal/transport/rest/server/restapi/operations"
)

type ApiHandler struct {
	Endpoint endpoint.BoilerplateEndpoint
}

func NewApiHandler(endpoint endpoint.BoilerplateEndpoint) *ApiHandler {
	return &ApiHandler{
		Endpoint: endpoint,
	}
}

func (h *ApiHandler) ConfigureHandlers(api *operations.BoilerplateMicroserviceAPI) {
	api.HealthcheckHandler = operations.HealthcheckHandlerFunc(h.HealthCheck)
	api.CarHandler = operations.CarHandlerFunc(h.Car)
}

func (h *ApiHandler) HealthCheck(_ operations.HealthcheckParams) middleware.Responder {
	return operations.NewHealthcheckOK()
}

func (h *ApiHandler) Car(params operations.CarParams) middleware.Responder {
	resp, err := h.Endpoint.GetCar(endpoint.GetCarRequest{
		ID: int(params.ID),
	})

	if err != nil {
		return operations.NewCarBadRequest().WithPayload(&models.Error{
			Code: 400,
			Message: err.Error(),
		})
	}

	data := make([]*models.Car, len(resp.Products))
	for i, v := range resp.Products {
		data[i] = &models.Car{
			ID: toInt64(v.Id),
			Type: &v.Type,
			Manufacturer: &v.Manufacturer,
		}
	}

	return operations.NewCarOK().WithPayload(&operations.CarOKBody{
		Data: data,
	})
}

func toInt64(i int) *int64 {
	i64 := int64(i)
	return &i64
}
