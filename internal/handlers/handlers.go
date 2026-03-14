package handlers

import (
	"AuthService/internal/models"
	"AuthService/internal/service"
	"encoding/json"
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

	response, err := h.authService.Register(input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка при регистрации: "+err.Error())
		return
	}

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
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, response)
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
