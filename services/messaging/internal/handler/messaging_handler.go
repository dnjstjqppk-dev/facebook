package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/messaging/internal/service"
)

type MessagingHandler struct {
	svc *service.MessagingService
}

func NewMessagingHandler(svc *service.MessagingService) *MessagingHandler {
	return &MessagingHandler{svc: svc}
}

func (h *MessagingHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/conversations/direct", h.HandleCreateDirect)
	mux.HandleFunc("POST /api/v1/conversations/group", h.HandleCreateGroup)
	mux.HandleFunc("POST /api/v1/messages", h.HandleSendMessage)
	mux.HandleFunc("POST /api/v1/messages/{id}/unsend", h.HandleUnsendMessage)
	mux.HandleFunc("POST /api/v1/messages/{id}/reactions", h.HandleReactMessage)
	mux.HandleFunc("POST /api/v1/presence", h.HandleSetPresence)
	mux.HandleFunc("GET /api/v1/presence/{id}", h.HandleGetPresence)
	mux.HandleFunc("POST /api/v1/calls", h.HandleInitiateCall)
	mux.HandleFunc("POST /api/v1/calls/{id}/answer", h.HandleAnswerCall)
	mux.HandleFunc("POST /api/v1/calls/{id}/end", h.HandleEndCall)
	mux.HandleFunc("GET /healthz", h.HandleHealthz)
}

func (h *MessagingHandler) HandleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "messaging"})
}

func (h *MessagingHandler) HandleCreateDirect(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req struct {
		RecipientID string `json:"recipient_id"`
		IsE2EE      bool   `json:"is_e2ee"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	conv, err := h.svc.Create1on1Conversation(r.Context(), userID, req.RecipientID, req.IsE2EE)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, conv)
}

func (h *MessagingHandler) HandleCreateGroup(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req struct {
		Title     string   `json:"title"`
		MemberIDs []string `json:"member_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	conv, err := h.svc.CreateGroupConversation(r.Context(), userID, req.Title, req.MemberIDs)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, conv)
}

func (h *MessagingHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req struct {
		ConversationID string `json:"conversation_id"`
		Content        string `json:"content"`
		MediaURL       string `json:"media_url"`
		ReplyToID      string `json:"reply_to_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	msg, err := h.svc.SendMessage(r.Context(), userID, req.ConversationID, req.Content, req.MediaURL, req.ReplyToID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, msg)
}

func (h *MessagingHandler) HandleUnsendMessage(w http.ResponseWriter, r *http.Request) {
	msgID := r.PathValue("id")
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	err := h.svc.UnsendMessage(r.Context(), msgID, userID)
	if err != nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"unsent": true})
}

func (h *MessagingHandler) HandleReactMessage(w http.ResponseWriter, r *http.Request) {
	msgID := r.PathValue("id")
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req struct {
		Reaction domain.ReactionType `json:"reaction"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	err := h.svc.ReactToMessage(r.Context(), msgID, userID, req.Reaction)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *MessagingHandler) HandleSetPresence(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req struct {
		Status domain.PresenceStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	err := h.svc.SetPresence(r.Context(), userID, req.Status)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"presence_updated": true})
}

func (h *MessagingHandler) HandleGetPresence(w http.ResponseWriter, r *http.Request) {
	targetUserID := r.PathValue("id")
	pres, err := h.svc.GetPresence(r.Context(), targetUserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, pres)
}

func (h *MessagingHandler) HandleInitiateCall(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req struct {
		ReceiverIDs []string        `json:"receiver_ids"`
		CallType    domain.CallType `json:"call_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	call, err := h.svc.InitiateCall(r.Context(), userID, req.ReceiverIDs, req.CallType)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, call)
}

func (h *MessagingHandler) HandleAnswerCall(w http.ResponseWriter, r *http.Request) {
	callID := r.PathValue("id")
	userID := r.Header.Get("X-User-ID")
	call, err := h.svc.AnswerCall(r.Context(), callID, userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, call)
}

func (h *MessagingHandler) HandleEndCall(w http.ResponseWriter, r *http.Request) {
	callID := r.PathValue("id")
	userID := r.Header.Get("X-User-ID")
	call, err := h.svc.EndCall(r.Context(), callID, userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, call)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

