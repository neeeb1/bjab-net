package meta

import "testing"

type item struct {
	id   string
	date string
}

func (i item) GetDate() string { return i.date }

func TestSortByDate_DescendingOrder(t *testing.T) {
	in := []item{
		{"a", "2024-01-01"},
		{"b", "2026-05-13"},
		{"c", "2025-06-15"},
	}
	got, err := SortByDate(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantIDs := []string{"b", "c", "a"}
	for i, w := range wantIDs {
		if got[i].id != w {
			t.Errorf("idx %d: got %q want %q", i, got[i].id, w)
		}
	}
}

func TestSortByDate_DoesNotMutateInput(t *testing.T) {
	in := []item{{"a", "2024-01-01"}, {"b", "2026-01-01"}}
	_, _ = SortByDate(in)
	if in[0].id != "a" || in[1].id != "b" {
		t.Errorf("input mutated: %+v", in)
	}
}

func TestSortByDate_EmptySlice(t *testing.T) {
	got, err := SortByDate([]item{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("want empty, got %d items", len(got))
	}
}

func TestSortByDate_MalformedDateReturnsError(t *testing.T) {
	in := []item{{"a", "not-a-date"}, {"b", "2025-01-01"}}
	_, err := SortByDate(in)
	if err == nil {
		t.Error("expected error for malformed date, got nil")
	}
}

func TestSortByDate_EqualDatesStable(t *testing.T) {
	in := []item{
		{"a", "2025-01-01"},
		{"b", "2025-01-01"},
		{"c", "2025-01-01"},
	}
	got, _ := SortByDate(in)
	for i, w := range []string{"a", "b", "c"} {
		if got[i].id != w {
			t.Errorf("idx %d: got %q want %q (sort not stable)", i, got[i].id, w)
		}
	}
}
