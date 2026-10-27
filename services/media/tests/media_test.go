package tests

import (
	"context"
	"testing"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/media/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/media/internal/service"
)

func setupTestMediaService() (*service.MediaService, repository.MediaRepository) {
	repo := repository.NewInMemoryMediaRepository()
	mediaSvc := service.NewMediaService(repo)
	return mediaSvc, repo
}

func TestChunkedUpload_FullCycle(t *testing.T) {
	mediaSvc, _ := setupTestMediaService()
	ctx := context.Background()

	// 1. Initiate 3-chunk upload for JPEG photo
	upload, err := mediaSvc.InitiateChunkedUpload(ctx, "user-photographer", "vacation.jpg", "image/jpeg", 15000000, 3)
	if err != nil {
		t.Fatalf("failed to initiate upload: %v", err)
	}
	if upload.Status != domain.UploadStatusUploading || upload.TotalChunks != 3 {
		t.Errorf("unexpected initial upload status: %v", upload)
	}

	// 2. Upload chunk 1 and 2
	_, err = mediaSvc.UploadChunk(ctx, upload.ID, 1)
	if err != nil {
		t.Fatalf("chunk 1 upload failed: %v", err)
	}
	_, err = mediaSvc.UploadChunk(ctx, upload.ID, 2)
	if err != nil {
		t.Fatalf("chunk 2 upload failed: %v", err)
	}

	// 3. Attempt to complete prematurely (should fail)
	_, err = mediaSvc.CompleteUpload(ctx, upload.ID)
	if err != service.ErrIncompleteChunks {
		t.Errorf("expected ErrIncompleteChunks, got: %v", err)
	}

	// 4. Upload chunk 3
	_, err = mediaSvc.UploadChunk(ctx, upload.ID, 3)
	if err != nil {
		t.Fatalf("chunk 3 upload failed: %v", err)
	}

	// 5. Complete upload
	completed, err := mediaSvc.CompleteUpload(ctx, upload.ID)
	if err != nil {
		t.Fatalf("complete upload failed: %v", err)
	}
	if completed.Status != domain.UploadStatusCompleted {
		t.Errorf("expected COMPLETED status, got: %s", completed.Status)
	}
	if !completed.EXIFSanitized {
		t.Errorf("expected EXIF metadata to be sanitized")
	}
}

func TestVideoUpload_HLSManifestGeneration(t *testing.T) {
	mediaSvc, _ := setupTestMediaService()
	ctx := context.Background()

	upload, err := mediaSvc.InitiateChunkedUpload(ctx, "user-filmmaker", "clip.mp4", "video/mp4", 50000000, 1)
	if err != nil {
		t.Fatalf("failed to initiate video upload: %v", err)
	}

	_, _ = mediaSvc.UploadChunk(ctx, upload.ID, 1)
	completed, err := mediaSvc.CompleteUpload(ctx, upload.ID)
	if err != nil {
		t.Fatalf("failed to complete video: %v", err)
	}

	if completed.HLSManifestURL == "" || completed.HLSManifestURL[:len("https://s3.storage.internal/hls/")] != "https://s3.storage.internal/hls/" {
		t.Errorf("expected valid HLS manifest URL for video, got: %s", completed.HLSManifestURL)
	}
}

func TestStory_CreationAnd24HourTTL(t *testing.T) {
	mediaSvc, repo := setupTestMediaService()
	ctx := context.Background()

	// 1. Create Story
	story, err := mediaSvc.CreateStory(ctx, "author-1", "https://cdn.example.com/story1.jpg", "Morning coffee!")
	if err != nil {
		t.Fatalf("failed to create story: %v", err)
	}

	// Verify 24-hour expiration duration
	diff := story.ExpiresAt.Sub(story.CreatedAt)
	if diff < 23*time.Hour || diff > 25*time.Hour {
		t.Errorf("expected 24-hour TTL, got: %v", diff)
	}

	// Verify active stories query includes this story
	active, err := mediaSvc.GetActiveStories(ctx, []string{"author-1"})
	if err != nil || len(active) != 1 {
		t.Fatalf("expected 1 active story, got: %d", len(active))
	}

	// 2. Test Expired Story
	expiredStory := &domain.Story{
		ID:        "expired-story-1",
		AuthorID:  "author-1",
		MediaURL:  "https://cdn.example.com/expired.jpg",
		ExpiresAt: time.Now().UTC().Add(-1 * time.Hour), // 1 hour in past
		CreatedAt: time.Now().UTC().Add(-25 * time.Hour),
	}
	_ = repo.SaveStory(ctx, expiredStory)

	// Expired story should be excluded from active stories
	active, _ = mediaSvc.GetActiveStories(ctx, []string{"author-1"})
	for _, s := range active {
		if s.ID == "expired-story-1" {
			t.Errorf("expired story should not be returned in active list")
		}
	}
}

