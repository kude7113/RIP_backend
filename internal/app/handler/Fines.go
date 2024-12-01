package handler

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

// вызываются функции из репы, которые идут в бд
// то есть как бы прослойка между эндпоинтами и данными, которые идут их бд

func (h *Handler) AllFines(ctx *gin.Context) {
	searchFines := ctx.Query("searchFines")

	userId, _ := ctx.Get("user_id")

	// Объявляем переменные resCount и resID
	var resCount, resID int

	// Если userId не равен 0, выполняем запросы на получение данных
	if userId.(float64) != 0 {
		var err error
		resCount, err = h.Repository.GetResolutionLength(int(userId.(float64)))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		resID, err = h.Repository.ResolutionByUserID(int(userId.(float64)))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// Запрос на получение штрафов
	var fines *[]ds.Fines
	var err error
	if searchFines == "" {
		fines, err = h.Repository.GetAllFines()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		fines, err = h.Repository.SearchFines(searchFines)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	// Формируем ответ
	answer := models.FinesListWithRes{
		Fines:    fines,
		ResCount: resCount,
		ResID:    resID,
	}
	ctx.JSON(http.StatusOK, answer)
}

func (h *Handler) FinesByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.String(http.StatusBadRequest, "Invalid ID")
		return
	}

	fine, err := h.Repository.GetFinesByID(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, fine)
}

func (h *Handler) CreateFines(ctx *gin.Context) {
	var request *ds.Fines
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newFine, err := h.Repository.CreateFine(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, newFine)
}

func (h *Handler) UpdateFines(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var request *ds.Fines
	err = ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	request.Fine_ID = id
	updateFine, err := h.Repository.UpdateFine(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, updateFine)
}

func (h *Handler) DeleteFines(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = h.Repository.DeleteFine(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Deleted",
	})
}

func (h *Handler) AddFinesToResolution(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	fineID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = h.Repository.AddFinesToResolution(int(userID.(float64)), fineID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "added",
	})
}

func (h *Handler) AllResolutions(ctx *gin.Context) {
	dateFromQuery := ctx.Query("date_from")
	dateToQuery := ctx.Query("date_to")
	statusQuery := ctx.Query("status")

	var dateFrom, dateTo *time.Time

	// Если date_from присутствует, парсим дату
	if dateFromQuery != "" {
		parsedDateFrom, err := time.Parse("2006-01-02", dateFromQuery)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_from format. Use YYYY-MM-DD."})
			return
		}
		dateFrom = &parsedDateFrom
	}

	// Если date_to присутствует, парсим дату
	if dateToQuery != "" {
		parsedDateTo, err := time.Parse("2006-01-02", dateToQuery)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_to format. Use YYYY-MM-DD."})
			return
		}
		dateTo = &parsedDateTo
	}

	// Запрашиваем список резолюций с учётом опциональных параметров
	FilteredRes, err := h.Repository.ResolutionsList(dateFrom, dateTo, statusQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, FilteredRes)
}

func (h *Handler) ResolutionByID(ctx *gin.Context) {
	resID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid resID"})
	}

	result, err := h.Repository.GetResByID(resID)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) UpdateResolution(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var request *ds.Resolutions
	err = ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	request.Resolution_ID = id
	updateRes, err := h.Repository.UpdateRes(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, updateRes)
}

func (h *Handler) SetStatusByUser(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	result, err := h.Repository.SetStatusByUser(int(userID.(float64)))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) SetStatusByAdmin(ctx *gin.Context) {
	resID, err := strconv.Atoi(ctx.Param("id"))
	newStatus := ds.ApprovedStatus
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.Repository.SetStatusByAdmin(resID, newStatus)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) DeleteResolution(ctx *gin.Context) {
	resID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.Repository.DeleteResolution(resID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) DeleteFR(ctx *gin.Context) {
	resFID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.Repository.DeleteFR(resFID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Deleted",
	})
}

func (h *Handler) UpdateFRCount(ctx *gin.Context) {
	resFID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid ID",
		})
		return
	}

	var request ds.Fine_Resolutions
	if err = ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request payload",
			"details": err.Error(),
		})
		return
	}

	result, err := h.Repository.UpdateFRCount(resFID, request.Number)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to update record",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) UpdateUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var request *ds.Users
	err = ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	request.User_ID = userID

	result, err := h.Repository.UpdateUser(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) UploadImage(ctx *gin.Context) {
	// Считываем id из параметра URL
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Получаем файл из запроса
	request, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: image file is required"})
		return
	}

	// Открываем файл
	file, err := request.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open image"})
		return
	}
	defer file.Close()

	// Генерация имени файла
	fileExtension := filepath.Ext(request.Filename)
	fileName := strconv.Itoa(id) + fileExtension

	// Загружаем файл через репозиторий, предварительно удаляя существующий, если он есть
	imageURL, err := h.Repository.UploadImageAndUpdateURL(id, fileName, file, request.Size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ImageURL": imageURL})
}

func (h *Handler) RegisterUser(ctx *gin.Context) {
	var newUser ds.Users
	err := ctx.ShouldBindJSON(&newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.Repository.CreateUser(newUser.Login, newUser.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (h *Handler) LoginUser(ctx *gin.Context) {
	var request ds.Users
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.Repository.LoginUser(request.Login, request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, token)
}

func (h *Handler) LogoutUser(ctx *gin.Context) {
	value, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No user id found"})
		return
	}

	tokenString := extractTokenFromHeader(ctx.Request)
	if tokenString == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	err := h.Repository.LogoutUser(int(value.(float64)), tokenString)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"logout": "success",
	})
}
