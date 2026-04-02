package api

import (
	"os"
	"strings"
	"testing"
	"time"
)

// dateVideo is a known date where APOD is a video (no hdurl).
var dateVideo = time.Date(2019, 12, 23, 0, 0, 0, 0, time.Local)

// dateImage is a known date where APOD is an image.
var dateImage = time.Date(2019, 12, 24, 0, 0, 0, 0, time.Local)

func TestGetInfoSpecificDate(t *testing.T) {
	client := NewClient("")
	info, err := client.GetInfo(dateVideo)
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}
	if info.MediaType != "video" {
		t.Errorf("expected media_type=video, got %s", info.MediaType)
	}
	if info.HDURL != "" {
		t.Errorf("expected empty hdurl for video, got %s", info.HDURL)
	}
}

func TestGetInfoImageDate(t *testing.T) {
	client := NewClient("")
	info, err := client.GetInfo(dateImage)
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}
	if info.MediaType != "image" {
		t.Errorf("expected media_type=image, got %s", info.MediaType)
	}
	if info.HDURL == "" {
		t.Error("expected non-empty hdurl for image")
	}
}

func TestDownloadImagePath(t *testing.T) {
	client := NewClient("")
	path, err := client.DownloadImage(dateImage)
	if err != nil {
		t.Fatalf("DownloadImage failed: %v", err)
	}
	if !strings.HasSuffix(path, ".jpg") {
		t.Errorf("expected .jpg suffix, got %s", path)
	}
	if !strings.Contains(path, "2019-12-24") {
		t.Errorf("expected path to contain date, got %s", path)
	}

	// Clean up
	os.Remove(path)
}

func TestDownloadImageVideoDate(t *testing.T) {
	client := NewClient("")
	_, err := client.DownloadImage(dateVideo)
	if err == nil {
		t.Error("expected error for video date, got nil")
	}
	if !strings.Contains(err.Error(), "image not found") {
		t.Errorf("expected 'image not found' error, got: %v", err)
	}
}

func TestDownloadImageCached(t *testing.T) {
	client := NewClient("")

	// First download
	path, err := client.DownloadImage(dateImage)
	if err != nil {
		t.Fatalf("first download failed: %v", err)
	}

	// Second download should use cache
	path2, err := client.DownloadImage(dateImage)
	if err != nil {
		t.Fatalf("cached download failed: %v", err)
	}
	if path != path2 {
		t.Errorf("cached path differs: %s vs %s", path, path2)
	}

	// Clean up
	os.Remove(path)
}
