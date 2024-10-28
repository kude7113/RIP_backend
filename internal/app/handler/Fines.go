package handler

import (
	"RIP/internal/app/ds"
	"RIP/internal/app/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// вызываются функции из репы, которые идут в бд
// то есть как бы прослойка между эндпоинтами и данными, которые идут их бд

func (h *Handler) AllFines(ctx *gin.Context) {
	searchFines := ctx.Query("searchFines")

	userId := 1

	resCount, err := h.Repository.GetResolutionLength(userId)
	resID, err := h.Repository.ResolutionByUserID(userId)

	var fines *[]ds.Fines
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
	}
	var request *ds.Fines
	err = ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	request.Fine_ID = id
	updateFine, err := h.Repository.UpdateFine(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, updateFine)
}

func (h *Handler) DeleteFines(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	err = h.Repository.DeleteFine(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Deleted",
	})
}

func (h *Handler) AddFinesToResolution(ctx *gin.Context) {
	userID := 1
	fineID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	err = h.Repository.AddFinesToResolution(userID, fineID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "added",
	})
}

func (h *Handler) AllResolutions(ctx *gin.Context) {
	allRes, err := h.Repository.ResolutionsList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, allRes)
}

func (h *Handler) ResolutionByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	result, err := h.Repository.GetResByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) UpdateResolution(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	var request *ds.Resolutions
	err = ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	request.Resolution_ID = id
	updateRes, err := h.Repository.UpdateRes(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, updateRes)
}

func (h *Handler) SetStatusByUser(ctx *gin.Context) {
	userID := 1

	result, err := h.Repository.SetStatusByUser(userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
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
	}

	result, err := h.Repository.SetStatusByAdmin(resID, newStatus)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) DeleteResolution(ctx *gin.Context) {
	resID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	result, err := h.Repository.DeleteResolution(resID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) DeleteFR(ctx *gin.Context) {
	resFID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	err = h.Repository.DeleteFR(resFID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Deleted",
	})
}
