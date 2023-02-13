package hostmeta

import (
	"bytes"
	_ "embed"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

var (
	//go:embed template/host-meta.xml.template
	hostMetaXMLTemplate string
	//go:embed template/host-meta.json.template
	hostMetaJSONTemplate string
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/.well-known/host-meta", handleHostMeta)
	r.GET("/.well-known/host-meta.json", handleHostMeta)
}

func handleHostMeta(c *gin.Context) {
	accept := c.GetHeader("Accept")
	c.Header("Cache-Control", "max-age=3600, public")
	switch accept {
	case "application/json":
		tmpl := template.New("host-meta.json.template")
		tmpl, err := tmpl.Parse(hostMetaJSONTemplate)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		buf := bytes.NewBuffer(nil)
		tmpl.Execute(buf, map[string]interface{}{
			"Host": c.Request.Host,
		})
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, buf.String())
	default:
		tmpl := template.New("host-meta.xml.template")
		tmpl, err := tmpl.Parse(hostMetaXMLTemplate)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		buf := bytes.NewBuffer(nil)
		tmpl.Execute(buf, map[string]interface{}{
			"Host": c.Request.Host,
		})
		c.Header("Content-Type", "application/xrd+xml; charset=utf-8")
		c.String(http.StatusOK, buf.String())
	}
}
