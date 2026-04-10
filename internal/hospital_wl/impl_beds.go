package hospital_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"slices"
)

type implBedsAPI struct {
}

func NewBedsApi() BedsAPI {
	return &implBedsAPI{}
}

func (o implBedsAPI) GetBeds(c *gin.Context) {
	updateWardFunc(c, func(c *gin.Context, ward *Ward) (*Ward, interface{}, int) {
		result := ward.Beds
		if result == nil {
			result = []Bed{}
		}
		return nil, result, http.StatusOK
	})
}

func (o implBedsAPI) CreateBed(c *gin.Context) {
	updateWardFunc(c, func(c *gin.Context, ward *Ward) (*Ward, interface{}, int) {
		var bed Bed

		if err := c.ShouldBindJSON(&bed); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if bed.Number == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Bed number is required",
			}, http.StatusBadRequest
		}

		if bed.Status == "" {
			bed.Status = "free"
		}

		if bed.Id == "" || bed.Id == "@new" {
			bed.Id = uuid.NewString()
		}

		conflictIndx := slices.IndexFunc(ward.Beds, func(b Bed) bool {
			return b.Id == bed.Id
		})

		if conflictIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Bed already exists",
			}, http.StatusConflict
		}

		ward.Beds = append(ward.Beds, bed)
		return ward, bed, http.StatusOK
	})
}

func (o implBedsAPI) GetBed(c *gin.Context) {
	updateWardFunc(c, func(c *gin.Context, ward *Ward) (*Ward, interface{}, int) {
		bedId := c.Param("bedId")

		bedIndx := slices.IndexFunc(ward.Beds, func(b Bed) bool {
			return b.Id == bedId
		})

		if bedIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Bed not found",
			}, http.StatusNotFound
		}

		return nil, ward.Beds[bedIndx], http.StatusOK
	})
}

func (o implBedsAPI) UpdateBed(c *gin.Context) {
	updateWardFunc(c, func(c *gin.Context, ward *Ward) (*Ward, interface{}, int) {
		bedId := c.Param("bedId")
		var bed Bed

		if err := c.ShouldBindJSON(&bed); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if bed.Id != "" && bed.Id != bedId {
			return nil, gin.H{
				"status":  http.StatusForbidden,
				"message": "Bed ID in request body does not match path parameter",
			}, http.StatusForbidden
		}

		bed.Id = bedId

		bedIndx := slices.IndexFunc(ward.Beds, func(b Bed) bool {
			return b.Id == bedId
		})

		if bedIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Bed not found",
			}, http.StatusNotFound
		}

		ward.Beds[bedIndx] = bed
		return ward, bed, http.StatusOK
	})
}

func (o implBedsAPI) DeleteBed(c *gin.Context) {
	updateWardFunc(c, func(c *gin.Context, ward *Ward) (*Ward, interface{}, int) {
		bedId := c.Param("bedId")

		bedIndx := slices.IndexFunc(ward.Beds, func(b Bed) bool {
			return b.Id == bedId
		})

		if bedIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Bed not found",
			}, http.StatusNotFound
		}

		ward.Beds = append(ward.Beds[:bedIndx], ward.Beds[bedIndx+1:]...)
		return ward, nil, http.StatusNoContent
	})
}
