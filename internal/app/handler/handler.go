package handler

import (
	"RIP/docs"
	"RIP/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type Handler struct {
	Repository *repository.Repository
	Logger     *logrus.Logger
}

func NewHandler(l *logrus.Logger, r *repository.Repository) *Handler {
	return &Handler{
		Logger:     l,
		Repository: r,
	}
}

const (
	FineDomain       = "/fine"
	ResolutionDomain = "/resolution"
	FinResDomain     = "/fr"
	UserDomain       = "/user"
)

func (h *Handler) RegisterHandler(router *gin.Engine) {

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Разрешаем все источники
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK) // Возвращаем 200 OK на OPTIONS
			return
		}

		c.Next()
	})

	docs.SwaggerInfo.Title = "Fines for scooter"
	docs.SwaggerInfo.Description = "API server"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = "localhost:8000"

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// домен услуги /Fines
	router.GET(FineDomain, h.RoleMiddleware(AdminRole, UserRole, GuestRole), h.AllFines)              // Список штрафов
	router.GET(FineDomain+"/:id", h.RoleMiddleware(AdminRole, UserRole, GuestRole), h.FinesByID)      // Штраф по ID
	router.POST(FineDomain+"/create", h.RoleMiddleware(AdminRole), h.CreateFines)                     // Добавление штрафа
	router.POST(FineDomain+"/img/:id", h.RoleMiddleware(AdminRole), h.UploadImage)                    // Добавление или замена изображения
	router.PUT(FineDomain+"/update/:id", h.RoleMiddleware(AdminRole), h.UpdateFines)                  // Редактирование штрафа
	router.DELETE(FineDomain+"/delete/:id", h.RoleMiddleware(AdminRole), h.DeleteFines)               // Удаление штрафа
	router.POST(FineDomain+"/add/:id", h.RoleMiddleware(UserRole, AdminRole), h.AddFinesToResolution) // Добавление штрафа в последнее постановление

	// домен заявки /Resolutions
	router.GET(ResolutionDomain, h.RoleMiddleware(AdminRole), h.AllResolutions)                              // Список постановлений
	router.GET(ResolutionDomain+"/:id", h.RoleMiddleware(AdminRole, UserRole), h.ResolutionByID)             // Постановление по ID
	router.PUT(ResolutionDomain+"/update/:id", h.RoleMiddleware(AdminRole, UserRole), h.UpdateResolution)    // Редактирование постановления
	router.PUT(ResolutionDomain+"/form", h.RoleMiddleware(UserRole, AdminRole), h.SetStatusByUser)           // Изменение статуса создателем
	router.PUT(ResolutionDomain+"/complete/:id", h.RoleMiddleware(AdminRole), h.SetStatusByAdmin)            // Изменение статуса админом
	router.DELETE(ResolutionDomain+"/delete/:id", h.RoleMiddleware(AdminRole, UserRole), h.DeleteResolution) // Удаление постановления

	// домен м-м
	router.DELETE(FinResDomain+"/delete/:id", h.RoleMiddleware(AdminRole, UserRole), h.DeleteFR)  // Удаление из Fin_Res
	router.PUT(FinResDomain+"/count/:id", h.RoleMiddleware(AdminRole, UserRole), h.UpdateFRCount) // Изменение поля в Fin_Res

	// домен пользователя
	router.POST(UserDomain+"/register", h.RoleMiddleware(AdminRole, GuestRole), h.RegisterUser)
	router.POST(UserDomain+"/login", h.RoleMiddleware(AdminRole, GuestRole), h.LoginUser)
	router.POST(UserDomain+"/logout", h.RoleMiddleware(AdminRole, UserRole), h.LogoutUser)
	router.PUT(UserDomain+"/update", h.RoleMiddleware(AdminRole, UserRole), h.UpdateUser)

	router.GET(UserDomain+"/protected", h.RoleMiddleware(AdminRole), func(ctx *gin.Context) {
		userID := ctx.MustGet("user_id").(float64)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "user is authorized",
			"user_id": userID,
		})
	})
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	h.Logger.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
