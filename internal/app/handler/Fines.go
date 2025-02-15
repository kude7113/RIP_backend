package handler

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/models"
	"RIP/internal/app/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

// AllFines godoc
// @Summary Get all fines
// @Description Get list of all fines or search fines by keyword
// @Tags Fines
// @Produce json
// @Param searchFines query string false "Search fines"
// @Success 200 {object} models.FinesListWithRes
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /fine [get]
func (h *Handler) AllFines(ctx *gin.Context) {
	searchFines := ctx.Query("searchFines")

	userId, _ := ctx.Get("user_id")

	var resCount, resID int

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

	answer := models.FinesListWithRes{
		Fines:    fines,
		ResCount: resCount,
		ResID:    resID,
	}
	ctx.JSON(http.StatusOK, answer)
}

// FinesByID godoc
// @Summary Get fine by ID
// @Description Get a fine by its ID
// @Tags Fines
// @Produce json
// @Param id path int true "Fine ID"
// @Success 200 {object} ds.Fines
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /fine/{id} [get]
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

// CreateFines godoc
// @Summary Create a new fine
// @Description Create a new fine
// @Tags Fines
// @Accept json
// @Produce json
// @Param fine body ds.Fines true "Fine data"
// @Success 201 {object} ds.Fines
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /fine/create [post]
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

// UpdateFines godoc
// @Summary Update a fine by ID
// @Description Update fine details by its ID
// @Tags Fines
// @Accept json
// @Produce json
// @Param id path int true "Fine ID"
// @Param fine body ds.Fines true "Fine data"
// @Success 200 {object} ds.Fines
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /fine/update/{id} [put]
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

// DeleteFines godoc
// @Summary Delete a fine by ID
// @Description Delete fine by its ID
// @Tags Fines
// @Produce json
// @Param id path int true "Fine ID"
// @Success 200 {string} string "Deleted"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /fines/delete/{id} [delete]
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

// AddFinesToResolution godoc
// @Summary Add a fine to a resolution
// @Description Add fine to a user's resolution
// @Tags Fines
// @Produce json
// @Param id path int true "Fine ID"
// @Success 200 {string} string "added"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /fine/add/{id} [post]
func (h *Handler) AddFinesToResolution(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	fineID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вызываем репозиторий, который теперь возвращает `resId`
	resId, err := h.Repository.AddFinesToResolution(int(userID.(float64)), fineID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ✅ Теперь возвращаем `resId` клиенту
	ctx.JSON(http.StatusOK, gin.H{
		"message": "added",
		"resId":   resId,
	})
}

// AllResolutions godoc
// @Summary Get all resolutions
// @Description Get list of all resolutions with optional filters for date range and status
// @Tags Resolutions
// @Produce json
// @Param date_from query string false "Filter by start date (YYYY-MM-DD)"
// @Param date_to query string false "Filter by end date (YYYY-MM-DD)"
// @Param status query string false "Filter by resolution status"
// @Success 200 {object} []ds.Resolutions
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /resolution [get]
func (h *Handler) AllResolutions(ctx *gin.Context) {
	dateFromQuery := ctx.Query("date_from")
	dateToQuery := ctx.Query("date_to")
	statusQuery := ctx.Query("status")

	var dateFrom, dateTo *time.Time

	if dateFromQuery != "" {
		parsedDateFrom, err := time.Parse("2006-01-02", dateFromQuery)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_from format. Use YYYY-MM-DD."})
			return
		}
		dateFrom = &parsedDateFrom
	}

	if dateToQuery != "" {
		parsedDateTo, err := time.Parse("2006-01-02", dateToQuery)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_to format. Use YYYY-MM-DD."})
			return
		}
		dateTo = &parsedDateTo
	}

	FilteredRes, err := h.Repository.ResolutionsList(dateFrom, dateTo, statusQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, FilteredRes)
}

// ResolutionByID godoc
// @Summary Get resolution by ID
// @Description Get a resolution by its ID
// @Tags Resolutions
// @Produce json
// @Param id path int true "Resolution ID"
// @Success 200 {object} ds.Resolutions
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /resolution/{id} [get]
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

