package persistence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/khaledhikmat/threat-detection-shared/models"
)

// type apiResponse struct {
// 	Index   string `json:"_index"`
// 	ID     string `json:"_id"`
// 	Version string `json:"_version"`
// 	Result  string `json:"result"`
// }

type basicAuthRoundTripper struct {
	Next     http.RoundTripper
	Username string
	Password string
}

func (b *basicAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(b.Username, b.Password)
	req.Header.Set("Content-Type", "application/json")
	return b.Next.RoundTrip(req)
}

type loggerRoundTripper struct {
	Next   http.RoundTripper
	Logger io.Writer
}

func (l *loggerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Fprintf(l.Logger, "[%s] %s %s %v\n", time.Now(), req.Method, req.URL, req.Header)
	resp, err := l.Next.RoundTrip(req)
	if err != nil {
		fmt.Fprintf(l.Logger, "😢 error: %v\n", err)
		return resp, err
	}

	fmt.Fprintf(l.Logger, "😀 response: %v\n", resp)
	return resp, err
}

func newOpenSearch() IService {
	if os.Getenv("OPEN_SEARCH_DOMAIN_ENDPOINT") == "" ||
		os.Getenv("OPEN_SEARCH_INDEX_NAME") == "" ||
		os.Getenv("OPEN_SEARCH_USERNAME") == "" ||
		os.Getenv("OPEN_SEARCH_PASSWORD") == "" {
		panic(fmt.Errorf("😢 missing required environment variables"))
	}

	return &opensearch{}
}

type opensearch struct {
}

func (p *opensearch) NewClip(clip models.RecordingClip) error {
	fmt.Printf("***** 💪 inserting a new clip with id %s\n", clip.ID)
	client := &http.Client{
		Transport: &loggerRoundTripper{
			Next: &basicAuthRoundTripper{
				Next:     http.DefaultTransport,
				Username: os.Getenv("OPEN_SEARCH_USERNAME"),
				Password: os.Getenv("OPEN_SEARCH_PASSWORD"),
			},
			Logger: os.Stdout,
		},
	}

	url := fmt.Sprintf("%s/%s/_doc/%s", os.Getenv("OPEN_SEARCH_DOMAIN_ENDPOINT"), os.Getenv("OPEN_SEARCH_INDEX_NAME"), clip.ID)
	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(clip)
	if err != nil {
		fmt.Printf("😢 error encoding clip %v\n", err)
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, payloadBuf)
	if err != nil {
		fmt.Printf("😢 error creating a new request %v\n", err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("😢 error inserting a new clip %v\n", err)
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("😢 error inserting a new clip %v\n", err)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	fmt.Printf("***** 😀 inserted a clip successfully 😎\n")
	return nil
}

func (p *opensearch) RetrieveClipCount(lastDays int) (int, error) {
	return 0, nil
}

func (p *opensearch) RetrieveClipsStatsByRegion(lastDays int) ([]models.ClipStats, error) {
	return []models.ClipStats{}, nil
}

func (p *opensearch) RetrieveClipsByRegion(region string, page, pageSize int) ([]models.RecordingClip, error) {
	return []models.RecordingClip{}, nil
}

func (p *opensearch) RetrieveTopCapturers(top int, lastDays int) ([]string, error) {
	return []string{}, nil
}

func (p *opensearch) RetrieveTopCameras(top int, lastDays int) ([]string, error) {
	return []string{}, nil
}

func (p *opensearch) RetrieveTopRegions(top int, lastDays int) ([]string, error) {
	return []string{}, nil
}

func (p *opensearch) RetrieveTopLocations(top int, lastDays int) ([]string, error) {
	return []string{}, nil
}

func (p *opensearch) RetrieveTopTags(top int, lastDays int) ([]string, error) {
	return []string{}, nil
}

func (p *opensearch) RetrieveTopAnalytics(top int, lastDays int) ([]string, error) {
	return []string{}, nil
}

func (p *opensearch) RetrieveTopAlertTypes(top int, lastDays int) ([]string, error) {
	return []string{}, nil
}

func (p *opensearch) RetrieveClipsByTag(tag string, lastDays int) ([]models.RecordingClip, error) {
	return []models.RecordingClip{}, nil
}

func (p *opensearch) RetrieveClipsByTags(tags []string, lastDays int) ([]models.RecordingClip, error) {
	return []models.RecordingClip{}, nil
}

func (p *opensearch) RetrieveClipsByAnalytic(analytic string, lastDays int) ([]models.RecordingClip, error) {
	return []models.RecordingClip{}, nil
}

func (p *opensearch) RetrieveClipsByAnalytics(analytics []string, lastDays int) ([]models.RecordingClip, error) {
	return []models.RecordingClip{}, nil
}

func (p *opensearch) RetrieveClipsByAlertType(typ string, lastDays int) ([]models.RecordingClip, error) {
	return []models.RecordingClip{}, nil
}

func (p *opensearch) RetrieveClipsByAlertTypes(typs []string, lastDays int) ([]models.RecordingClip, error) {
	return []models.RecordingClip{}, nil
}

func (p *opensearch) Finalize() {
}
