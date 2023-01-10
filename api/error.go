package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ApiError struct {
	Field string
	Msg   string
}

// errorResponse wraps error messages
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// errMsgForTag returns an exact error message depending on a tag of validation error
func errMsgForTag(tag string, val interface{}) string {
	switch tag {
	case "required":
		return "This field is required"
	case "currency":
		return fmt.Sprintf(
			"Invalid currency type: %v, should be: %v",
			val,
			utils.Currencies,
		)
	}
	return ""
}

// respondWithValidationError responds with a prettier field validation error message
func respondWithValidationError(ctx *gin.Context, err error) {
	var vErrs validator.ValidationErrors
	if errors.As(err, &vErrs) {
		out := make([]ApiError, len(vErrs))
		for i, vErr := range vErrs {
			out[i] = ApiError{
				Field: vErr.Field(),
				Msg:   errMsgForTag(vErr.Tag(), vErr.Value()),
			}
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": out})
		return
	}
	ctx.JSON(http.StatusBadRequest, errorResponse(err))
}
