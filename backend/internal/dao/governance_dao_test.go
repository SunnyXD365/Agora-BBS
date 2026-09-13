package dao

import "testing"

func TestReadingEligibleForOrdinaryTopic(t *testing.T) {
	if !readingEligible(true, false, 0, 10) {
		t.Fatal("ordinary topic should become eligible after reaching the bottom")
	}
}

func TestReadingEligibleForLongOrSensitiveTopic(t *testing.T) {
	if readingEligible(true, true, 9, 10) {
		t.Fatal("topic requiring dwell should remain ineligible before the dwell threshold")
	}
	if !readingEligible(true, true, 10, 10) {
		t.Fatal("topic requiring dwell should become eligible at the dwell threshold")
	}
	if readingEligible(false, true, 10, 10) {
		t.Fatal("reply dwell alone must not make a session eligible")
	}
}
