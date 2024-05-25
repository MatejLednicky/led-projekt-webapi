package led_projekt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
  "github.com/MatejLednicky/led-projekt-webapi/internal/db_service"
)

// Nasledujúci kód je kópiou vygenerovaného a zakomentovaného kódu zo súboru api_led_treatment_list.go

// CreateTreatment - Saves new treatment
func (this *implLedTreatmentListAPI) CreateTreatment(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists { 
    ctx.JSON(http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Treatment])
	if !ok { 
    ctx.JSON( http.StatusInternalServerError, 
      gin.H{
				"status":  "Internal Server Error",
				"message": "db context is not of required type",
				"error":   "cannot cast db context to db_service.DbService",
			})
		return
	}

	treatment := Treatment{}
	err := ctx.BindJSON(&treatment)
	if err != nil {
		ctx.JSON( http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		return
	}

	if treatment.Id == "@new" {
		treatment.Id = uuid.New().String()
	}

	err = db.CreateDocument(ctx, treatment.Id, &treatment)

	switch err {
	case nil:
		ctx.JSON(http.StatusCreated, treatment)
	case db_service.ErrConflict:
		ctx.JSON(http.StatusConflict,
			gin.H{
				"status":  "Conflict",
				"message": "Treatment already exists",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create treatment in database",
				"error":   err.Error(),
			},
		)
	}
}

// DeleteTreatment - Deletes specific treatment
func (this *implLedTreatmentListAPI) DeleteTreatment(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Treatment])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	treatmentId := ctx.Param("treatmentId")
	err := db.DeleteDocument(ctx, treatmentId)

	switch err {
	case nil:
		ctx.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Treatment not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to delete treatment from database",
				"error":   err.Error(),
			})
	}
}

// GetTreatments - Provides the treatments list
func (this *implLedTreatmentListAPI) GetTreatments(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Treatment])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	treatments, err := db.FindDocuments(ctx)

	switch err {
	case nil:
		ctx.JSON(
			http.StatusOK,
			treatments,
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to find all treatments from database",
				"error":   err.Error(),
			},
		)
	}
}

// GetTreatmentDetail - Provides details about treatment
func (this *implLedTreatmentListAPI) GetTreatmentDetail(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Treatment])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	treatmentId := ctx.Param("treatmentId")

	if treatmentId == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "Treatment ID is required",
			})
		return
	}

	treatment, err := db.FindDocument(ctx, treatmentId)

	switch err {
	case nil:
		ctx.JSON(
			http.StatusOK,
			treatment,
		)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Treatment not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to load treatment from database",
				"error":   err.Error(),
			})
	}
}

// UpdateTreatment - Updates specific treatment
func (this *implLedTreatmentListAPI) UpdateTreatment(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Treatment])
	if !ok {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	updatedTreatment := Treatment{}

	if err := ctx.ShouldBindJSON(&updatedTreatment); err != nil {
		ctx.JSON(http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			},
		)
		return
	}

	if !updatedTreatment.StartDate.IsZero() && !updatedTreatment.EndDate.IsZero() {
		if updatedTreatment.EndDate.Before(updatedTreatment.StartDate) {
			ctx.JSON(http.StatusBadRequest,
				gin.H{
					"status":  http.StatusBadRequest,
					"message": "End date is before start date",
					"error":   "End date is before start date",
				},
			)
			return
		}
	}

	treatmentId := ctx.Param("treatmentId")

	if treatmentId == "" {
		ctx.JSON(http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "Treatment ID is required",
				"error":   "Treatment ID is required",
			},
		)
		return
	}


	err := db.UpdateDocument(ctx, treatmentId, &updatedTreatment)

	switch err {
	case nil:
		ctx.JSON(http.StatusOK, updatedTreatment)
	case db_service.ErrNotFound:
		ctx.JSON(http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Treatment not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to update treatment",
				"error":   err.Error(),
			})

	}
}
