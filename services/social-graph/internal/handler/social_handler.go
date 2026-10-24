package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dnjstjqppk-dev/facebook/pkg/auth"
	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/social-graph/internal/service"
)

type SocialHandler struct {
	socialSvc *service.SocialGraphService
	jwtMgr    *auth.JWTManager
}

func NewSocialHandler(socialSvc *service.SocialGraphService, jwtMgr *auth.JWTManager) *SocialHandler {
	return &SocialHandler{
		socialSvc: socialSvc,
		jwtMgr:    jwtMgr,
	}
}

func (h *SocialHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/friends/request", h.HandleSendFriendRequest)
	mux.HandleFunc("POST /api/v1/friends/accept", h.HandleAcceptFriendRequest)
	mux.HandleFunc("POST /api/v1/friends/reject", h.HandleRejectFriendRequest)
	mux.HandleFunc("POST /api/v1/friends/unfriend", h.HandleUnfriend)
	mux.HandleFunc("GET /api/v1/friends", h.HandleGetFriends)
	mux.HandleFunc("GET /api/v1/friends/pending", h.HandleGetPendingRequests)
	mux.HandleFunc("POST /api/v1/users/block", h.HandleBlockUser)
	mux.HandleFunc("GET /healthz", h.HandleHealthz)
}

func (h *SocialHandler) authenticate(r *http.Request) (*auth.JWTClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return nil, auth.ErrInvalidToken
	}
	tokenStr := authHeader[7:]
	return h.jwtMgr.VerifyToken(tokenStr)
}

func (h *SocialHandler) HandleSendFriendRequest(w http.ResponseWriter, r *http.Request) {
	claims, err := h.authenticate(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var action domain.FriendRequestAction
	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	rel, err := h.socialSvc.SendFriendRequest(r.Context(), claims.UserID, action.TargetUserID)
	if err != nil {
		status := http.StatusBadRequest
		if err == service.ErrSelfFriendship || err == service.ErrBlockedUser {
			status = http.StatusForbidden
		} else if err == service.ErrAlreadyFriends || err == service.ErrPendingRequestExists {
			status = http.StatusConflict
		}
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, rel)
}

func (h *SocialHandler) HandleAcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	claims, err := h.authenticate(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var action domain.FriendRequestAction
	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	rel, err := h.socialSvc.AcceptFriendRequest(r.Context(), claims.UserID, action.TargetUserID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, rel)
}

func (h *SocialHandler) HandleRejectFriendRequest(w http.ResponseWriter, r *http.Request) {
	claims, err := h.authenticate(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var action domain.FriendRequestAction
	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.socialSvc.RejectFriendRequest(r.Context(), claims.UserID, action.TargetUserID); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "request rejected"})
}

func (h *SocialHandler) HandleUnfriend(w http.ResponseWriter, r *http.Request) {
	claims, err := h.authenticate(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var action domain.FriendRequestAction
	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.socialSvc.Unfriend(r.Context(), claims.UserID, action.TargetUserID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "unfriended"})
}

func (h *SocialHandler) HandleGetFriends(w http.ResponseWriter, r *http.Request) {
	claims, err := h.authenticate(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	friends, err := h.socialSvc.GetFriends(r.Context(), claims.UserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"friends": friends,
		"count":   len(friends),
	})
}

func (h *SocialHandler) HandleGetPendingRequests(w http.ResponseWriter, r *http.Request) {
	claims, err := h.authenticate(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	pending, err := h.socialSvc.GetPendingRequests(r.Context(), claims.UserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"pending_requests": pending,
		"count":            len(pending),
	})
}

func (h *SocialHandler) HandleBlockUser(w http.ResponseWriter, r *http.Request) {
	claims, err := h.authenticate(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var action domain.FriendRequestAction
	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.socialSvc.BlockUser(r.Context(), claims.UserID, action.TargetUserID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "user blocked"})
}

func (h *SocialHandler) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "social-graph-service",
		"version": "1.0.0",
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
