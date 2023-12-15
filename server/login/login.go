package login

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	if creds.Username == "admin" && creds.Password == "password123" {
		c.JSON(http.StatusOK, gin.H{"message": "Login successful!"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})

	}
}

//UNIMPLEMENTED
