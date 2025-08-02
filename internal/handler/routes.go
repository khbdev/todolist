package handler

import (
	"todolist/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
    r *gin.Engine,
    userHandler *UserHandler,
    profileHandler *ProfileHandler,
    settingHandler *SettingHandler,
    userUC *usecase.UserUsecase,
) {
    authMiddleware := AuthMiddleware(userUC)

    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "golang todo-app",
        })
    })

    // Auth
    r.POST("/register", userHandler.Register)
    r.POST("/login", userHandler.Login)
    r.POST("/lagout", authMiddleware, userHandler.Logout)
    // Profile
    r.GET("/profile/", authMiddleware, profileHandler.GetMyProfile)
    r.PUT("/profile/", authMiddleware, profileHandler.UpdateProfile)

    // Setting
    r.GET("/setting/", authMiddleware, settingHandler.GetSetting)
    r.PUT("/setting/", authMiddleware, settingHandler.UpdateSetting)

    // Static files
    r.Static("/storage/images", "pkg/storage/images")
}


