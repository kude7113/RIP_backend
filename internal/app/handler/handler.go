package handler

import (
	"RIP/internal/app/repository"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	/*
		router.GET("/", h.AllFines)
		router.GET("/more/:id", h.FinesByID)
		router.POST("/delete/:id", h.DeleteResolution)
		router.POST("/add/:id", h.AddFinesToRes)
		router.GET("/resolution/:id", h.GetResolution)
	*/

	// домен услуги /Fines
	router.GET(FineDomain, h.AllFines)                         // Список штрафов
	router.GET(FineDomain+"/:id", h.FinesByID)                 // Штраф по ID
	router.POST(FineDomain+"/create", h.CreateFines)           // Добавление штрафа
	router.POST(FineDomain+"/img/:id", h.UploadImage)          // Добавление или замена изображения
	router.PUT(FineDomain+"/update/:id", h.UpdateFines)        // Редактирование штрафа
	router.DELETE(FineDomain+"/delete/:id", h.DeleteFines)     // Удаление штрафа
	router.POST(FineDomain+"/add/:id", h.AddFinesToResolution) // Добавление штрафа в последнее постановление

	// домен заявки /Resolutions
	router.GET(ResolutionDomain, h.AllResolutions)                    // Список постановлений
	router.GET(ResolutionDomain+"/:id", h.ResolutionByID)             // Постановление по ID
	router.PUT(ResolutionDomain+"/update/:id", h.UpdateResolution)    // Редактирование постановления
	router.PUT(ResolutionDomain+"/form", h.SetStatusByUser)           // Изменение статуса создателем
	router.PUT(ResolutionDomain+"/complete/:id", h.SetStatusByAdmin)  // Изменение статуса админом
	router.DELETE(ResolutionDomain+"/delete/:id", h.DeleteResolution) // Удаление постановления

	// домен м-м
	router.DELETE(FinResDomain+"/delete/:id", h.DeleteFR)  // Удаление из Fin_Res
	router.PUT(FinResDomain+"/count/:id", h.UpdateFRCount) // Изменение поля в Fin_Res

	// домен пользователя
	router.POST(UserDomain+"/register", h.RegistrUser)
	//router.POST(UserDomain+"/login", h.LoginUser)
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
