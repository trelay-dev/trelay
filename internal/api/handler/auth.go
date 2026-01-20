package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aftaab/trelay/internal/api/response"
	"github.com/aftaab/trelay/internal/core/auth"
)

type AuthHandler struct {
	jwtManager *auth.JWTManager
	apiKeyHash string
}

func NewAuthHandler(jwtManager *auth.JWTManager, apiKeyHash string) *AuthHandler {
	return &AuthHandler{
		jwtManager: jwtManager,
		apiKeyHash: apiKeyHash,
	}
}

type loginRequest struct {
	APIKey string `json:"api_key"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.APIKey == "" {
		response.ValidationError(w, "api_key", "API key is required")
		return
	}

	if !auth.ValidateAPIKey(req.APIKey, h.apiKeyHash) {
		response.Unauthorized(w, "invalid API key")
		return
	}

	accessToken, err := h.jwtManager.GenerateAccessToken()
	if err != nil {
		response.InternalError(w)
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken()
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		response.ValidationError(w, "refresh_token", "refresh token is required")
		return
	}

	claims, err := h.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(w, "invalid refresh token")
		return
	}

	if !claims.IsRefreshToken() {
		response.Unauthorized(w, "invalid token type")
		return
	}

	accessToken, err := h.jwtManager.GenerateAccessToken()
	if err != nil {
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, tokenResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
		TokenType:    "Bearer",
	})
}
