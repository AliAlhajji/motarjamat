package controllers

import (
	"net/http"

	"github.com/AliAlhajji/Motarjamat/models"
	"github.com/gin-gonic/gin"
)

type SettingsServer interface {
	Update(*models.SiteSettings) error
}

type settingsController struct {
	settingsServer SettingsServer
}

func NewSettingsController(categoryServer SettingsServer) *settingsController {
	return &settingsController{
		settingsServer: categoryServer,
	}
}

func (c *settingsController) AdminHome(ctx *gin.Context) {
	data := gin.H{}
	data["title"] = "Admin Home"

	data["settings"] = ctx.MustGet("siteSettings")

	ctx.HTML(http.StatusOK, "admin_home.html", data)
}

func (c *settingsController) UpdateSettings(ctx *gin.Context) {
	var settings models.SiteSettings
	data := gin.H{}

	err := ctx.Bind(&settings)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"err": err})
		return
	}

	err = c.settingsServer.Update(&settings)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "error.html", gin.H{"err": err})
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/admin")
	data["title"] = "Admin Home"
	data["msg"] = "Updated"
	data["settings"] = ctx.MustGet("siteSettings")
	ctx.HTML(http.StatusOK, "admin_home.html", data)
}
