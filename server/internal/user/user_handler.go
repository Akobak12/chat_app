package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (handler *Handler) CreateUser(context *gin.Context) {
	var userReq CreateUserReq
	if err := context.ShouldBindJSON(&userReq); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := handler.Service.CreateUser(context.Request.Context(), &userReq)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, res)
}

func (handler *Handler) Login(context *gin.Context) {
	var userReq LoginUserReq
	if err := context.ShouldBindJSON(&userReq); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRes, err := handler.Service.Login(context.Request.Context(), &userReq)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.SetCookie("jwt", userRes.AccessToken, 60*60*24, "/", "localhost", false, false)
	context.JSON(http.StatusOK, userRes) // Return the full user response including the access token
}

func (handler *Handler) Logout(context *gin.Context) {
	context.SetCookie("jwt", "", -1, "", "", false, true)
	context.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
