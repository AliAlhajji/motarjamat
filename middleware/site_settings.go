package middleware

import (
	"net/http"

	"github.com/AliAlhajji/Motarjamat/models"
	"github.com/AliAlhajji/Motarjamat/utils"
	"github.com/gin-gonic/gin"
)

type settingsServer interface {
	GetSettings() (*models.SiteSettings, error)
}

type siteSettingsMiddleware struct {
	settingsServer settingsServer
}

func NewSettingsMiddleware(settingsServer settingsServer) *siteSettingsMiddleware {
	return &siteSettingsMiddleware{settingsServer: settingsServer}
}

func (m *siteSettingsMiddleware) GetSettings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		settings, err := m.settingsServer.GetSettings()
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"err": err})
			ctx.Abort()
			return
		}

		ctx.Set("siteSettings", settings)

		//If the site is closed and the user is not admin, show closed message
		user := utils.GetContextUser(ctx)
		if (user == nil || user.Role != "admin") && !settings.IsRunning {
			ctx.HTML(http.StatusOK, "closed_site.html", gin.H{"title": "site is closed", "announcement": settings.Announcement})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
