package gormongo

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router   *gin.Engine
	Registry *Registry
}

func NewApp(registry *Registry) *App {
	return &App{
		Router:   gin.Default(),
		Registry: registry,
	}
}

func (a *App) Run(addr string) {
	a.Router.Run(addr)
}

func GetModel[T IModel](r *Registry, name string) (T, error) {
	model, err := r.Get(name)
	if err != nil {
		var zero T
		return zero, err
	}
	typed, ok := model.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("model %s type mismatch", name)
	}
	return typed, nil
}
