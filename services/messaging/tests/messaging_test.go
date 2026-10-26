package tests

import (
	"context"
	"testing"

	"github.com/dnjstjqppk-dev/facebook/pkg/domain"
	"github.com/dnjstjqppk-dev/facebook/services/messaging/internal/repository"
	"github.com/dnjstjqppk-dev/facebook/services/messaging/internal/service"
)

func setupTestMessagingService() (*service.MessagingService, repository.MessagingRepository) {
	repo := repository.NewInMemoryMessagingRepository()
	svc := service.NewMessagingService(repo)
	return svc, repo
}

func TestCreate1on1_AndMessaging(t *testing.T) {
	svc, repo := setupTestMessagingService()
	ctx := context.Background()

	// 1. Create 1:1 conversation
	conv, err := svc.Create1on1Conversation(ctx, "alice", "bob", false)
	if err != nil {
		t.Fatalf("failed to create 1:1 conversation: %v", err)
	}

	if conv.Type != domain.ConvOneToOne || len(conv.Members) != 2 {
		t.Errorf("conversation structure mismatch: %v", conv)
	}

	// 2. Alice sends a message
	msg, err := svc.SendMessage(ctx, "alice", conv.ID, "Hey Bob, how are you?", "", "")
	if err != nil {
		t.Fatalf("failed to send message: %v", err)
	}

	if msg.Content != "Hey Bob, how are you?" || msg.SenderID != "alice" {
		t.Errorf("message fields mismatch: %v", msg)
	}

	// Verify conversation snippet updated
	c, _ := repo.GetConversation(ctx, conv.ID)
	if c.LastMessageSnippet != "Hey Bob, how are you?" {
		t.Errorf("expected last snippet updated, got: %s", c.LastMessageSnippet)
	}
}

func TestGroupChat_WithAdmin(t *testing.T) {
	svc, _ := setupTestMessagingService()
	ctx := context.Background()

	conv, err := svc.CreateGroupConversation(ctx, "admin-user", "Weekend Football", []string{"member-1", "member-2"})
	if err != nil {
		t.Fatalf("failed to create group chat: %v", err)
	}

	if conv.Type != domain.ConvGroup || conv.Title != "Weekend Football" {
		t.Errorf("group metadata mismatch")
	}
	if conv.Members["admin-user"] != domain.MemberRoleAdmin {
		t.Errorf("expected creator to be ADMIN, got %s", conv.Members["admin-user"])
	}
	if conv.Members["member-1"] != domain.MemberRoleMember {
		t.Errorf("expected member-1 to be MEMBER")
	}
}

func TestDeliveryAndReadReceipts(t *testing.T) {
	svc, repo := setupTestMessagingService()
	ctx := context.Background()

	conv, _ := svc.Create1on1Conversation(ctx, "alice", "bob", false)
	msg, _ := svc.SendMessage(ctx, "alice", conv.ID, "Meeting at 3PM", "", "")

	// 1. Bob receives message
	err := svc.MarkDelivered(ctx, msg.ID, "bob")
	if err != nil {
		t.Fatalf("failed to mark delivered: %v", err)
	}

	m, _ := repo.GetMessage(ctx, msg.ID)
	if _, delivered := m.DeliveredTo["bob"]; !delivered {
		t.Errorf("expected Bob in DeliveredTo")
	}

	// 2. Bob reads message
	err = svc.MarkRead(ctx, conv.ID, msg.ID, "bob")
	if err != nil {
		t.Fatalf("failed to mark read: %v", err)
	}

	m, _ = repo.GetMessage(ctx, msg.ID)
	if _, read := m.ReadBy["bob"]; !read {
		t.Errorf("expected Bob in ReadBy")
	}
}

func TestUnsendMessage(t *testing.T) {
	svc, repo := setupTestMessagingService()
	ctx := context.Background()

	conv, _ := svc.Create1on1Conversation(ctx, "alice", "bob", false)
	msg, _ := svc.SendMessage(ctx, "alice", conv.ID, "Oops, wrong chat!", "", "")

	// Unauthorized unsend attempt by Bob
	err := svc.UnsendMessage(ctx, msg.ID, "bob")
	if err != service.ErrCannotUnsend {
		t.Errorf("expected ErrCannotUnsend for non-author, got: %v", err)
	}

	// Authorized unsend by Alice
	err = svc.UnsendMessage(ctx, msg.ID, "alice")
	if err != nil {
		t.Fatalf("unsend by author failed: %v", err)
	}

	m, _ := repo.GetMessage(ctx, msg.ID)
	if !m.IsUnsent || m.Content != "[This message was unsent]" {
		t.Errorf("message content not marked as unsent")
	}
}

