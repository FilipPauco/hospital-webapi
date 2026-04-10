package hospital_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"slices"
)

type implVisitsAPI struct {
}

func NewVisitsApi() VisitsAPI {
	return &implVisitsAPI{}
}

func (o implVisitsAPI) GetVisits(c *gin.Context) {
	updateWardFunc(c, func(c *gin.Context, ward *Ward) (*Ward, interface{}, int) {
		result := ward.Visits
		if result == nil {
			result = []Visit{}
		}
		return nil, result, http.StatusOK
	})
}

func (o implVisitsAPI) CreateVisit(c *gin.Context) {
	updateWardFunc(c, func(c *gin.Context, ward *Ward) (*Ward, interface{}, int) {
		var visit Visit

		if err := c.ShouldBindJSON(&visit); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if visit.PatientId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Patient ID is required",
			}, http.StatusBadRequest
		}

		if visit.Date == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Date is required",
			}, http.StatusBadRequest
		}

		if visit.Time == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Time is required",
			}, http.StatusBadRequest
		}

		if visit.Id == "" || visit.Id == "@new" {
			visit.Id = uuid.NewString()
		}

		conflictIndx := slices.IndexFunc(ward.Visits, func(v Visit) bool {
			return v.Id == visit.Id
		})

		if conflictIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Visit already exists",
			}, http.StatusConflict
		}

		ward.Visits = append(ward.Visits, visit)
		return ward, visit, http.StatusOK
	})
}

func (o implVisitsAPI) GetVisit(c *gin.Context) {
	updateWardFunc(c, func(c *gin.Context, ward *Ward) (*Ward, interface{}, int) {
		visitId := c.Param("visitId")

		visitIndx := slices.IndexFunc(ward.Visits, func(v Visit) bool {
			return v.Id == visitId
		})

		if visitIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Visit not found",
			}, http.StatusNotFound
		}

		return nil, ward.Visits[visitIndx], http.StatusOK
	})
}

func (o implVisitsAPI) UpdateVisit(c *gin.Context) {
	updateWardFunc(c, func(c *gin.Context, ward *Ward) (*Ward, interface{}, int) {
		visitId := c.Param("visitId")
		var visit Visit

		if err := c.ShouldBindJSON(&visit); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if visit.Id != "" && visit.Id != visitId {
			return nil, gin.H{
				"status":  http.StatusForbidden,
				"message": "Visit ID in request body does not match path parameter",
			}, http.StatusForbidden
		}

		visit.Id = visitId

		visitIndx := slices.IndexFunc(ward.Visits, func(v Visit) bool {
			return v.Id == visitId
		})

		if visitIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Visit not found",
			}, http.StatusNotFound
		}

		ward.Visits[visitIndx] = visit
		return ward, visit, http.StatusOK
	})
}

func (o implVisitsAPI) DeleteVisit(c *gin.Context) {
	updateWardFunc(c, func(c *gin.Context, ward *Ward) (*Ward, interface{}, int) {
		visitId := c.Param("visitId")

		visitIndx := slices.IndexFunc(ward.Visits, func(v Visit) bool {
			return v.Id == visitId
		})

		if visitIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Visit not found",
			}, http.StatusNotFound
		}

		ward.Visits = append(ward.Visits[:visitIndx], ward.Visits[visitIndx+1:]...)
		return ward, nil, http.StatusNoContent
	})
}
