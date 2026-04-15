package handlers

import (
	apperrors "AuthService/internal/errors"
	"AuthService/internal/models"
	"AuthService/internal/service"
	"encoding/json"
	"errors"
	"net/http"
)

type Handlers struct {
	authService *service.AuthService
}

func NewHandler(authService *service.AuthService) *Handlers {
	return &Handlers{authService: authService}
}

func respondWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJson(w, statusCode, map[string]string{"error": message})
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var input models.RegisterUserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный формат данных")
		return
	}

	if err := input.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authService.Register(input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка при регистрации: "+err.Error())
		return
	}
	response := models.NewRegisterResponse(user, "Пользователь создан. Для активации аккаунта необходимо подтвердить почту")
	respondWithJson(w, http.StatusCreated, response)
}

func (h *Handlers) Refresh(w http.ResponseWriter, r *http.Request) {
	var input models.RefreshInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный формат данных")
		return
	}

	if err := input.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.authService.Refresh(input)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Ошибка при выдаче новых токенов: "+err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, response)
}

func (h *Handlers) Verify(w http.ResponseWriter, r *http.Request) {
	var input models.VerifyInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный формат данных")
		return
	}

	if err := input.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	response, err := h.authService.Verify(input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, response)
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var input models.LoginUserInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некоррректный формат данных")
		return
	}

	if err := input.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.authService.Login(input)
	if err != nil {
		var notVerifiedErr *apperrors.ErrEmailNotVerified

		if errors.As(err, &notVerifiedErr) {
			response := models.AuthNotVerifyResponse{
				Email:   notVerifiedErr.Email,
				Message: notVerifiedErr.Error(),
			}
			respondWithJson(w, http.StatusForbidden, response)
			return
		}
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, response)
}

func (h *Handlers) GenerateNewEmailCode(w http.ResponseWriter, r *http.Request) {
	var input models.GenerateNewCodeInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный формат данных")
		return
	}

	if err := input.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	err := h.authService.NewEmailConfirmationCode(input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	var input models.RefreshInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный формат данных")
		return
	}

	if err := input.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.authService.Logout(input); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
