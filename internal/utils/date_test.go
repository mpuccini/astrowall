package utils

import (
	"testing"
	"time"
)

func TestParseDateFull(t *testing.T) {
	got, err := ParseDate("20190706")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2019, 7, 6, 0, 0, 0, 0, time.Local)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseDateDashed(t *testing.T) {
	got, err := ParseDate("2019-12-24")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2019, 12, 24, 0, 0, 0, 0, time.Local)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseDateMinified(t *testing.T) {
	got, err := ParseDate("2019-7-6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2019, 7, 6, 0, 0, 0, 0, time.Local)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseDateExtra(t *testing.T) {
	// Extra numeric groups beyond the first three should be ignored
	got, err := ParseDate("2019-7-6-99-88")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2019, 7, 6, 0, 0, 0, 0, time.Local)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseDateShort(t *testing.T) {
	_, err := ParseDate("190706")
	if err == nil {
		t.Error("expected error for 6-digit date, got nil")
	}
}

func TestParseDateIncomplete(t *testing.T) {
	_, err := ParseDate("2019-07")
	if err == nil {
		t.Error("expected error for incomplete date, got nil")
	}
}
