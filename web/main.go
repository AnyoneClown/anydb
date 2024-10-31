package web

import (
	"net/http"

	"github.com/AnyoneClown/anydb/config"
	"github.com/AnyoneClown/anydb/web/gintemplrenderer"
	"github.com/AnyoneClown/anydb/web/templates"
	"github.com/gin-gonic/gin"
)

func Web() {
	engine := gin.Default()
	engine.LoadHTMLFiles("./home.html")

	ginHtmlRenderer := engine.HTMLRender
	engine.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	// Disable trusted proxy warning.
	engine.SetTrustedProxies(nil)

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
	}
	engine.Run(":8080")
}
