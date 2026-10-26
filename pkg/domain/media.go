package domain

import "time"

type UploadStatus string

const (
	UploadStatusPending    UploadStatus = "PENDING"
	UploadStatusUploading  UploadStatus = "UPLOADING"
	UploadStatusProcessing UploadStatus = "PROCESSING"
	UploadStatusCompleted  UploadStatus = "COMPLETED"
	UploadStatusFailed     UploadStatus = "FAILED"
)

type MediaFile struct {
	ID             string       `json:"id"`
	UploaderID     string       `json:"uploader_id"`
	FileName       string       `json:"file_name"`
	MimeType       string       `json:"mime_type"`
	ByteSize       int64        `json:"byte_size"`
	StorageURL     string       `json:"storage_url"`
	HLSManifestURL string       `json:"hls_manifest_url,omitempty"`
	EXIFSanitized  bool         `json:"exif_sanitized"`
	Status         UploadStatus `json:"status"`
	TotalChunks    int          `json:"total_chunks"`
	UploadedChunks []int        `json:"uploaded_chunks"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type Story struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"author_id"`
	MediaURL  string    `json:"media_url"`
	Caption   string    `json:"caption,omitempty"`
	ViewCount int       `json:"view_count"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Story) IsExpired(now time.Time) bool {
	return now.After(s.ExpiresAt)
}

type StoryView struct {
	StoryID  string    `json:"story_id"`
	ViewerID string    `json:"viewer_id"`
	ViewedAt time.Time `json:"viewed_at"`
}

type Reel struct {
	ID           string    `json:"id"`
	AuthorID     string    `json:"author_id"`
	VideoURL     string    `json:"video_url"`
	ThumbnailURL string    `json:"thumbnail_url"`
	AudioTrack   string    `json:"audio_track,omitempty"`
	Caption      string    `json:"caption,omitempty"`
	AspectRatio  string    `json:"aspect_ratio"` // Standard 9:16
	LikesCount   int       `json:"likes_count"`
	CreatedAt    time.Time `json:"created_at"`
}

type MediaAlbum struct {
	ID          string    `json:"id"`
	AuthorID    string    `json:"author_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	MediaIDs    []string  `json:"media_ids"`
	CreatedAt   time.Time `json:"created_at"`
}

type LiveStreamStatus string

const (
	LiveStatusCreated LiveStreamStatus = "CREATED"
	LiveStatusLive    LiveStreamStatus = "LIVE"
	LiveStatusEnded   LiveStreamStatus = "ENDED"
)

type LiveStream struct {
	ID          string           `json:"id"`
	HostID      string           `json:"host_id"`
	Title       string           `json:"title"`
	StreamKey   string           `json:"stream_key"`
	IngestURL   string           `json:"ingest_url"`   // rtmp://ingest.live/...
	PlaybackURL string           `json:"playback_url"` // https://cdn.live/.../index.m3u8
	Status      LiveStreamStatus `json:"status"`
	ViewerCount int              `json:"viewer_count"`
	StartedAt   *time.Time       `json:"started_at,omitempty"`
	EndedAt     *time.Time       `json:"ended_at,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}
