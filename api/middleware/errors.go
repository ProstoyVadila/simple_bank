package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ProstoyVadila/simple_bank/e"
	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// TODO: improve error handling structure
// Errors middleware handle all errors from handlers
func Errors() gin.HandlerFunc {
	log.Info().Msg("Setting Errors middleware")
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) > 0 {
			for _, gErr := range ctx.Errors {
				log.Error().Msg(fmt.Sprintf("Error: %v", gErr.Type))
				switch gErr.Type {
				case gin.ErrorTypePrivate:
					err := gErr.Unwrap()
					var vErrs validator.ValidationErrors
					if errors.As(err, &vErrs) {
						RespondWithValidationError(ctx, err)
					}
					switch err.(type) {
					case *pq.Error:
						RespondWithPqError(ctx, err)
					case e.ErrAccountNotFound, e.ErrInvalidCurrencyType:
						ctx.JSON(http.StatusForbidden, ErrorResponse(err))
					}
				case gin.ErrorTypeBind:
				case gin.ErrorTypePublic:
				case gin.ErrorTypeRender:
				case gin.ErrorTypeAny:
				default:
					ctx.JSON(http.StatusInternalServerError, ErrorResponse(gErr))
				}
			}
		}
	}
}

type ApiError struct {
	Field string
	Msg   string
}

// ErrorResponse wraps error messages
func ErrorResponse(err error) gin.H {
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

// RespondWithValidationError responds with a prettier field validation error message
func RespondWithValidationError(ctx *gin.Context, err error) {
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
	ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
}

// respondWithPqError responds with value
func RespondWithPqError(ctx *gin.Context, err error) {
	pqErr, ok := err.(*pq.Error)
	if ok {
		switch pqErr.Code.Name() {
		case "unique_violation", "foreign_key_violation":
			ctx.JSON(http.StatusForbidden, ErrorResponse(pqErr))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(pqErr))
			return
		}
	}
	ctx.JSON(http.StatusInternalServerError, ErrorResponse(pqErr))
}

func RespondWithTransferError(ctx *gin.Context, err error) {

}
