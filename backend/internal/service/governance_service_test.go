package service

import (
	"testing"

	"agora-backend/internal/config"
)

func TestDevelopmentPolicyUsesShortThresholds(t *testing.T) {
	service := NewGovernanceService(nil, nil, &config.Config{
		AppEnv: "development", GovernanceCoolingSeconds: 60,
		GovernanceReplyDwellSeconds: 10, GovernanceLongTopicChars: 800,
	}, nil)
	policy := service.Policy()
	if policy.Level1ReadSeconds != 60 || policy.Level2ReadSeconds != 180 || policy.Level3ReadSeconds != 300 {
		t.Fatalf("development thresholds = %d/%d/%d", policy.Level1ReadSeconds, policy.Level2ReadSeconds, policy.Level3ReadSeconds)
	}
	if policy.HeartbeatSeconds != 5 || policy.CoolingSeconds != 60 {
		t.Fatalf("unexpected timing policy: %+v", policy)
	}
}

func TestProductionPolicyUsesFullThresholds(t *testing.T) {
	service := NewGovernanceService(nil, nil, &config.Config{AppEnv: "production"}, nil)
	policy := service.Policy()
	if policy.Level1ReadSeconds != 1800 || policy.Level2ReadSeconds != 7200 || policy.Level3ReadSeconds != 18000 {
		t.Fatalf("production thresholds = %d/%d/%d", policy.Level1ReadSeconds, policy.Level2ReadSeconds, policy.Level3ReadSeconds)
	}
}
