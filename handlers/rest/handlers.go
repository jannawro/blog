package rest

import (
	"github.com/jannawro/blog/article"
)

type Handler struct {
	service *article.Service
}

func NewHandler(service *article.Service) *Handler {
	return &Handler{
		service: service,
	}
}
