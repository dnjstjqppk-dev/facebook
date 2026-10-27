package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/dnjstjqppk-dev/facebook/services/media/internal/service"
)

type MediaHandler struct {
	mediaService *service.MediaService
}

func NewMediaHandler(mediaService *service.MediaService) *MediaHandler {
	return &MediaHandler{mediaService: mediaService}
}

func (h *MediaHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/media/upload/init", h.HandleInitUpload)
	mux.HandleFunc("POST /api/v1/media/upload/{id}/chunk", h.HandleUploadChunk)
	mux.HandleFunc("POST /api/v1/media/upload/{id}/complete", h.HandleCompleteUpload)
	mux.HandleFunc("POST /api/v1/stories", h.HandleCreateStory)
	mux.HandleFunc("GET /api/v1/stories", h.HandleGetStories)
	mux.HandleFunc("POST /api/v1/stories/{id}/view", h.HandleViewStory)
	mux.HandleFunc("GET /api/v1/stories/{id}/viewers", h.HandleGetStoryViewers)
	mux.HandleFunc("POST /api/v1/reels", h.HandleCreateReel)
	mux.HandleFunc("GET /api/v1/reels", h.HandleListReels)
	mux.HandleFunc("POST /api/v1/live", h.HandleCreateLiveStream)
	mux.HandleFunc("POST /api/v1/live/{id}/start", h.HandleStartLiveStream)
	mux.HandleFunc("POST /api/v1/live/{id}/end", h.HandleEndLiveStream)
	mux.HandleFunc("GET /healthz", h.HandleHealthz)
}

func (h *MediaHandler) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "media"})
}

func (h *MediaHandler) HandleInitUpload(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req struct {
		FileName    string `json:"file_name"`
		MimeType    string `json:"mime_type"`
		ByteSize    int64  `json:"byte_size"`
		TotalChunks int    `json:"total_chunks"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	upload, err := h.mediaService.InitiateChunkedUpload(r.Context(), userID, req.FileName, req.MimeType, req.ByteSize, req.TotalChunks)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, upload)
}

func (h *MediaHandler) HandleUploadChunk(w http.ResponseWriter, r *http.Request) {
	uploadID := r.PathValue("id")
	chunkStr := r.URL.Query().Get("chunk")
	chunkNum, err := strconv.Atoi(chunkStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "valid chunk parameter required"})
		return
	}

	upload, err := h.mediaService.UploadChunk(r.Context(), uploadID, chunkNum)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, upload)
}

func (h *MediaHandler) HandleCompleteUpload(w http.ResponseWriter, r *http.Request) {
	uploadID := r.PathValue("id")
	upload, err := h.mediaService.CompleteUpload(r.Context(), uploadID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, upload)
}

func (h *MediaHandler) HandleCreateStory(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req struct {
		MediaURL string `json:"media_url"`
		Caption  string `json:"caption"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	story, err := h.mediaService.CreateStory(r.Context(), userID, req.MediaURL, req.Caption)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, story)
}

func (h *MediaHandler) HandleGetStories(w http.ResponseWriter, r *http.Request) {
	authorIDsStr := r.URL.Query().Get("author_ids")
	var authorIDs []string
	if authorIDsStr != "" {
		authorIDs = strings.Split(authorIDsStr, ",")
	}

	stories, err := h.mediaService.GetActiveStories(r.Context(), authorIDs)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, stories)
}

func (h *MediaHandler) HandleViewStory(w http.ResponseWriter, r *http.Request) {
	storyID := r.PathValue("id")
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	recorded, err := h.mediaService.ViewStory(r.Context(), storyID, userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"view_recorded": recorded})
}

func (h *MediaHandler) HandleGetStoryViewers(w http.ResponseWriter, r *http.Request) {
	storyID := r.PathValue("id")
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	viewers, err := h.mediaService.GetStoryViewers(r.Context(), storyID, userID)
	if err != nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, viewers)
}

func (h *MediaHandler) HandleCreateReel(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req struct {
		VideoURL     string `json:"video_url"`
		ThumbnailURL string `json:"thumbnail_url"`
		AudioTrack   string `json:"audio_track"`
		Caption      string `json:"caption"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	reel, err := h.mediaService.CreateReel(r.Context(), userID, req.VideoURL, req.ThumbnailURL, req.AudioTrack, req.Caption)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, reel)
}

func (h *MediaHandler) HandleListReels(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	reels, err := h.mediaService.ListReels(r.Context(), limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, reels)
}

func (h *MediaHandler) HandleCreateLiveStream(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	stream, err := h.mediaService.CreateLiveStream(r.Context(), userID, req.Title)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, stream)
}

func (h *MediaHandler) HandleStartLiveStream(w http.ResponseWriter, r *http.Request) {
	streamID := r.PathValue("id")
	userID := r.Header.Get("X-User-ID")
	stream, err := h.mediaService.StartLiveStream(r.Context(), streamID, userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, stream)
}

func (h *MediaHandler) HandleEndLiveStream(w http.ResponseWriter, r *http.Request) {
	streamID := r.PathValue("id")
	userID := r.Header.Get("X-User-ID")
	stream, err := h.mediaService.EndLiveStream(r.Context(), streamID, userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, stream)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

