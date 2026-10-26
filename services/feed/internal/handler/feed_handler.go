package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/feed/internal/service"
)

type FeedHandler struct {
	feedService *service.FeedService
}

func NewFeedHandler(feedService *service.FeedService) *FeedHandler {
	return &FeedHandler{feedService: feedService}
}

func (h *FeedHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/posts", h.HandleCreatePost)
	mux.HandleFunc("GET /api/v1/posts/", h.HandlePostOperations)
	mux.HandleFunc("PATCH /api/v1/posts/", h.HandlePostOperations)
	mux.HandleFunc("DELETE /api/v1/posts/", h.HandlePostOperations)
	mux.HandleFunc("POST /api/v1/posts/{id}/reactions", h.HandleAddReaction)
	mux.HandleFunc("DELETE /api/v1/posts/{id}/reactions", h.HandleDeleteReaction)
	mux.HandleFunc("POST /api/v1/posts/{id}/comments", h.HandleAddComment)
	mux.HandleFunc("GET /api/v1/posts/{id}/comments", h.HandleGetComments)
	mux.HandleFunc("GET /api/v1/feed", h.HandleGetFeed)
	mux.HandleFunc("GET /api/v1/memories", h.HandleGetMemories)
	mux.HandleFunc("POST /api/v1/saved", h.HandleSavePost)
	mux.HandleFunc("GET /api/v1/saved", h.HandleGetSaved)
	mux.HandleFunc("GET /healthz", h.HandleHealthz)
}

func (h *FeedHandler) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "feed"})
}

func (h *FeedHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req domain.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	post, err := h.feedService.CreatePost(r.Context(), userID, &req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, post)
}

func (h *FeedHandler) HandlePostOperations(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/posts/")
	parts := strings.Split(path, "/")
	postID := parts[0]
	if postID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "post ID required"})
		return
	}

	// Sub-route handling
	if len(parts) > 1 {
		sub := parts[1]
		if sub == "reactions" {
			if r.Method == http.MethodPost {
				h.handleAddReactionWithID(w, r, postID)
			} else if r.Method == http.MethodDelete {
				h.handleDeleteReactionWithID(w, r, postID)
			}
			return
		}
		if sub == "comments" {
			if r.Method == http.MethodPost {
				h.handleAddCommentWithID(w, r, postID)
			} else if r.Method == http.MethodGet {
				h.handleGetCommentsWithID(w, r, postID)
			}
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		// Not implemented directly here, can be added if needed
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "direct get by id via feed"})
	case http.MethodPatch:
		userID := r.Header.Get("X-User-ID")
		var req domain.UpdatePostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
			return
		}
		updated, err := h.feedService.UpdatePost(r.Context(), userID, postID, &req)
		if err != nil {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, updated)
	case http.MethodDelete:
		userID := r.Header.Get("X-User-ID")
		err := h.feedService.DeletePost(r.Context(), userID, postID)
		if err != nil {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
	}
}

func (h *FeedHandler) HandleAddReaction(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	h.handleAddReactionWithID(w, r, postID)
}

func (h *FeedHandler) handleAddReactionWithID(w http.ResponseWriter, r *http.Request, postID string) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req domain.AddReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	err := h.feedService.ReactToPost(r.Context(), userID, postID, req.ReactionType)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *FeedHandler) HandleDeleteReaction(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	h.handleDeleteReactionWithID(w, r, postID)
}

func (h *FeedHandler) handleDeleteReactionWithID(w http.ResponseWriter, r *http.Request, postID string) {
	userID := r.Header.Get("X-User-ID")
	err := h.feedService.RemoveReaction(r.Context(), userID, postID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"removed": true})
}

func (h *FeedHandler) HandleAddComment(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	h.handleAddCommentWithID(w, r, postID)
}

func (h *FeedHandler) handleAddCommentWithID(w http.ResponseWriter, r *http.Request, postID string) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req domain.AddCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	comment, err := h.feedService.AddComment(r.Context(), userID, postID, &req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, comment)
}

func (h *FeedHandler) HandleGetComments(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	h.handleGetCommentsWithID(w, r, postID)
}

func (h *FeedHandler) handleGetCommentsWithID(w http.ResponseWriter, r *http.Request, postID string) {
	comments, err := h.feedService.GetComments(r.Context(), postID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, comments)
}

func (h *FeedHandler) HandleGetFeed(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	cursor := r.URL.Query().Get("cursor")
	mode := r.URL.Query().Get("mode") // chronological or ranked
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	res, err := h.feedService.GetFeed(r.Context(), userID, cursor, limit, mode)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (h *FeedHandler) HandleGetMemories(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	memories, err := h.feedService.GetMemories(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, memories)
}

func (h *FeedHandler) HandleSavePost(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var payload struct {
		PostID         string `json:"post_id"`
		CollectionName string `json:"collection_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	err := h.feedService.SavePost(r.Context(), userID, payload.CollectionName, payload.PostID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"saved": true})
}

func (h *FeedHandler) HandleGetSaved(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	collections, err := h.feedService.GetSavedCollections(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, collections)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
