package domain

import "time"

type ConversationType string

const (
	ConvOneToOne         ConversationType = "ONE_TO_ONE"
	ConvGroup            ConversationType = "GROUP"
	ConvBroadcastChannel ConversationType = "BROADCAST_CHANNEL"
)

type MemberRole string

const (
	MemberRoleAdmin  MemberRole = "ADMIN"
	MemberRoleMember MemberRole = "MEMBER"
)

type Conversation struct {
	ID                 string                  `json:"id"`
	Type               ConversationType        `json:"type"`
	Title              string                  `json:"title,omitempty"`
	AvatarURL          string                  `json:"avatar_url,omitempty"`
	CreatorID          string                  `json:"creator_id"`
	Members            map[string]MemberRole   `json:"members"` // userID -> role
	LastMessageSnippet string                  `json:"last_message_snippet,omitempty"`
	LastMessageAt      time.Time               `json:"last_message_at"`
	IsE2EE             bool                    `json:"is_e2ee"`
	CreatedAt          time.Time               `json:"created_at"`
	UpdatedAt          time.Time               `json:"updated_at"`
}

type Message struct {
	ID               string                      `json:"id"`
	ConversationID   string                      `json:"conversation_id"`
	SenderID         string                      `json:"sender_id"`
	Content          string                      `json:"content"`
	MediaURL         string                      `json:"media_url,omitempty"`
	ReplyToMessageID string                      `json:"reply_to_message_id,omitempty"`
	IsUnsent         bool                        `json:"is_unsent"`
	IsE2EE           bool                        `json:"is_e2ee"`
	Reactions        map[string]ReactionType     `json:"reactions"` // userID -> reaction
	DeliveredTo      map[string]time.Time        `json:"delivered_to"`
	ReadBy           map[string]time.Time        `json:"read_by"`
	CreatedAt        time.Time                   `json:"created_at"`
}

type PresenceStatus string

const (
	PresenceOnline  PresenceStatus = "ONLINE"
	PresenceOffline PresenceStatus = "OFFLINE"
	PresenceAway    PresenceStatus = "AWAY"
)

type UserPresence struct {
	UserID     string         `json:"user_id"`
	Status     PresenceStatus `json:"status"`
	LastSeenAt time.Time      `json:"last_seen_at"`
}

type TypingEvent struct {
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	IsTyping       bool      `json:"is_typing"`
	Timestamp      time.Time `json:"timestamp"`
}

type CallType string

const (
	CallAudio CallType = "AUDIO"
	CallVideo CallType = "VIDEO"
)

type CallStatus string

const (
	CallInitiated CallStatus = "INITIATED"
	CallConnected CallStatus = "CONNECTED"
	CallEnded     CallStatus = "ENDED"
	CallRejected  CallStatus = "REJECTED"
)

type CallSession struct {
	ID          string     `json:"id"`
	CallerID    string     `json:"caller_id"`
	ReceiverIDs []string   `json:"receiver_ids"`
	CallType    CallType   `json:"call_type"`
	Status      CallStatus `json:"status"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	EndedAt     *time.Time `json:"ended_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

