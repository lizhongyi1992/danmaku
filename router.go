package main

import "github.com/gin-gonic/gin"

type App struct {
}

func NewApp(config Config) *App {
	app := &App{}
	return app
}

func (p *App) Stop() {
}

func (p *App) danmaku_all(c *gin.Context) {
}