func TestMessageReactions(t *testing.T) {
	svc, repo := setupTestMessagingService()
	ctx := context.Background()

	conv, _ := svc.Create1on1Conversation(ctx, "alice", "bob", false)
	msg, _ := svc.SendMessage(ctx, "alice", conv.ID, "Check this out!", "", "")

	err := svc.ReactToMessage(ctx, msg.ID, "bob", domain.ReactionLove)
	if err != nil {
		t.Fatalf("failed to react to message: %v", err)
	}

	m, _ := repo.GetMessage(ctx, msg.ID)
	if m.Reactions["bob"] != domain.ReactionLove {
		t.Errorf("expected Bob's reaction LOVE, got: %s", m.Reactions["bob"])
	}
}

func TestUserPresence_OnlineOffline(t *testing.T) {
	svc, _ := setupTestMessagingService()
	ctx := context.Background()

	_ = svc.SetPresence(ctx, "user-online", domain.PresenceOnline)
	pres, err := svc.GetPresence(ctx, "user-online")
	if err != nil || pres.Status != domain.PresenceOnline {
		t.Errorf("expected presence ONLINE, got: %v", pres)
	}

	_ = svc.SetPresence(ctx, "user-online", domain.PresenceOffline)
	pres, _ = svc.GetPresence(ctx, "user-online")
	if pres.Status != domain.PresenceOffline {
		t.Errorf("expected presence OFFLINE, got: %s", pres.Status)
	}
}

func TestBroadcastChannel_OneToMany(t *testing.T) {
	svc, _ := setupTestMessagingService()
	ctx := context.Background()

	channel, err := svc.CreateBroadcastChannel(ctx, "creator-pro", "Official Updates")
	if err != nil {
		t.Fatalf("failed to create broadcast channel: %v", err)
	}

	if channel.Type != domain.ConvBroadcastChannel || channel.Members["creator-pro"] != domain.MemberRoleAdmin {
		t.Errorf("broadcast channel metadata mismatch")
	}
}

func TestWebRTCCallLifecycle_InitiateAnswerEnd(t *testing.T) {
	svc, _ := setupTestMessagingService()
	ctx := context.Background()

	// 1. Caller initiates video call
	call, err := svc.InitiateCall(ctx, "caller-1", []string{"receiver-2"}, domain.CallVideo)
	if err != nil {
		t.Fatalf("failed to initiate call: %v", err)
	}
	if call.Status != domain.CallInitiated || call.CallType != domain.CallVideo {
		t.Errorf("unexpected initial call state")
	}

	// 2. Receiver answers call
	answered, err := svc.AnswerCall(ctx, call.ID, "receiver-2")
	if err != nil || answered.Status != domain.CallConnected || answered.StartedAt == nil {
		t.Fatalf("failed to answer call: %v", err)
	}

	// 3. Either party ends call
	ended, err := svc.EndCall(ctx, call.ID, "caller-1")
	if err != nil || ended.Status != domain.CallEnded || ended.EndedAt == nil {
		t.Fatalf("failed to end call: %v", err)
	}
}

func TestE2EEConversationFlag(t *testing.T) {
	svc, _ := setupTestMessagingService()
	ctx := context.Background()

	conv, err := svc.Create1on1Conversation(ctx, "agent-007", "agent-m", true)
	if err != nil {
		t.Fatalf("failed to create E2EE chat: %v", err)
	}
	if !conv.IsE2EE {
		t.Errorf("expected conversation IsE2EE = true")
	}

	msg, _ := svc.SendMessage(ctx, "agent-007", conv.ID, "Encrypted payload ciphertext", "", "")
	if !msg.IsE2EE {
		t.Errorf("expected message IsE2EE = true")
	}
}

func TestNonMemberCannotSendMessage(t *testing.T) {
	svc, _ := setupTestMessagingService()
	ctx := context.Background()

	conv, _ := svc.Create1on1Conversation(ctx, "alice", "bob", false)

	// Mallory tries to inject a message into Alice and Bob's private conversation
	_, err := svc.SendMessage(ctx, "mallory", conv.ID, "Malicious eavesdrop", "", "")
	if err != service.ErrNotConversationMember {
		t.Errorf("expected ErrNotConversationMember, got: %v", err)
	}
}
