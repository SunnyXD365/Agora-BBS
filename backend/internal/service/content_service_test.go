package service

import "testing"

func TestEscapeLike(t *testing.T) {
	if got := escapeLike(`50%_done\next`); got != `50\%\_done\\next` {
		t.Fatalf("escapeLike() = %q", got)
	}
}

func TestNormalizePublicPage(t *testing.T) {
	page, pageSize := normalizePublicPage(-1, 25)
	if page != 1 || pageSize != 10 {
		t.Fatalf("normalizePublicPage() = %d/%d", page, pageSize)
	}
	_, pageSize = normalizePublicPage(2, 50)
	if pageSize != 50 {
		t.Fatalf("expected page size 50, got %d", pageSize)
	}
}
