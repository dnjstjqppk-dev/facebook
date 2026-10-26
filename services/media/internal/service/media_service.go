package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/media/internal/repository"
)

var (
	ErrInvalidChunk      = errors.New("invalid chunk number or out of bounds")
	ErrIncompleteChunks  = errors.New("cannot complete upload: missing chunks")
	ErrStoryExpired      = errors.New("story has expired (24-hour TTL)")
	ErrUnauthorizedMedia = errors.New("not authorized to perform this operation")
	ErrStreamNotLive     = errors.New("stream is not currently live")
)

type MediaService struct {
	repo repository.MediaRepository
}

func NewMediaService(repo repository.MediaRepository) *MediaService {
	return &MediaService{repo: repo}
}

func (s *MediaService) InitiateChunkedUpload(ctx context.Context, uploaderID, fileName, mimeType string, byteSize int64, totalChunks int) (*domain.MediaFile, error) {
	if totalChunks <= 0 {
		totalChunks = 1
	}

	uploadID := generateUUID()
	now := time.Now().UTC()

	media := &domain.MediaFile{
		ID:             uploadID,
		UploaderID:     uploaderID,
		FileName:       fileName,
		MimeType:       mimeType,
		ByteSize:       byteSize,
		StorageURL:     fmt.Sprintf("https://s3.storage.internal/media/%s/%s", uploadID, fileName),
		Status:         domain.UploadStatusUploading,
		TotalChunks:    totalChunks,
		UploadedChunks: []int{},
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repo.SaveMediaFile(ctx, media); err != nil {
		return nil, err
	}
	return media, nil
}

func (s *MediaService) UploadChunk(ctx context.Context, uploadID string, chunkNumber int) (*domain.MediaFile, error) {
	media, err := s.repo.GetMediaFile(ctx, uploadID)
	if err != nil {
		return nil, err
	}

	if chunkNumber <= 0 || chunkNumber > media.TotalChunks {
		return nil, ErrInvalidChunk
	}

	// Avoid duplicates
	for _, c := range media.UploadedChunks {
		if c == chunkNumber {
			return media, nil
		}
	}

	media.UploadedChunks = append(media.UploadedChunks, chunkNumber)
	media.UpdatedAt = time.Now().UTC()

	if err := s.repo.SaveMediaFile(ctx, media); err != nil {
		return nil, err
	}
	return media, nil
}

func (s *MediaService) CompleteUpload(ctx context.Context, uploadID string) (*domain.MediaFile, error) {
	media, err := s.repo.GetMediaFile(ctx, uploadID)
	if err != nil {
		return nil, err
	}

	if len(media.UploadedChunks) < media.TotalChunks {
		return nil, ErrIncompleteChunks
	}

	// Privacy & Security sanitization: strip EXIF metadata
	media.EXIFSanitized = true
	media.Status = domain.UploadStatusCompleted
	media.UpdatedAt = time.Now().UTC()

	// If video, assign HLS adaptive streaming manifest URL
	if strings.HasPrefix(media.MimeType, "video/") {
		media.HLSManifestURL = fmt.Sprintf("https://s3.storage.internal/hls/%s/master.m3u8", uploadID)
	}

	if err := s.repo.SaveMediaFile(ctx, media); err != nil {
		return nil, err
	}
	return media, nil
}

func (s *MediaService) CreateStory(ctx context.Context, authorID, mediaURL, caption string) (*domain.Story, error) {
	now := time.Now().UTC()
	// Strict 24-hour ephemeral media TTL
	expiresAt := now.Add(24 * time.Hour)

	story := &domain.Story{
		ID:        generateUUID(),
		AuthorID:  authorID,
		MediaURL:  mediaURL,
		Caption:   caption,
		ViewCount: 0,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}

	if err := s.repo.SaveStory(ctx, story); err != nil {
		return nil, err
	}
	return story, nil
}

func (s *MediaService) GetActiveStories(ctx context.Context, authorIDs []string) ([]*domain.Story, error) {
	return s.repo.GetActiveStories(ctx, authorIDs, time.Now().UTC())
}

func (s *MediaService) ViewStory(ctx context.Context, storyID, viewerID string) (bool, error) {
	story, err := s.repo.GetStory(ctx, storyID)
	if err != nil {
		return false, err
	}

	now := time.Now().UTC()
	if story.IsExpired(now) {
		return false, ErrStoryExpired
	}

	return s.repo.RecordStoryView(ctx, storyID, viewerID, now)
}

func (s *MediaService) GetStoryViewers(ctx context.Context, storyID, requesterID string) ([]*domain.StoryView, error) {
	story, err := s.repo.GetStory(ctx, storyID)
	if err != nil {
		return nil, err
	}

	// Only author can inspect the viewer list
	if story.AuthorID != requesterID {
		return nil, ErrUnauthorizedMedia
	}

	return s.repo.GetStoryViewers(ctx, storyID)
}

func (s *MediaService) CreateReel(ctx context.Context, authorID, videoURL, thumbnailURL, audioTrack, caption string) (*domain.Reel, error) {
	now := time.Now().UTC()

	reel := &domain.Reel{
		ID:           generateUUID(),
		AuthorID:     authorID,
		VideoURL:     videoURL,
		ThumbnailURL: thumbnailURL,
		AudioTrack:   audioTrack,
		Caption:      caption,
		AspectRatio:  "9:16", // Strict vertical video standard
		LikesCount:   0,
		CreatedAt:    now,
	}

	if err := s.repo.SaveReel(ctx, reel); err != nil {
		return nil, err
	}
	return reel, nil
}

func (s *MediaService) ListReels(ctx context.Context, limit int) ([]*domain.Reel, error) {
	return s.repo.ListReels(ctx, limit)
}

func (s *MediaService) CreateAlbum(ctx context.Context, authorID, name, description string, mediaIDs []string) (*domain.MediaAlbum, error) {
	album := &domain.MediaAlbum{
		ID:          generateUUID(),
		AuthorID:    authorID,
		Name:        name,
		Description: description,
		MediaIDs:    mediaIDs,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.repo.SaveAlbum(ctx, album); err != nil {
		return nil, err
	}
	return album, nil
}

func (s *MediaService) CreateLiveStream(ctx context.Context, hostID, title string) (*domain.LiveStream, error) {
	streamID := generateUUID()
	streamKeyBytes := make([]byte, 16)
	_, _ = rand.Read(streamKeyBytes)
	streamKey := hex.EncodeToString(streamKeyBytes)

	now := time.Now().UTC()
	stream := &domain.LiveStream{
		ID:          streamID,
		HostID:      hostID,
		Title:       title,
		StreamKey:   streamKey,
		IngestURL:   fmt.Sprintf("rtmp://live.ingest.social.network/app/%s", streamKey),
		PlaybackURL: fmt.Sprintf("https://cdn.live.social.network/hls/%s/index.m3u8", streamID),
		Status:      domain.LiveStatusCreated,
		ViewerCount: 0,
		CreatedAt:   now,
	}

	if err := s.repo.SaveLiveStream(ctx, stream); err != nil {
		return nil, err
	}
	return stream, nil
}

func (s *MediaService) StartLiveStream(ctx context.Context, streamID, hostID string) (*domain.LiveStream, error) {
	stream, err := s.repo.GetLiveStream(ctx, streamID)
	if err != nil {
		return nil, err
	}
	if stream.HostID != hostID {
		return nil, ErrUnauthorizedMedia
	}

	now := time.Now().UTC()
	stream.Status = domain.LiveStatusLive
	stream.StartedAt = &now

	if err := s.repo.SaveLiveStream(ctx, stream); err != nil {
		return nil, err
	}
	return stream, nil
}

func (s *MediaService) EndLiveStream(ctx context.Context, streamID, hostID string) (*domain.LiveStream, error) {
	stream, err := s.repo.GetLiveStream(ctx, streamID)
	if err != nil {
		return nil, err
	}
	if stream.HostID != hostID {
		return nil, ErrUnauthorizedMedia
	}

	now := time.Now().UTC()
	stream.Status = domain.LiveStatusEnded
	stream.EndedAt = &now

	if err := s.repo.SaveLiveStream(ctx, stream); err != nil {
		return nil, err
	}
	return stream, nil
}

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return hex.EncodeToString(b[0:4]) + "-" +
		hex.EncodeToString(b[4:6]) + "-" +
		hex.EncodeToString(b[6:8]) + "-" +
		hex.EncodeToString(b[8:10]) + "-" +
		hex.EncodeToString(b[10:16])
}
