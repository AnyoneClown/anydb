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

	engine.GET("/", func(c *gin.Context) {
		r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, templates.DBConfigView(config.Configs))
		c.Render(http.StatusOK, r)
	})

	engine.Run(":8080")
}
