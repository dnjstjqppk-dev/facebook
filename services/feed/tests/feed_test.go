package tests

import (
	"context"
	"testing"
	"time"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/feed/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/feed/internal/service"
)

type mockFriendLister struct {
	friends map[string][]string
}

func (m *mockFriendLister) GetFriendIDs(ctx context.Context, userID string) ([]string, error) {
	return m.friends[userID], nil
}

func setupTestFeedService() (*service.FeedService, *repository.InMemoryFeedRepository, *mockFriendLister) {
	repo := repository.NewInMemoryFeedRepository()
	friendLister := &mockFriendLister{
		friends: make(map[string][]string),
	}
	feedSvc := service.NewFeedService(repo, friendLister)
	return feedSvc, repo, friendLister
}

func TestCreatePost_TextAndMedia(t *testing.T) {
	feedSvc, _, _ := setupTestFeedService()
	ctx := context.Background()

	req := &domain.CreatePostRequest{
		Content:       "Hello Facebook Parity World!",
		MediaURLs:     []string{"https://cdn.example.com/photo1.jpg"},
		Feeling:       "excited",
		Location:      "Seoul, South Korea",
		TaggedUserIDs: []string{"user-2", "user-3"},
		Privacy:       domain.PrivacyPublic,
	}

	post, err := feedSvc.CreatePost(ctx, "user-1", req)
	if err != nil {
		t.Fatalf("expected post creation to succeed, got: %v", err)
	}

	if post.AuthorID != "user-1" {
		t.Errorf("expected author user-1, got: %s", post.AuthorID)
	}
	if len(post.MediaURLs) != 1 || post.MediaURLs[0] != "https://cdn.example.com/photo1.jpg" {
		t.Errorf("media URLs mismatch")
	}
	if post.Feeling != "excited" || post.Location != "Seoul, South Korea" {
		t.Errorf("metadata mismatch")
	}
}

func TestEditPost_WithHistory(t *testing.T) {
	feedSvc, _, _ := setupTestFeedService()
	ctx := context.Background()

	post, err := feedSvc.CreatePost(ctx, "author-1", &domain.CreatePostRequest{
		Content: "Original Version",
		Privacy: domain.PrivacyFriends,
	})
	if err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	// Update Post
	updated, err := feedSvc.UpdatePost(ctx, "author-1", post.ID, &domain.UpdatePostRequest{
		Content: "Edited Version",
		Privacy: domain.PrivacyPublic,
	})
	if err != nil {
		t.Fatalf("update post failed: %v", err)
	}

	if updated.Content != "Edited Version" {
		t.Errorf("content not updated")
	}
	if len(updated.EditHistory) != 1 {
		t.Fatalf("expected 1 edit history entry, got %d", len(updated.EditHistory))
	}
	if updated.EditHistory[0].PreviousContent != "Original Version" {
		t.Errorf("previous content mismatch in edit history")
	}

	// Unauthorized edit attempt
	_, err = feedSvc.UpdatePost(ctx, "attacker-99", post.ID, &domain.UpdatePostRequest{
		Content: "Hacked Content",
	})
	if err != service.ErrUnauthorizedPost {
		t.Errorf("expected ErrUnauthorizedPost, got: %v", err)
	}
}

func TestDeletePost(t *testing.T) {
	feedSvc, _, _ := setupTestFeedService()
	ctx := context.Background()

	post, _ := feedSvc.CreatePost(ctx, "user-1", &domain.CreatePostRequest{
		Content: "To be deleted",
	})

	err := feedSvc.DeletePost(ctx, "user-1", post.ID)
	if err != nil {
		t.Fatalf("failed to delete post: %v", err)
	}

	// Should not be visible in feed
	feed, err := feedSvc.GetFeed(ctx, "user-1", "", 10, "chronological")
	if err != nil {
		t.Fatalf("failed to get feed: %v", err)
	}
	for _, item := range feed.Items {
		if item.Post.ID == post.ID {
			t.Errorf("deleted post should not be returned in feed")
		}
	}
}

func TestFanOutOnWrite_Timeline(t *testing.T) {
	feedSvc, _, friends := setupTestFeedService()
	ctx := context.Background()

	friends.friends["user-alice"] = []string{"user-bob", "user-charlie"}

	post, err := feedSvc.CreatePost(ctx, "user-alice", &domain.CreatePostRequest{
		Content: "Friends-only update",
		Privacy: domain.PrivacyFriends,
	})
	if err != nil {
		t.Fatalf("create post failed: %v", err)
	}

	// Verify Bob's feed has Alice's post
	feedBob, err := feedSvc.GetFeed(ctx, "user-bob", "", 10, "chronological")
	if err != nil || len(feedBob.Items) != 1 {
		t.Fatalf("expected 1 post in Bob's feed, got: %v", len(feedBob.Items))
	}
	if feedBob.Items[0].Post.ID != post.ID {
		t.Errorf("expected post ID %s, got: %s", post.ID, feedBob.Items[0].Post.ID)
	}

	// Verify Charlie's feed has Alice's post
	feedCharlie, err := feedSvc.GetFeed(ctx, "user-charlie", "", 10, "chronological")
	if err != nil || len(feedCharlie.Items) != 1 {
		t.Fatalf("expected 1 post in Charlie's feed, got: %v", len(feedCharlie.Items))
	}
}

