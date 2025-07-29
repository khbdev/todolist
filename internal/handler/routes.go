package handler

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, userHandler *UserHandler) {
    r.POST("/register", userHandler.Register)
    r.POST("/login", userHandler.Login)
}