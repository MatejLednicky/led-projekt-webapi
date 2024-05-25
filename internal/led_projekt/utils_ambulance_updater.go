package led_projekt

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/MatejLednicky/led-projekt-webapi/internal/db_service"
)

type treatmentUpdater = func(
    ctx *gin.Context,
    treatment *Treatment,
) (updatedTreatment *Treatment, responseContent interface{}, status int)

func updateTreatmentFunc(ctx *gin.Context, updater treatmentUpdater) {
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

    treatment, err := db.FindDocument(ctx, treatmentId)

    switch err {
    case nil:
        // continue
    case db_service.ErrNotFound:
        ctx.JSON(
            http.StatusNotFound,
            gin.H{
                "status":  "Not Found",
                "message": "Treatment not found",
                "error":   err.Error(),
            },
        )
        return
    default:
        ctx.JSON(
            http.StatusBadGateway,
            gin.H{
                "status":  "Bad Gateway",
                "message": "Failed to load treatment from database",
                "error":   err.Error(),
            })
        return
    }

    if !ok {
        ctx.JSON(
            http.StatusInternalServerError,
            gin.H{
                "status":  "Internal Server Error",
                "message": "Failed to cast treatment from database",
                "error":   "Failed to cast treatment from database",
            })
        return
    }

    updatedTreatment, responseObject, status := updater(ctx, treatment)

    if updatedTreatment != nil {
        err = db.UpdateDocument(ctx, treatmentId, updatedTreatment)
    } else {
        err = nil // redundant but for clarity
    }

    switch err {
    case nil:
        if responseObject != nil {
            ctx.JSON(status, responseObject)
        } else {
            ctx.AbortWithStatus(status)
        }
    case db_service.ErrNotFound:
        ctx.JSON(
            http.StatusNotFound,
            gin.H{
                "status":  "Not Found",
                "message": "Treatment was deleted while processing the request",
                "error":   err.Error(),
            },
        )
    default:
        ctx.JSON(
            http.StatusBadGateway,
            gin.H{
                "status":  "Bad Gateway",
                "message": "Failed to update treatment in database",
                "error":   err.Error(),
            })
    }

}