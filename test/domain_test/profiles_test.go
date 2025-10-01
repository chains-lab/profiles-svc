package domain_test

import (
	"context"
	"testing"

	"github.com/chains-lab/profiles-svc/internal/domain/services/profile"
	"github.com/google/uuid"
)

func TestProfiles(t *testing.T) {
	s, err := newSetup(t)
	if err != nil {
		t.Fatalf("newSetup: %v", err)
	}

	cleanDb(t)

	ctx := context.Background()

	firstID := uuid.New()
	secondID := uuid.New()

	first, err := s.domain.profile.Create(ctx, firstID, "first")
	if err != nil {
		t.Fatalf("CreateProfile first: %v", err)
	}

	second, err := s.domain.profile.Create(ctx, secondID, "second")
	if err != nil {
		t.Fatalf("CreateProfile second: %v", err)
	}

	if first.UserID == second.UserID {
		t.Fatalf("expected different IDs, got same: %v", first.UserID)
	}

	first, err = s.domain.profile.GetByID(ctx, firstID)
	if err != nil {
		t.Fatalf("GetProfileByUserID first: %v", err)
	}

	if first.UserID != firstID {
		t.Fatalf("GetProfileByUserID first: expected ID %v, got %v", firstID, first.UserID)
	}

	second, err = s.domain.profile.GetByID(ctx, secondID)
	if err != nil {
		t.Fatalf("GetProfileByUserID second: %v", err)
	}

	if second.UserID != secondID {
		t.Fatalf("GetProfileByUserID second: expected ID %v, got %v", secondID, second.UserID)
	}

	avatar := "avatar"
	newFirst := "new_first"
	description := "description"

	first, err = s.domain.profile.Update(ctx, firstID, profile.Update{
		Avatar:      &avatar,
		Pseudonym:   &newFirst,
		Description: &description,
	})
	if err != nil {
		t.Fatalf("UpdateProfile first: %v", err)
	}

	if *first.Avatar != avatar {
		t.Fatalf("UpdateProfile first: expected avatar %s, got %s", avatar, *first.Avatar)
	}
	if *first.Pseudonym != newFirst {
		t.Fatalf("UpdateProfile first: expected pseudonym %s, got %s", newFirst, *first.Pseudonym)
	}
	if *first.Description != description {
		t.Fatalf("UpdateProfile first: expected description %s, got %s", description, *first.Description)
	}

	second, err = s.domain.profile.ResetUserProfile(ctx, secondID)
	if err != nil {
		t.Fatalf("ResetUserProfile second: %v", err)
	}
	if second.Avatar != nil {
		t.Fatalf("ResetUserProfile second: expected avatar nil, got %v", *second.Avatar)
	}
	if second.Pseudonym != nil {
		t.Fatalf("ResetUserProfile second: expected pseudonym nil, got %v", *second.Pseudonym)
	}
	if second.Description != nil {
		t.Fatalf("ResetUserProfile second: expected description nil, got %v", *second.Description)
	}

	if second.Username != "second" {
		t.Fatalf("ResetUserProfile second: expected username %s, got %s", "second", second.Username)
	}

	first, err = s.domain.profile.ResetUsername(ctx, firstID)
	if err != nil {
		t.Fatalf("ResetUsername first: %v", err)
	}
	if first.Username == "first" {
		t.Fatalf("ResetUsername first: expected username not %s, got %s", "first", first.Username)
	}

	first, err = s.domain.profile.UpdateOfficial(ctx, firstID, false)
	if err != nil {
		t.Fatalf("UpdateProfileOfficial first to false: %v", err)
	}

	if first.Official {
		t.Fatalf("UpdateProfileOfficial first to false: expected official false, got true")
	}

	first, err = s.domain.profile.UpdateOfficial(ctx, firstID, true)
	if err != nil {
		t.Fatalf("UpdateProfileOfficial first to true: %v", err)
	}

	if !first.Official {
		t.Fatalf("UpdateProfileOfficial first to true: expected official true, got false")
	}
}
