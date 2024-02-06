package routes

import (
	"github.com/gin-gonic/gin"
	"unhush-backend/controllers"
)

func AppRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/login-fetch-profile", controllers.LoginAndGetProfile())
}