func TestSixReactions_Idempotent(t *testing.T) {
	feedSvc, repo, _ := setupTestFeedService()
	ctx := context.Background()

	post, _ := feedSvc.CreatePost(ctx, "author-1", &domain.CreatePostRequest{
		Content: "Reaction post",
	})

	// React LIKE
	err := feedSvc.ReactToPost(ctx, "user-1", post.ID, domain.ReactionLike)
	if err != nil {
		t.Fatalf("failed to react LIKE: %v", err)
	}

	p, _ := repo.GetPost(ctx, post.ID)
	if p.ReactionCounts[domain.ReactionLike] != 1 {
		t.Errorf("expected 1 LIKE, got: %d", p.ReactionCounts[domain.ReactionLike])
	}

	// Change reaction to LOVE (should decrement LIKE, increment LOVE)
	err = feedSvc.ReactToPost(ctx, "user-1", post.ID, domain.ReactionLove)
	if err != nil {
		t.Fatalf("failed to change reaction to LOVE: %v", err)
	}

	p, _ = repo.GetPost(ctx, post.ID)
	if p.ReactionCounts[domain.ReactionLike] != 0 {
		t.Errorf("expected 0 LIKE, got: %d", p.ReactionCounts[domain.ReactionLike])
	}
	if p.ReactionCounts[domain.ReactionLove] != 1 {
		t.Errorf("expected 1 LOVE, got: %d", p.ReactionCounts[domain.ReactionLove])
	}

	// Remove reaction
	err = feedSvc.RemoveReaction(ctx, "user-1", post.ID)
	if err != nil {
		t.Fatalf("failed to remove reaction: %v", err)
	}

	p, _ = repo.GetPost(ctx, post.ID)
	if p.ReactionCounts[domain.ReactionLove] != 0 {
		t.Errorf("expected 0 LOVE after removal, got: %d", p.ReactionCounts[domain.ReactionLove])
	}
}

func TestTwoTierNestedComments(t *testing.T) {
	feedSvc, _, _ := setupTestFeedService()
	ctx := context.Background()

	post, _ := feedSvc.CreatePost(ctx, "author-1", &domain.CreatePostRequest{
		Content: "Post with comments",
	})

	// 1. Root comment (Tier 1)
	rootComment, err := feedSvc.AddComment(ctx, "user-1", post.ID, &domain.AddCommentRequest{
		Content: "This is a root comment",
	})
	if err != nil {
		t.Fatalf("add root comment failed: %v", err)
	}
	if rootComment.ParentCommentID != "" {
		t.Errorf("root comment should have empty parent ID")
	}

	// 2. Nested reply (Tier 2)
	reply, err := feedSvc.AddComment(ctx, "user-2", post.ID, &domain.AddCommentRequest{
		ParentCommentID: rootComment.ID,
		Content:         "This is a nested reply",
	})
	if err != nil {
		t.Fatalf("add reply failed: %v", err)
	}
	if reply.ParentCommentID != rootComment.ID {
		t.Errorf("reply should point to root comment")
	}

	// 3. Attempting 3rd tier nesting should collapse to Tier 2 (root comment ID)
	deepReply, err := feedSvc.AddComment(ctx, "user-3", post.ID, &domain.AddCommentRequest{
		ParentCommentID: reply.ID,
		Content:         "Deep reply attempting tier 3",
	})
	if err != nil {
		t.Fatalf("deep reply failed: %v", err)
	}
	if deepReply.ParentCommentID != rootComment.ID {
		t.Errorf("3rd tier nesting must collapse to root comment ID %s, got %s", rootComment.ID, deepReply.ParentCommentID)
	}

	// Verify comments list
	comments, err := feedSvc.GetComments(ctx, post.ID)
	if err != nil || len(comments) != 3 {
		t.Fatalf("expected 3 total comments, got %d", len(comments))
	}
}

