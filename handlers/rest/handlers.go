package rest

import (
	"github.com/your-username/your-repo-name/article"
)

type Handler struct {
	service *article.Service
}

func NewHandler(service *article.Service) *Handler {
	return &Handler{
		service: service,
	}
}
