package well_known

import (
	"bytes"
	_ "embed"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

var (
	//go:embed template/host-meta.xml.template
	hostMetaXMLTemplate string
	//go:embed template/host-meta.json.template
	hostMetaJSONTemplate string
)

type service struct{}

func New() *service {
	return &service{}
}

func (e *service) Register(r *gin.Engine) {
	wellKnownRoutes := r.Group("/.well-known")
	wellKnownRoutes.GET("/host-meta", e.handleHostMeta)
	wellKnownRoutes.GET("/host-meta.json", e.handleHostMeta)
	wellKnownRoutes.GET("/webfinger", e.handleWebfinger)
}

func (e *service) handleHostMeta(c *gin.Context) {
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

func (e *service) handleWebfinger(c *gin.Context) {
	host := c.Request.Host
	resource := c.Query("resource")
	if resource == "" {
		c.String(http.StatusBadRequest, "resource is required")
		return
	}
	sub, err := url.Parse(resource)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid resource")
		return
	}
	if sub.Scheme != "acct" {
		c.String(http.StatusBadRequest, "invalid resource")
		return
	}
	username := strings.Split(sub.Opaque, "@")[0]

	res := gin.H{
		"subject": "acct:" + username + "@" + host,
		"aliases": []string{
			"https://" + host + "/@" + username,
			"https://" + host + "/users/" + username,
		},
		"links": []interface{}{
			map[string]string{
				"rel":  "self",
				"type": "application/activity+json",
				"href": "https://" + host + "/users/" + username,
			},
		},
	}
	c.Header("Content-Type", "application/jrd+json")
	c.Header("Cache-Control", "max-age=3600, public")
	c.JSON(http.StatusOK, res)
}