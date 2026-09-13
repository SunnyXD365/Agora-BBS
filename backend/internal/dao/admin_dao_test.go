package dao

import "testing"

func TestAdminSortColumnUsesWhitelist(t *testing.T) {
	allowed := map[string]string{"created_at": "u.created_at", "level": "tp.unlock_level"}
	if got := adminSortColumn("level", allowed, "u.created_at"); got != "tp.unlock_level" {
		t.Fatalf("allowed sort = %q", got)
	}
	if got := adminSortColumn("created_at; DROP TABLE users", allowed, "u.created_at"); got != "u.created_at" {
		t.Fatalf("untrusted sort did not fall back: %q", got)
	}
}

func TestAdminSortDirectionUsesWhitelist(t *testing.T) {
	if got := adminSortDirection("asc"); got != "ASC" {
		t.Fatalf("ascending direction = %q", got)
	}
	if got := adminSortDirection("desc; DROP TABLE users"); got != "DESC" {
		t.Fatalf("untrusted direction did not fall back: %q", got)
	}
}
