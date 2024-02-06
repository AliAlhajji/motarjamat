package controllers

import (
	"github.com/AliAlhajji/Motarjamat/middleware"
	"github.com/AliAlhajji/Motarjamat/models"
	"github.com/gin-gonic/gin"
)

func getContextUser(ctx *gin.Context) *models.User {
	contextUser, ok := ctx.Get(middleware.ContextUser)
	if !ok {
		return nil
	}

	user, ok := contextUser.(*models.User)
	if !ok {
		return nil

	}

	return user
}
