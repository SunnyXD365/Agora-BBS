package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"agora-backend/internal/config"
	"agora-backend/internal/dao"
	"agora-backend/internal/model"
	"agora-backend/internal/workflow"
)

var (
	ErrReplyLocked            = errors.New("reply permission is not unlocked")
	ErrTopicLocked            = errors.New("topic creation permission is not unlocked")
	ErrReadingRequired        = errors.New("complete the required reading before replying")
	ErrInvalidReadingResource = errors.New("invalid reading resource")
)

type GovernanceService struct {
	dao     *dao.GovernanceDAO
	users   *dao.UserDAO
	cfg     *config.Config
	reviews workflow.BlindReviewStarter
}

func NewGovernanceService(governanceDAO *dao.GovernanceDAO, users *dao.UserDAO, cfg *config.Config, reviews workflow.BlindReviewStarter) *GovernanceService {
	return &GovernanceService{dao: governanceDAO, users: users, cfg: cfg, reviews: reviews}
}

func (s *GovernanceService) Policy() model.GovernancePolicy {
	l1, l2, l3 := 1800, 7200, 18000
	if s.cfg.AppEnv == "development" {
		l1, l2, l3 = 60, 180, 300
	}
	return model.GovernancePolicy{
		CoolingSeconds:    s.cfg.GovernanceCoolingSeconds,
		ReplyDwellSeconds: s.cfg.GovernanceReplyDwellSeconds,
		HeartbeatSeconds:  5,
		LongTopicChars:    s.cfg.GovernanceLongTopicChars,
		Level1ReadSeconds: l1, Level2ReadSeconds: l2, Level3ReadSeconds: l3,
	}
}

func (s *GovernanceService) StartReading(ctx context.Context, userID int64, req *model.StartReadingReq) (*model.ReadingSession, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}
	publicID := hex.EncodeToString(buf)
	var result *model.ReadingSession
	var err error
	if req.TopicID != nil && *req.TopicID > 0 && req.Resource == "" {
		result, err = s.dao.StartTopicReading(ctx, publicID, userID, *req.TopicID)
	} else if req.TopicID == nil && req.Resource == "forum-guide" {
		result, err = s.dao.StartResourceReading(ctx, publicID, userID, "guide", req.Resource)
	} else {
		return nil, ErrInvalidReadingResource
	}
	if err != nil {
		return nil, err
	}
	return s.decorateReadingRequirement(ctx, result)
}

func (s *GovernanceService) Heartbeat(ctx context.Context, userID int64, publicID string, req *model.ReadingHeartbeatReq) (*model.ReadingSession, error) {
	result, err := s.dao.Heartbeat(ctx, publicID, userID, req.Progress, req.ReplyFocused)
	if err != nil {
		return nil, err
	}
	return s.decorateReadingRequirement(ctx, result)
}

func (s *GovernanceService) CompleteReading(ctx context.Context, userID int64, publicID string, req *model.CompleteReadingReq) (*model.ReadingSession, error) {
	result, _, err := s.dao.CompleteReading(ctx, publicID, userID, s.cfg.GovernanceReplyDwellSeconds, s.cfg.GovernanceLongTopicChars, req)
	if err != nil {
		return nil, err
	}
	policy := s.Policy()
	if err := s.dao.UpdateUnlockLevel(ctx, userID, s.cfg.AppEnv == "development", policy.Level1ReadSeconds, policy.Level2ReadSeconds, policy.Level3ReadSeconds); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *GovernanceService) decorateReadingRequirement(ctx context.Context, result *model.ReadingSession) (*model.ReadingSession, error) {
	if result.ResourceType != "topic" || result.TopicID == nil {
		result.RequiresReplyDwell = false
		return result, nil
	}
	required, err := s.dao.TopicRequiresReading(ctx, *result.TopicID, s.cfg.GovernanceLongTopicChars)
	if err != nil {
		return nil, err
	}
	result.RequiresReplyDwell = required
	return result, nil
}

func (s *GovernanceService) SaveOnboarding(ctx context.Context, userID int64, req *model.OnboardingReq) error {
	if err := s.dao.SaveOnboarding(ctx, userID, req.Statement, req.BackgroundTag); err != nil {
		return err
	}
	return s.reviews.StartBlindReview(ctx, "user", userID)
}

func (s *GovernanceService) RequireReviewPermission(ctx context.Context, userID int64) error {
	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil || (user.UnlockLevel < 3 && user.Role != "admin") {
		return errors.New("blind review permission is not unlocked")
	}
	return nil
}

func (s *GovernanceService) RequireTopicPermission(ctx context.Context, userID int64) error {
	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil || (user.UnlockLevel < 2 && user.Role != "admin") {
		return ErrTopicLocked
	}
	return nil
}

func (s *GovernanceService) RequireReplyPermission(ctx context.Context, userID, topicID int64) error {
	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil || (user.UnlockLevel < 1 && user.Role != "admin") {
		return ErrReplyLocked
	}
	required, err := s.dao.TopicRequiresReading(ctx, topicID, s.cfg.GovernanceLongTopicChars)
	if err != nil {
		return err
	}
	if required {
		eligible, err := s.dao.HasEligibleReading(ctx, userID, topicID)
		if err != nil {
			return err
		}
		if !eligible {
			return ErrReadingRequired
		}
	}
	return nil
}

func (s *GovernanceService) RequireFeedbackPermission(ctx context.Context, userID int64) error {
	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil || (user.UnlockLevel < 1 && user.Role != "admin") {
		return ErrReplyLocked
	}
	return nil
}
