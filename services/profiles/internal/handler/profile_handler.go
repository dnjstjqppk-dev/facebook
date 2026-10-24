package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dnjstjqppk-dev/facebook/pkg/auth"
	"github.com/dnjstjqppk-dev/facebook/services/profiles/internal/service"
)

type ProfileHandler struct {
	profileSvc *service.ProfileService
	jwtMgr     *auth.JWTManager
}

func NewProfileHandler(profileSvc *service.ProfileService, jwtMgr *auth.JWTManager) *ProfileHandler {
	return &ProfileHandler{
		profileSvc: profileSvc,
		jwtMgr:     jwtMgr,
	}
}

func (h *ProfileHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/profiles/", h.HandleGetProfile)
	mux.HandleFunc("PUT /api/v1/profiles/me", h.HandleUpdateProfile)
	mux.HandleFunc("GET /healthz", h.HandleHealthz)
}

func (h *ProfileHandler) authenticateOptional(r *http.Request) *auth.JWTClaims {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return nil
	}
	claims, err := h.jwtMgr.VerifyToken(authHeader[7:])
	if err != nil {
		return nil
	}
	return claims
}

func (h *ProfileHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/api/v1/profiles/")
	if username == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username required"})
		return
	}

	viewerID := ""
	if claims := h.authenticateOptional(r); claims != nil {
		viewerID = claims.UserID
	}

	p, err := h.profileSvc.GetProfile(r.Context(), viewerID, username)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "profile not found"})
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func (h *ProfileHandler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := h.authenticateOptional(r)
	if claims == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req service.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	updated, err := h.profileSvc.UpdateProfile(r.Context(), claims.UserID, &req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *ProfileHandler) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "profiles-service",
		"version": "1.0.0",
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
