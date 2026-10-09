package pagination

import (
	"github.com/gin-gonic/gin"
)

// Params holds parsed pagination parameters.
type Params struct {
	Page    int `form:"page"`
	PerPage int `form:"limit"`
}

// GetParams extracts and normalizes pagination query params from the Gin context.
func GetParams(c *gin.Context) Params {
	var params Params

	if err := c.ShouldBindQuery(&params); err != nil {
		return Params{
			Page:    1,
			PerPage: 10,
		}
	}

	if params.Page < 1 {
		params.Page = 1
	}

	if params.PerPage < 1 {
		params.PerPage = 10
	}

	return params
}
