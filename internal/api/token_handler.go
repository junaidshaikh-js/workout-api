package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/junaidshaikh-js/workout-api/internal/store"
	"github.com/junaidshaikh-js/workout-api/internal/tokens"
	"github.com/junaidshaikh-js/workout-api/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type CreateTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req CreateTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		h.logger.Printf("ERROR: decodeCreateToken: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "invalid request"})
		return
	}

	// fetch user if exists
	user, err := h.userStore.GetUserByUserName(req.Username)

	if err != nil || user == nil {
		h.logger.Printf("ERROR: getUserByUsername: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}

	isPasswordCorrect, err := user.PasswordHash.Matches(req.Password)

	if err != nil {
		h.logger.Printf("ERROR: passwordHashMatch: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}

	if !isPasswordCorrect {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelop{"error": "invalid credentials"})
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)

	if err != nil {
		h.logger.Printf("ERROR: createNewToken: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelop{"auth_token": token})
}
