package domain

import "time"

type RelationshipType string

const (
	RelFriend   RelationshipType = "FRIEND"
	RelFollow   RelationshipType = "FOLLOW"
	RelBlock    RelationshipType = "BLOCK"
	RelRestrict RelationshipType = "RESTRICT"
)

type RelationshipStatus string

const (
	RelStatusPending  RelationshipStatus = "PENDING"
	RelStatusAccepted RelationshipStatus = "ACCEPTED"
	RelStatusDeclined RelationshipStatus = "DECLINED"
)

type Relationship struct {
	ID               string             `json:"id"`
	UserID           string             `json:"user_id"`
	TargetUserID     string             `json:"target_user_id"`
	RelationshipType RelationshipType   `json:"relationship_type"`
	Status           RelationshipStatus `json:"status"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

type FriendRequestAction struct {
	TargetUserID string `json:"target_user_id"`
}

type FriendResponse struct {
	User     *User    `json:"user"`
	Profile  *Profile `json:"profile"`
	Since    time.Time `json:"since"`
}

type PendingRequestResponse struct {
	RequestID string    `json:"request_id"`
	Requester *Profile  `json:"requester"`
	CreatedAt time.Time `json:"created_at"`
}
