package domain

import "time"

// ============================================================================
// Groups
// ============================================================================

type GroupPrivacy string

const (
	GroupPrivacyPublic  GroupPrivacy = "PUBLIC"
	GroupPrivacyPrivate GroupPrivacy = "PRIVATE"
)

type GroupRole string

const (
	GroupRoleAdmin     GroupRole = "ADMIN"
	GroupRoleModerator GroupRole = "MODERATOR"
	GroupRoleMember    GroupRole = "MEMBER"
)

type GroupPostStatus string

const (
	GroupPostPending  GroupPostStatus = "PENDING_APPROVAL"
	GroupPostApproved GroupPostStatus = "APPROVED"
	GroupPostRejected GroupPostStatus = "REJECTED"
)

type GroupRule struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Group struct {
	ID                 string               `json:"id"`
	Name               string               `json:"name"`
	Description        string               `json:"description"`
	Privacy            GroupPrivacy         `json:"privacy"`
	CoverPhotoURL      string               `json:"cover_photo_url,omitempty"`
	CreatedBy          string               `json:"created_by"`
	Members            map[string]GroupRole `json:"members"` // userID -> role
	PendingMembers     []string             `json:"pending_members"`
	Rules              []GroupRule          `json:"rules,omitempty"`
	RequirePostApproval bool                `json:"require_post_approval"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
}

type GroupPost struct {
	ID        string          `json:"id"`
	GroupID   string          `json:"group_id"`
	AuthorID  string          `json:"author_id"`
	Content   string          `json:"content"`
	Status    GroupPostStatus `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
}

// ============================================================================
// Pages
// ============================================================================

type PageRole string

const (
	PageRoleAdmin     PageRole = "ADMIN"
	PageRoleEditor    PageRole = "EDITOR"
	PageRoleModerator PageRole = "MODERATOR"
)

type Page struct {
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	Category       string              `json:"category"`
	Description    string              `json:"description"`
	AvatarURL      string              `json:"avatar_url,omitempty"`
	CoverPhotoURL  string              `json:"cover_photo_url,omitempty"`
	Website        string              `json:"website,omitempty"`
	CreatedBy      string              `json:"created_by"`
	Roles          map[string]PageRole `json:"roles"` // userID -> PageRole
	FollowersCount int                 `json:"followers_count"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

type PageInsights struct {
	PageID         string  `json:"page_id"`
	FollowersCount int     `json:"followers_count"`
	TotalPostReach int     `json:"total_post_reach"`
	EngagementRate float64 `json:"engagement_rate"`
	Period         string  `json:"period"`
}

// ============================================================================
// Events
// ============================================================================

type RSVPStatus string

const (
	RSVPGoing      RSVPStatus = "GOING"
	RSVPInterested RSVPStatus = "INTERESTED"
	RSVPDeclined   RSVPStatus = "DECLINED"
)

type Event struct {
	ID           string                `json:"id"`
	HostID       string                `json:"host_id"`
	Title        string                `json:"title"`
	Description  string                `json:"description"`
	Location     string                `json:"location"`
	IsOnline     bool                  `json:"is_online"`
	StartsAt     time.Time             `json:"starts_at"`
	EndsAt       time.Time             `json:"ends_at"`
	RSVPs        map[string]RSVPStatus `json:"rsvps"` // userID -> status
	CreatedAt    time.Time             `json:"created_at"`
}