func TestStory_UniqueViewsTracking(t *testing.T) {
	mediaSvc, repo := setupTestMediaService()
	ctx := context.Background()

	story, _ := mediaSvc.CreateStory(ctx, "author-1", "https://cdn.example.com/story.jpg", "Hello!")

	// Viewer 1 watches
	newView, err := mediaSvc.ViewStory(ctx, story.ID, "viewer-alice")
	if err != nil || !newView {
		t.Fatalf("first view should be recorded as new")
	}

	// Viewer 1 watches again (idempotent, should not increment view count)
	newView, err = mediaSvc.ViewStory(ctx, story.ID, "viewer-alice")
	if err != nil || newView {
		t.Fatalf("repeat view should not be recorded as new")
	}

	// Viewer 2 watches
	_, _ = mediaSvc.ViewStory(ctx, story.ID, "viewer-bob")

	st, _ := repo.GetStory(ctx, story.ID)
	if st.ViewCount != 2 {
		t.Errorf("expected view count 2, got: %d", st.ViewCount)
	}

	// Test Author access to viewer list
	viewers, err := mediaSvc.GetStoryViewers(ctx, story.ID, "author-1")
	if err != nil || len(viewers) != 2 {
		t.Fatalf("expected 2 viewers, got: %d", len(viewers))
	}

	// Non-author attempting to view viewers list should be rejected
	_, err = mediaSvc.GetStoryViewers(ctx, story.ID, "eavesdropper-99")
	if err != service.ErrUnauthorizedMedia {
		t.Errorf("expected ErrUnauthorizedMedia, got: %v", err)
	}
}

func TestReel_VerticalAspectRatioStandard(t *testing.T) {
	mediaSvc, _ := setupTestMediaService()
	ctx := context.Background()

	reel, err := mediaSvc.CreateReel(ctx, "creator-pro", "https://cdn.example.com/dance.mp4", "https://cdn.example.com/thumb.jpg", "Track: Upbeat Beats #12", "Dancing in the rain #shorts")
	if err != nil {
		t.Fatalf("failed to create reel: %v", err)
	}

	if reel.AspectRatio != "9:16" {
		t.Errorf("reel must adhere to 9:16 vertical standard, got %s", reel.AspectRatio)
	}
	if reel.AudioTrack != "Track: Upbeat Beats #12" {
		t.Errorf("audio track mismatch")
	}

	reelsList, err := mediaSvc.ListReels(ctx, 10)
	if err != nil || len(reelsList) != 1 {
		t.Fatalf("expected 1 reel in list, got: %d", len(reelsList))
	}
}

func TestPhotoAlbum_CreationAndManagement(t *testing.T) {
	mediaSvc, _ := setupTestMediaService()
	ctx := context.Background()

	album, err := mediaSvc.CreateAlbum(ctx, "author-1", "Summer Trip 2026", "Photos from Jeju", []string{"media-1", "media-2", "media-3"})
	if err != nil {
		t.Fatalf("failed to create album: %v", err)
	}

	if album.Name != "Summer Trip 2026" || len(album.MediaIDs) != 3 {
		t.Errorf("album attributes mismatch")
	}
}

func TestLiveStream_RTMPandHLSFlow(t *testing.T) {
	mediaSvc, _ := setupTestMediaService()
	ctx := context.Background()

	stream, err := mediaSvc.CreateLiveStream(ctx, "streamer-host", "Live Gaming Tournament Final")
	if err != nil {
		t.Fatalf("failed to create live stream: %v", err)
	}

	if stream.Status != domain.LiveStatusCreated {
		t.Errorf("expected status CREATED, got: %s", stream.Status)
	}
	if stream.StreamKey == "" || stream.IngestURL[:len("rtmp://")] != "rtmp://" {
		t.Errorf("invalid RTMP ingest URL or stream key")
	}

	// Start stream
	started, err := mediaSvc.StartLiveStream(ctx, stream.ID, "streamer-host")
	if err != nil || started.Status != domain.LiveStatusLive || started.StartedAt == nil {
		t.Fatalf("stream start failed: %v", err)
	}

	// End stream
	ended, err := mediaSvc.EndLiveStream(ctx, stream.ID, "streamer-host")
	if err != nil || ended.Status != domain.LiveStatusEnded || ended.EndedAt == nil {
		t.Fatalf("stream end failed: %v", err)
	}
}

func TestInvalidChunk_Rejected(t *testing.T) {
	mediaSvc, _ := setupTestMediaService()
	ctx := context.Background()

	upload, _ := mediaSvc.InitiateChunkedUpload(ctx, "user-1", "file.png", "image/png", 1000, 2)

	// Chunk 3 for totalChunks 2 should fail
	_, err := mediaSvc.UploadChunk(ctx, upload.ID, 3)
	if err != service.ErrInvalidChunk {
		t.Errorf("expected ErrInvalidChunk, got: %v", err)
	}
}

