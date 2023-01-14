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
				switch gErr.Type {
				case gin.ErrorTypePrivate:
					log.Error().Msg("Error: gin.ErrorTypePrivate")
					err := gErr.Unwrap()
					var vErrs validator.ValidationErrors
					if errors.As(err, &vErrs) {
						RespondWithValidationError(ctx, err)
						return
					}
					switch err.(type) {
					case *pq.Error:
						RespondWithPqError(ctx, err)
						return
					// TODO: refactor
					case e.ErrEntityNotFound, e.ErrInvalidCurrencyType:
						ctx.JSON(http.StatusForbidden, errorResponse(err))
						return
					// TODO: check sql.Errs
					default:
						ctx.JSON(http.StatusInternalServerError, errorResponse(err))
						return
					}
				case gin.ErrorTypeBind:
					log.Error().Msg("Error: gin.ErrorTypeBind")
				case gin.ErrorTypePublic:
					log.Error().Msg("Error: gin.ErrorTypePublic")
				case gin.ErrorTypeRender:
					log.Error().Msg("Error: gin.ErrorTypeRecorder")
				case gin.ErrorTypeAny:
					log.Error().Msg("Error: gin.ErrorTypeAny")
				default:
					log.Error().Msg("Error: in Default")
					ctx.JSON(http.StatusInternalServerError, errorResponse(gErr))
				}
			}
		}
	}
}

type ApiError struct {
	Field string
	Msg   string
}

// errorResponse wraps error messages
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// errMsgForTag returns an exact error message depending on a tag of validation error
func errMsgForBindingTag(tag string, val interface{}) string {
	switch tag {
	case "required":
		return "This field is required"
	case "currency":
		return fmt.Sprintf("Invalid currency type: %v, should be on of: %v", val, utils.Currencies)
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too big"
	case "email":
		return "Invalid email"
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
				Msg:   errMsgForBindingTag(vErr.Tag(), vErr.Value()),
			}
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": out})
		return
	}
	ctx.JSON(http.StatusBadRequest, errorResponse(err))
}

// respondWithPqError responds with value
func RespondWithPqError(ctx *gin.Context, err error) {
	pqErr, ok := err.(*pq.Error)
	if ok {
		switch pqErr.Code.Name() {
		case "unique_violation", "foreign_key_violation":
			ctx.JSON(http.StatusForbidden, errorResponse(pqErr))
			return
		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(pqErr))
			return
		}
	}
	ctx.JSON(http.StatusInternalServerError, errorResponse(pqErr))
}

func RespondWithTransferError(ctx *gin.Context, err error) {

}