// UpdateResolution godoc
// @Summary Update a resolution by ID
// @Description Update resolution details by its ID
// @Tags Resolutions
// @Accept json
// @Produce json
// @Param id path int true "Resolution ID"
// @Param resolution body ds.Resolutions true "Resolution data"
// @Success 200 {object} ds.Resolutions
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /resolutions/update/{id} [put]
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

// SetStatusByUser godoc
// @Summary Set status of resolution by user
// @Description Set status of a resolution for a user
// @Tags Resolutions
// @Produce json
// @Success 200 {object} ds.Resolutions
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /resolution/form [put]
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

// SetStatusByAdmin godoc
// @Summary Set status of resolution by admin
// @Description Set status of a resolution by admin
// @Tags Resolutions
// @Produce json
// @Param id path int true "Resolution ID"
// @Success 200 {object} ds.Resolutions
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /resolution/complete/ [put]
func (h *Handler) SetStatusByAdmin(ctx *gin.Context) {
	// Получаем id из параметров URL и преобразуем в число
	resID, err := strconv.Atoi(ctx.Param("id"))
	newStatus := ds.ApprovedStatus
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Обновляем статус в репозитории и получаем обновлённую структуру ds.Resolutions
	result, err := h.Repository.SetStatusByAdmin(resID, newStatus)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Генерируем QR‑код для обновлённого заказа
	qrCode, err := services.GenerateResolutionQR(*result)
	println(qrCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка генерации QR-кода: " + err.Error(),
		})
		return
	}
	result.Qr = qrCode
	_, err = h.Repository.UpdateRes(result)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// DeleteResolution godoc
// @Summary Delete a resolution by ID
// @Description Delete resolution by its ID
// @Tags Resolutions
// @Produce json
// @Param id path int true "Resolution ID"
// @Success 200 {object} ds.Resolutions
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /resolution/delete/{id} [delete]
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

// DeleteFR godoc
// @Summary Delete fine-resolution link by ID
// @Description Delete fine-resolution link by its ID
// @Tags Fines-Resolutions
// @Produce json
// @Param id path int true "Fine-Resolution ID"
// @Success 200 {string} string "Deleted"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /fr/delete/{id} [delete]
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

// UpdateFRCount godoc
// @Summary Update fine-resolution count
// @Description Update the count of a fine-resolution link
// @Tags Fines-Resolutions
// @Accept json
// @Produce json
// @Param id path int true "Fine-Resolution ID"
// @Param request body ds.Fine_Resolutions true "Fine-Resolution data"
// @Success 200 {object} ds.Fine_Resolutions
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /fr/count/{id} [put]
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

// UpdateUser godoc
// @Summary Update user by ID
// @Description Update user details by its ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body ds.Users true "User data"
// @Success 200 {object} ds.Users
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users/update [put]
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

// UploadImage godoc
// @Summary Upload an image by ID
// @Description Upload an image for a specific ID
// @Tags Images
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "ID"
// @Param image formData file true "Image file"
// @Success 200 {object} map[string]string "ImageURL"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /fine/img/{id} [post]
func (h *Handler) UploadImage(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	request, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: image file is required"})
		return
	}

	file, err := request.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open image"})
		return
	}
	defer file.Close()

	fileExtension := filepath.Ext(request.Filename)
	fileName := strconv.Itoa(id) + fileExtension

	imageURL, err := h.Repository.UploadImageAndUpdateURL(id, fileName, file, request.Size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ImageURL": imageURL})
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user by providing login and password
// @Tags Users
// @Accept json
// @Produce json
// @Param user body ds.Users true "User data"
// @Success 201 {object} ds.Users
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /user/register [post]
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

// LoginUser godoc
// @Summary Login a user
// @Description Login a user and generate a token
// @Tags Users
// @Accept json
// @Produce json
// @Param user body ds.Users true "User credentials"
// @Success 200 {string} string "token"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /user/login [post]
func (h *Handler) LoginUser(ctx *gin.Context) {
	var request ds.Users
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, isAdmin, err := h.Repository.LoginUser(request.Login, request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token":   token,
		"login":   request.Login,
		"isAdmin": isAdmin,
	})
}

// LogoutUser godoc
// @Summary Logout a user
// @Description Logout a user by invalidating the token
// @Tags Users
// @Produce json
// @Success 200 {string} string "logout: success"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /user/logout [post]
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
