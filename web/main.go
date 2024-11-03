/*
Copyright © 2024 Denys <https://github.com/AnyoneClown>
This is my license. There are many like it, but this one is mine.
My license is my best friend. It is my life. I must master it as I must
master my life.
*/
package web

import (
	"net/http"
	"time"

	"github.com/AnyoneClown/anydb/config"
	"github.com/AnyoneClown/anydb/utils"
	"github.com/AnyoneClown/anydb/web/gintemplrenderer"
	"github.com/AnyoneClown/anydb/web/templates"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Web() {
	engine := gin.Default()

	ginHtmlRenderer := engine.HTMLRender
	engine.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	// Disable trusted proxy warning.
	engine.SetTrustedProxies(nil)

	// Configure gin to use custom zup logger and give output to the file
	engine.Use(ginzap.Ginzap(utils.Log, time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(utils.Log, true))

	// Custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("port", portValidator)
	}

	// Main Page
	engine.GET("/", func(c *gin.Context) {
		r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, templates.DBConfigView(config.Configs))
		c.Render(http.StatusOK, r)
	})

	// API for configs
	api := engine.Group("/api")
	{
		handler := &Handler{}
		api.GET("/configs", handler.GetConfigs)
		api.GET("/configs/:id", handler.GetConfig)
		api.POST("/configs", handler.CreateConfig)
		api.DELETE("/configs/:id", handler.DeleteConfig)
		api.PUT("/configs/:id", handler.UpdateConfig)
		api.POST("/configs/select/:id", handler.SelectConfig)
	}
	engine.Run(":8080")
}
