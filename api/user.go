package api

import (
	"net/http"
	"time"

	db "github.com/ProstoyVadila/simple_bank/db/sqlc"
	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
}

// createAccount creates a new account in db
func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	args := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}
	user, err := s.store.CreateUser(ctx, args)
	if err != nil {
		ctx.Error(err)
		return
	}
	resp := createUserResponse{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusOK, resp)
}