func TestCursorBasedPagination(t *testing.T) {
	feedSvc, _, _ := setupTestFeedService()
	ctx := context.Background()

	// Create 5 posts
	for i := 1; i <= 5; i++ {
		_, _ = feedSvc.CreatePost(ctx, "user-author", &domain.CreatePostRequest{
			Content: "Post number " + string(rune('0'+i)),
		})
		time.Sleep(2 * time.Millisecond) // distinct timestamps
	}

	// Page 1: limit 2
	page1, err := feedSvc.GetFeed(ctx, "user-author", "", 2, "chronological")
	if err != nil || len(page1.Items) != 2 {
		t.Fatalf("expected 2 items on page 1, got %d", len(page1.Items))
	}
	if !page1.HasMore || page1.NextCursor == "" {
		t.Fatalf("expected HasMore to be true and NextCursor present")
	}

	// Page 2: with NextCursor
	page2, err := feedSvc.GetFeed(ctx, "user-author", page1.NextCursor, 2, "chronological")
	if err != nil || len(page2.Items) != 2 {
		t.Fatalf("expected 2 items on page 2, got %d", len(page2.Items))
	}

	// Deduplication verify: page 1 items must differ from page 2 items
	if page1.Items[0].Post.ID == page2.Items[0].Post.ID || page1.Items[1].Post.ID == page2.Items[0].Post.ID {
		t.Errorf("pagination duplicate detected between pages")
	}
}

func TestRankedFeed_EdgeRankAlgorithm(t *testing.T) {
	feedSvc, _, _ := setupTestFeedService()
	ctx := context.Background()

	// Post 1: Low engagement
	p1, _ := feedSvc.CreatePost(ctx, "user-author", &domain.CreatePostRequest{
		Content: "Ordinary post",
	})

	// Post 2: High engagement (lots of reactions and comments)
	p2, _ := feedSvc.CreatePost(ctx, "user-author", &domain.CreatePostRequest{
		Content: "Viral post with lots of engagement",
	})
	_ = feedSvc.ReactToPost(ctx, "user-a", p2.ID, domain.ReactionLove)
	_ = feedSvc.ReactToPost(ctx, "user-b", p2.ID, domain.ReactionHaha)
	_, _ = feedSvc.AddComment(ctx, "user-c", p2.ID, &domain.AddCommentRequest{Content: "Nice!"})

	// Ranked feed should place post 2 ahead of post 1
	rankedFeed, err := feedSvc.GetFeed(ctx, "user-author", "", 10, "ranked")
	if err != nil || len(rankedFeed.Items) < 2 {
		t.Fatalf("expected at least 2 items, got: %v", len(rankedFeed.Items))
	}

	if rankedFeed.Items[0].Post.ID != p2.ID {
		t.Errorf("expected high engagement post %s first in ranked feed, got %s", p2.ID, rankedFeed.Items[0].Post.ID)
	}
	if rankedFeed.Items[1].Post.ID != p1.ID {
		t.Errorf("expected low engagement post %s second, got %s", p1.ID, rankedFeed.Items[1].Post.ID)
	}
	if rankedFeed.Items[0].RankingScore <= rankedFeed.Items[1].RankingScore {
		t.Errorf("score for p2 (%f) should exceed p1 (%f)", rankedFeed.Items[0].RankingScore, rankedFeed.Items[1].RankingScore)
	}
}

func TestSavedCollections(t *testing.T) {
	feedSvc, _, _ := setupTestFeedService()
	ctx := context.Background()

	post, _ := feedSvc.CreatePost(ctx, "author-1", &domain.CreatePostRequest{
		Content: "Recipe for chocolate cake",
	})

	err := feedSvc.SavePost(ctx, "user-collector", "Recipes", post.ID)
	if err != nil {
		t.Fatalf("failed to save post: %v", err)
	}

	collections, err := feedSvc.GetSavedCollections(ctx, "user-collector")
	if err != nil || len(collections) != 1 {
		t.Fatalf("expected 1 collection, got %d", len(collections))
	}

	if collections[0].Name != "Recipes" || len(collections[0].PostIDs) != 1 {
		t.Errorf("collection data mismatch")
	}
	if collections[0].PostIDs[0] != post.ID {
		t.Errorf("expected saved post ID %s, got: %s", post.ID, collections[0].PostIDs[0])
	}
}

func TestMemories_OnThisDay(t *testing.T) {
	feedSvc, repo, _ := setupTestFeedService()
	ctx := context.Background()

	now := time.Now().UTC()
	oneYearAgo := now.AddDate(-1, 0, 0)

	// Inject a past memory post directly into repo
	pastPost := &domain.Post{
		ID:             "past-post-1",
		AuthorID:       "user-memories",
		Content:        "Remembering our trip 1 year ago!",
		Privacy:        domain.PrivacyFriends,
		ReactionCounts: make(map[domain.ReactionType]int),
		CreatedAt:      oneYearAgo,
		UpdatedAt:      oneYearAgo,
	}
	_ = repo.CreatePost(ctx, pastPost)

	// Modern post today (should not be in memories)
	_, _ = feedSvc.CreatePost(ctx, "user-memories", &domain.CreatePostRequest{
		Content: "Today's current post",
	})

	memories, err := feedSvc.GetMemories(ctx, "user-memories")
	if err != nil {
		t.Fatalf("failed to get memories: %v", err)
	}

	if len(memories) != 1 {
		t.Fatalf("expected exactly 1 memory from past years, got: %d", len(memories))
	}
	if memories[0].ID != "past-post-1" {
		t.Errorf("expected past-post-1, got: %s", memories[0].ID)
	}
}
