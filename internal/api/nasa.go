package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mpuccini/astrowall/internal/config"
)

const (
	baseURL    = "https://api.nasa.gov/"
	endpoint   = "planetary/apod"
	defaultKey = "DEMO_KEY"
)

// APODResponse represents the JSON response from the NASA APOD API.
type APODResponse struct {
	Date           string `json:"date"`
	Title          string `json:"title"`
	Explanation    string `json:"explanation"`
	URL            string `json:"url"`
	HDURL          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	Copyright      string `json:"copyright"`
	ServiceVersion string `json:"service_version"`
}

// Client interacts with the NASA APOD API.
type Client struct {
	APIKey   string
	SaveDir  string
	HTTPClient *http.Client
}

// NewClient creates a new API client.
// Key precedence: explicit apiKey > NASA_API_KEY env var > config file > DEMO_KEY.
func NewClient(apiKey string) *Client {
	if apiKey == "" {
		apiKey = os.Getenv("NASA_API_KEY")
	}
	if apiKey == "" {
		apiKey = config.Load().APIKey
	}
	if apiKey == "" {
		apiKey = defaultKey
	}

	homeDir, _ := os.UserHomeDir()
	saveDir := filepath.Join(homeDir, "Pictures", "NASA")
	os.MkdirAll(saveDir, 0o755)

	return &Client{
		APIKey:     apiKey,
		SaveDir:    saveDir,
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
	}
}

// GetInfo fetches APOD metadata for the given date.
func (c *Client) GetInfo(date time.Time) (*APODResponse, error) {
	url := fmt.Sprintf("%s%s?date=%s&hd=true&api_key=%s",
		baseURL, endpoint, date.Format("2006-01-02"), c.APIKey)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not download meta-info for APOD %s: %w",
			date.Format("2006-01-02"), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d for date %s",
			resp.StatusCode, date.Format("2006-01-02"))
	}

	var apod APODResponse
	if err := json.NewDecoder(resp.Body).Decode(&apod); err != nil {
		return nil, fmt.Errorf("could not parse API response: %w", err)
	}

	return &apod, nil
}

// DownloadImage downloads the HD APOD image and returns the local file path.
// If the image was already downloaded, it returns the cached path without
// hitting the API again.
func (c *Client) DownloadImage(date time.Time) (string, error) {
	// Check cache first by scanning directory for files matching the date prefix.
	// This avoids an API call when the image is already downloaded.
	datePrefix := date.Format("2006-01-02") + "_"
	if entries, err := os.ReadDir(c.SaveDir); err == nil {
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), datePrefix) && strings.HasSuffix(e.Name(), ".jpg") {
				filePath := filepath.Join(c.SaveDir, e.Name())
				fmt.Println("Today's image has already been downloaded and is now being set as background.")
				return filePath, nil
			}
		}
	}

	info, err := c.GetInfo(date)
	if err != nil {
		return "", err
	}

	if info.HDURL == "" {
		return "", fmt.Errorf("image not found for the selected date %s (media type: %s)",
			date.Format("2006-01-02"), info.MediaType)
	}

	fmt.Printf("Title: %s\n", info.Title)

	safeName := strings.ReplaceAll(info.Title, " ", "-")
	fileName := fmt.Sprintf("%s_%s.jpg", date.Format("2006-01-02"), safeName)
	filePath := filepath.Join(c.SaveDir, fileName)

	fmt.Printf("Downloading - %s (%s)\n", info.Title, date.Format("2006-01-02"))

	resp, err := c.HTTPClient.Get(info.HDURL)
	if err != nil {
		return "", fmt.Errorf("could not download: %s: %w", info.Title, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("could not download: %s (HTTP %d)", info.Title, resp.StatusCode)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("could not create file %s: %w", filePath, err)
	}
	defer out.Close()

	totalSize := resp.ContentLength
	written := int64(0)
	buf := make([]byte, 4096)

	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := out.Write(buf[:n]); writeErr != nil {
				return "", fmt.Errorf("could not write to file: %w", writeErr)
			}
			written += int64(n)
			if totalSize > 0 {
				pct := float64(written) / float64(totalSize) * 100
				fmt.Printf("\r  [%-50s] %6.2f%% (%d / %d bytes)",
					strings.Repeat("=", int(pct/2))+">",
					pct, written, totalSize)
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return "", fmt.Errorf("error during download: %w", readErr)
		}
	}
	fmt.Println()

	return filePath, nil
}
