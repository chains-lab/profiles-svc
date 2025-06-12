package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const profilesCollection = "profiles"

type ProfileModel struct {
	ID          uuid.UUID `db:"id"`
	Username    string    `db:"username"`
	Pseudonym   *string   `db:"pseudonym,omitempty"`
	Description *string   `db:"description,omitempty"`
	AvatarURL   *string   `db:"avatar_url,omitempty"`
	Official    bool      `db:"official"`
	UpdatedAt   time.Time `db:"updated_at"`
	CreatedAt   time.Time `db:"created_at"`
}

type Profiles struct {
	client   *redis.Client
	lifeTime time.Duration
	log      *logrus.Entry
}

func NewProfiles(lifeTime time.Duration, client *redis.Client, log *logrus.Logger) Profiles {
	return Profiles{
		client:   client,
		lifeTime: lifeTime,
		log:      log.WithField("component", "redis"),
	}
}

func idKey(id uuid.UUID) string {
	return profilesCollection + ":id:" + id.String()
}

func usernameKey(username string) string {
	return profilesCollection + ":username:" + username
}

func (p *Profiles) Set(ctx context.Context, profile ProfileModel) {
	keyByID := idKey(profile.ID)
	keyByUser := usernameKey(profile.Username)

	data := map[string]interface{}{
		"username":    profile.Username,
		"pseudonym":   profile.Pseudonym,
		"description": profile.Description,
		"avatar_url":  profile.AvatarURL,
		"official":    profile.Official,
		"updated_at":  profile.UpdatedAt.Format(time.RFC3339),
		"created_at":  profile.CreatedAt.Format(time.RFC3339),
	}

	// Выполняем всё в одном pipeline: запись + TTL
	pipe := p.client.Pipeline()
	pipe.HSet(ctx, keyByID, data)
	pipe.Expire(ctx, keyByID, p.lifeTime)
	pipe.HSet(ctx, keyByUser, "id", profile.ID.String())
	pipe.Expire(ctx, keyByUser, p.lifeTime)
	if _, err := pipe.Exec(ctx); err != nil {
		p.log.WithError(err).Error("redis pipeline exec failed")
	}
}

func (p *Profiles) GetByID(ctx context.Context, id uuid.UUID) (ProfileModel, error) {
	keyByID := idKey(id)
	data, err := p.client.HGetAll(ctx, keyByID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			p.log.WithError(err).Errorf("get profile by id %s from redis failed", id)
		}

		return ProfileModel{}, redis.Nil
	}
	if len(data) == 0 {
		return ProfileModel{}, redis.Nil
	}

	res, err := parseProfile(id.String(), data)
	if err != nil {
		p.log.WithError(err).Errorf("parse profile from redis failed: %v", err)
		return ProfileModel{}, redis.Nil
	}

	return res, nil
}

func (p *Profiles) GetByUsername(ctx context.Context, username string) (ProfileModel, error) {
	keyByUser := usernameKey(username)
	ID, err := p.client.HGet(ctx, keyByUser, "id").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			p.log.WithError(err).Errorf("get profile by username %s from redis failed", username)
		}

		return ProfileModel{}, redis.Nil
	}

	data, err := p.client.HGetAll(ctx, ID).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			p.log.WithError(err).Errorf("get profile data by ID %s from redis failed", ID)
		}

		return ProfileModel{}, redis.Nil
	}
	if len(data) == 0 {
		return ProfileModel{}, redis.Nil
	}

	res, err := parseProfile(ID, data)
	if err != nil {
		p.log.WithError(err).Errorf("parse profile from redis failed: %v", err)
		return ProfileModel{}, redis.Nil
	}

	return res, nil
}

func parseProfile(id string, data map[string]string) (ProfileModel, error) {
	profile := ProfileModel{
		ID:       uuid.MustParse(id),
		Username: data["username"],
		Official: data["official"] == "true",
	}

	var err error
	profile.UpdatedAt, err = time.Parse(time.RFC3339, data["updated_at"])
	if err != nil {
		return ProfileModel{}, fmt.Errorf("parse updated_at: %w", err)
	}
	profile.CreatedAt, err = time.Parse(time.RFC3339, data["created_at"])
	if err != nil {
		return ProfileModel{}, fmt.Errorf("parse created_at: %w", err)
	}
	if data["pseudonym"] != "" {
		pseudonym := data["pseudonym"]
		profile.Pseudonym = &pseudonym
	}
	if data["description"] != "" {
		description := data["description"]
		profile.Description = &description
	}
	if data["avatar_url"] != "" {
		avatarURL := data["avatar_url"]
		profile.AvatarURL = &avatarURL
	}

	return profile, nil
}
