package controller

import (
	"net/http"
	"pamagi/domain"
	"pamagi/domain/dto"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthUsecase domain.AuthUsecase
}

// Register godoc
// @Summary Register user baru
// @Description Mendaftarkan akun menggunakan email dan password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var request dto.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	// Menampung objek AuthResponse berupa token hasil pendaftaran sukses
	response, err := ac.AuthUsecase.Register(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	// Kembalikan status 200 OK beserta token lengkap
	c.JSON(http.StatusOK, response)
}

// Login godoc
// @Summary Login user
// @Description Autentikasi user dan mendapatkan JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Router /login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var request dto.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	response, err := ac.AuthUsecase.Login(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}


// Logout godoc
// @Summary Logout user
// @Description Menghapus sesi token pengguna dari sistem (Redis)
// @Tags Auth
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} domain.SuccessResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /logout [post]
func (ac *AuthController) Logout(c *gin.Context) {
	// Ambil ID user dari Middleware
	userID := c.GetString("x-user-id")
	
	err := ac.AuthUsecase.Logout(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Gagal melakukan logout"})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Logout berhasil"})
}

// GetProfile godoc
// @Summary Ambil Data User Saat Ini
// @Description Mendapatkan info profil, bahasa ibu, dan bahasa yang dipelajari
// @Tags Profile
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Router /users/me [get]
func (ac *AuthController) GetProfile(c *gin.Context) {
	userID := c.GetString("x-user-id")
	
	response, err := ac.AuthUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Gagal mengambil data profil"})
		return
	}

	c.JSON(http.StatusOK, response)
}


// RefreshToken godoc
// @Summary Refresh Access Token
// @Description Mendapatkan access token baru menggunakan refresh token tanpa harus login ulang
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Router /refresh [post]
func (ac *AuthController) Refresh(c *gin.Context) {
	var request dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	// Panggil usecase
	response, err := ac.AuthUsecase.RefreshToken(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}