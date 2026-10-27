package domain

import "time"

type PrivacyLevel string

const (
	PrivacyPublic  PrivacyLevel = "PUBLIC"
	PrivacyFriends PrivacyLevel = "FRIENDS"
	PrivacyOnlyMe  PrivacyLevel = "ONLY_ME"
)

type ReactionType string

const (
	ReactionLike  ReactionType = "LIKE"
	ReactionLove  ReactionType = "LOVE"
	ReactionCare  ReactionType = "CARE"
	ReactionHaha  ReactionType = "HAHA"
	ReactionWow   ReactionType = "WOW"
	ReactionSad   ReactionType = "SAD"
	ReactionAngry ReactionType = "ANGRY"
)

type EditHistoryEntry struct {
	PreviousContent string    `json:"previous_content"`
	EditedAt        time.Time `json:"edited_at"`
}

type Post struct {
	ID             string             `json:"id"`
	AuthorID       string             `json:"author_id"`
	Content        string             `json:"content"`
	MediaURLs      []string           `json:"media_urls,omitempty"`
	Feeling        string             `json:"feeling,omitempty"`
	Location       string             `json:"location,omitempty"`
	TaggedUserIDs  []string           `json:"tagged_user_ids,omitempty"`
	Privacy        PrivacyLevel       `json:"privacy"`
	OriginalPostID string             `json:"original_post_id,omitempty"` // For re-shares
	IsPinned       bool               `json:"is_pinned"`
	ReactionCounts map[ReactionType]int `json:"reaction_counts"`
	CommentCount   int                `json:"comment_count"`
	ShareCount     int                `json:"share_count"`
	EditHistory    []EditHistoryEntry `json:"edit_history,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	DeletedAt      *time.Time         `json:"deleted_at,omitempty"`
}

type Reaction struct {
	PostID       string       `json:"post_id"`
	UserID       string       `json:"user_id"`
	ReactionType ReactionType `json:"reaction_type"`
	CreatedAt    time.Time    `json:"created_at"`
}

type Comment struct {
	ID              string     `json:"id"`
	PostID          string     `json:"post_id"`
	AuthorID        string     `json:"author_id"`
	ParentCommentID string     `json:"parent_comment_id,omitempty"` // Empty if root comment
	Content         string     `json:"content"`
	Mentions        []string   `json:"mentions,omitempty"`
	ReactionCounts  map[ReactionType]int `json:"reaction_counts"`
	RepliesCount    int        `json:"replies_count"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type SavedCollection struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	PostIDs   []string  `json:"post_ids"`
	CreatedAt time.Time `json:"created_at"`
}

type CreatePostRequest struct {
	Content        string       `json:"content"`
	MediaURLs      []string     `json:"media_urls,omitempty"`
	Feeling        string       `json:"feeling,omitempty"`
	Location       string       `json:"location,omitempty"`
	TaggedUserIDs  []string     `json:"tagged_user_ids,omitempty"`
	Privacy        PrivacyLevel `json:"privacy"`
	OriginalPostID string       `json:"original_post_id,omitempty"`
}

type UpdatePostRequest struct {
	Content   string       `json:"content"`
	Privacy   PrivacyLevel `json:"privacy"`
	MediaURLs []string     `json:"media_urls,omitempty"`
}

type AddReactionRequest struct {
	ReactionType ReactionType `json:"reaction_type"`
}

type AddCommentRequest struct {
	ParentCommentID string   `json:"parent_comment_id,omitempty"`
	Content         string   `json:"content"`
	Mentions        []string `json:"mentions,omitempty"`
}

type FeedItem struct {
	Post         *Post        `json:"post"`
	UserReaction ReactionType `json:"user_reaction,omitempty"` // Current user's reaction if any
	RankingScore float64      `json:"ranking_score,omitempty"`
}

type FeedResponse struct {
	Items      []*FeedItem `json:"items"`
	NextCursor string      `json:"next_cursor,omitempty"`
	HasMore    bool        `json:"has_more"`
}

