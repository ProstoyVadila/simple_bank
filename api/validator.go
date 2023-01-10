package api

import (
	// "reflect"
	// "sync"

	"errors"
	"fmt"
	"net/http"

	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

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

func registerValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", isValidCurrency)
	}
}

// isValidCurrency validate currency field
var isValidCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return utils.IsSupportedCurrency(currency)
	}
	return false
}
