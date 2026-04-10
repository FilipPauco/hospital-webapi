package hospital_wl

import (
	"net/http"

	"github.com/FilipPauco/hospital-webapi/internal/db_service"
	"github.com/gin-gonic/gin"
)

type wardUpdater = func(
	ctx *gin.Context,
	ward *Ward,
) (updatedWard *Ward, responseContent interface{}, status int)

func updateWardFunc(ctx *gin.Context, updater wardUpdater) {
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

	db, ok := value.(db_service.DbService[Ward])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of required type",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	wardId := ctx.Param("wardId")

	ward, err := db.FindDocument(ctx, wardId)

	switch err {
	case nil:
		// continue
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Ward not found",
				"error":   err.Error(),
			},
		)
		return
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to load ward from database",
				"error":   err.Error(),
			},
		)
		return
	}

	if !ctx.IsAborted() {
		updatedWard, responseObject, status := updater(ctx, ward)

		if updatedWard != nil {
			err = db.UpdateDocument(ctx, wardId, updatedWard)
		}

		switch err {
		case nil:
			if responseObject != nil {
				ctx.JSON(status, responseObject)
			} else if status == http.StatusNoContent {
				ctx.AbortWithStatus(status)
			}
		case db_service.ErrNotFound:
			ctx.JSON(
				http.StatusNotFound,
				gin.H{
					"status":  "Not Found",
					"message": "Ward was deleted while processing the request",
					"error":   err.Error(),
				},
			)
		case db_service.ErrConflict:
			ctx.JSON(
				http.StatusConflict,
				gin.H{
					"status":  "Conflict",
					"message": "Ward has been modified concurrently",
					"error":   err.Error(),
				},
			)
		default:
			ctx.JSON(
				http.StatusBadGateway,
				gin.H{
					"status":  "Bad Gateway",
					"message": "Failed to update ward in database",
					"error":   err.Error(),
				},
			)
		}
	}
}
