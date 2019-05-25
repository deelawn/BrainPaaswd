package users

import "github.com/deelawn/BrainPaaswd/services"

type Service struct {
	services.Service
}

func NewService(service services.Service) *Service {

	return &Service{service}
}