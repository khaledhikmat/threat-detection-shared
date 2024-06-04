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
	_, err := search[models.RecordingClip, map[string]interface{}](http.MethodPut, "_doc", clip.ID, clip)
	if err != nil {
		fmt.Printf("😢 error encoding clip %v\n", err)
		return err
	}
	fmt.Printf("***** 😀 inserted a clip successfully 😎\n")
	return nil
}

func (p *opensearch) RetrieveClipCount(lastPeriods int) (int, error) {
	fmt.Printf("***** 💪 Retrieve clips count %d\n", lastPeriods)
	if lastPeriods == -1 {
		lastPeriods = 1000
	}

	payload := map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"indexTime": map[string]interface{}{
					"gte": fmt.Sprintf("now-%dm/m", lastPeriods),
					"lte": "now/m",
				},
			},
		},
	}

	o, err := search[map[string]interface{}, map[string]interface{}](http.MethodGet, "_count", "", payload)
	if err != nil {
		fmt.Printf("😢 error retrieving clips count %v\n", err)
		return 0, err
	}

	return int(o["count"].(float64)), nil
}

func (p *opensearch) RetrieveClipsStatsByRegion(lastPeriods int) ([]models.ClipStats, error) {
	fmt.Printf("***** 💪 Retrieve clips stats by region %d\n", lastPeriods)
	stats := []models.ClipStats{}

	// TODO: I am sure this is not the right way to do it, but I am not sure how to do it correctly
	// Phase 1: Retrieve all region buckets in the lastPeriods
	payload1 := map[string]interface{}{
		"size": 0,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"indexTime": map[string]interface{}{
								"gte": fmt.Sprintf("now-%dm/m", lastPeriods),
								"lte": "now/m",
							},
						},
					},
				},
				"filter":   []map[string]interface{}{},
				"should":   []map[string]interface{}{},
				"must_not": []map[string]interface{}{},
			},
		},
		"aggs": map[string]interface{}{
			"region_buckets": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "region",
					"order": map[string]interface{}{
						"_count": "desc",
					},
					"size": 1000,
				},
			},
		},
	}

	o1, err := search[map[string]interface{}, map[string]interface{}](http.MethodGet, "_search", "", payload1)
	if err != nil {
		fmt.Printf("😢 error retrieving regions count %v\n", err)
		return stats, err
	}

	// WARNING: Need to make sure we check all data elements as this can produce a panic
	regions := []string{}
	aggregations, ok := o1["aggregations"].(map[string]interface{})
	if !ok {
		fmt.Printf("😢 expected aggregations element to be a map of interface{}....found nothing %v\n", err)
		return stats, nil
	}

	regionBuckets, ok := aggregations["region_buckets"].(map[string]interface{})
	if !ok {
		fmt.Printf("😢 expected region_buckets element  to be a map of interface{}....found nothing %v\n", err)
		return stats, nil
	}

	_, ok = regionBuckets["buckets"]
	if !ok {
		fmt.Printf("😢 expected buckets element....found nothing %v\n", err)
		return stats, nil
	}

	// This is really strange, why is assertion failing on []map[string]interface{} but not on []interface{}?
	buckets, ok := regionBuckets["buckets"].([]interface{})
	if !ok {
		fmt.Printf("😢 expected buckets element to be an array of map of interface{}....found nothing %v\n", err)
		return stats, nil
	}

	for _, bucket := range buckets {
		bucketValue, ok := bucket.(map[string]interface{})
		if ok {
			regions = append(regions, bucketValue["key"].(string))
		}
	}

	//regionBuckets := o1["aggregations"].(map[string]interface{})["region_buckets"].(map[string]interface{})["buckets"].([]interface{})
	// for _, bucket := range regionBuckets {
	// bucketValue, ok := bucket.(map[string]interface{})
	// if ok {
	// 		regions = append(regions, bucketValue["key"].(string))
	// }

	// Phase 2: Retrieve clips stats for each region in the lastPeriods
	for _, region := range regions {
		stat := models.ClipStats{
			Region: region,
		}

		payload2 := map[string]interface{}{
			"size": 0,
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{
							"range": map[string]interface{}{
								"indexTime": map[string]interface{}{
									"gte": fmt.Sprintf("now-%dm/m", lastPeriods),
									"lte": "now/m",
								},
							},
						},
						{
							"match": map[string]interface{}{
								"region": region,
							},
						},
					},
					"filter":   []map[string]interface{}{},
					"should":   []map[string]interface{}{},
					"must_not": []map[string]interface{}{},
				},
			},
			"aggs": map[string]interface{}{
				"cameras": map[string]interface{}{
					"cardinality": map[string]interface{}{
						"field": "cameraId",
					},
				},
				"clips": map[string]interface{}{
					"value_count": map[string]interface{}{
						"field": "_id",
					},
				},
				"frames": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "frames",
					},
				},
				"tags": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "tagsCount",
					},
				},
				"alerts": map[string]interface{}{
					"sum": map[string]interface{}{
						"field": "alertsCount",
					},
				},
			},
		}

		o2, err := search[map[string]interface{}, map[string]interface{}](http.MethodGet, "_search", "", payload2)
		if err != nil {
			fmt.Printf("😢 error retrieving aggregations count %v\n", err)
			return stats, err
		}

		stat.Cameras = int(o2["aggregations"].(map[string]interface{})["cameras"].(map[string]interface{})["value"].(float64))
		stat.Clips = int(o2["aggregations"].(map[string]interface{})["clips"].(map[string]interface{})["value"].(float64))
		stat.Frames = int(o2["aggregations"].(map[string]interface{})["frames"].(map[string]interface{})["value"].(float64))
		stat.Tags = int(o2["aggregations"].(map[string]interface{})["tags"].(map[string]interface{})["value"].(float64))
		stat.Alerts = int(o2["aggregations"].(map[string]interface{})["alerts"].(map[string]interface{})["value"].(float64))

		stats = append(stats, stat)
	}

	return stats, nil
}

func (p *opensearch) RetrieveAlertedClips(top, lastPeriods int) ([]models.RecordingClip, error) {
	payload := map[string]interface{}{
		"size": top,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"alertsCount": 1,
						},
					},
					{
						"range": map[string]interface{}{
							"indexTime": map[string]interface{}{
								"gte": fmt.Sprintf("now-%dm/m", lastPeriods),
								"lte": "now/m",
							},
						},
					},
				},
				"filter":   []map[string]interface{}{},
				"should":   []map[string]interface{}{},
				"must_not": []map[string]interface{}{},
			},
		},
		"sort": []map[string]interface{}{
			{
				"indexTime": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}

	return retrieveClips(payload)
}

func (p *opensearch) RetrieveClipsByRegion(region string, lastPeriods, page, pageSize int) ([]models.RecordingClip, error) {
	payload := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"region": region,
						},
					},
					{
						"range": map[string]interface{}{
							"indexTime": map[string]interface{}{
								"gte": fmt.Sprintf("now-%dm/m", lastPeriods),
								"lte": "now/m",
							},
						},
					},
				},
				"filter":   []map[string]interface{}{},
				"should":   []map[string]interface{}{},
				"must_not": []map[string]interface{}{},
			},
		},
		"sort": []map[string]interface{}{
			{
				"indexTime": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}

	return retrieveClips(payload)
}

func (p *opensearch) RetrieveClipByID(id string) (models.RecordingClip, error) {
	payload := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"id": id,
			},
		},
	}

	clips, err := retrieveClips(payload)
	if err != nil {
		return models.RecordingClip{}, err
	}

	if len(clips) == 0 {
		return models.RecordingClip{}, fmt.Errorf("clip with id %s not found", id)
	}

	return clips[0], nil
}

func (p *opensearch) RetrieveTopCapturers(top int, lastPeriods int) ([]string, error) {
	return []string{}, nil
}

func (p *opensearch) RetrieveTopCameras(top int, lastPeriods int) ([]string, error) {
	return []string{}, nil
}

func (p *opensearch) RetrieveTopRegions(top int, lastPeriods int) ([]string, error) {
	return []string{}, nil
}

func (p *opensearch) RetrieveTopLocations(top int, lastPeriods int) ([]string, error) {
	return []string{}, nil
}

func (p *opensearch) RetrieveTopTags(top int, lastPeriods int) ([]string, error) {
	return []string{}, nil
}

func (p *opensearch) RetrieveTopAnalytics(top int, lastPeriods int) ([]string, error) {
	return []string{}, nil
}

func (p *opensearch) RetrieveTopAlertTypes(top int, lastPeriods int) ([]string, error) {
	return []string{}, nil
}

func (p *opensearch) RetrieveClipsByTag(tag string, lastPeriods int) ([]models.RecordingClip, error) {
	return []models.RecordingClip{}, nil
}

func (p *opensearch) RetrieveClipsByTags(tags []string, lastPeriods int) ([]models.RecordingClip, error) {
	return []models.RecordingClip{}, nil
}

func (p *opensearch) RetrieveClipsByAnalytic(analytic string, lastPeriods int) ([]models.RecordingClip, error) {
	return []models.RecordingClip{}, nil
}

func (p *opensearch) RetrieveClipsByAnalytics(analytics []string, lastPeriods int) ([]models.RecordingClip, error) {
	return []models.RecordingClip{}, nil
}

func (p *opensearch) RetrieveClipsByAlertType(typ string, lastPeriods int) ([]models.RecordingClip, error) {
	return []models.RecordingClip{}, nil
}

func (p *opensearch) RetrieveClipsByAlertTypes(typs []string, lastPeriods int) ([]models.RecordingClip, error) {
	return []models.RecordingClip{}, nil
}

func (p *opensearch) Finalize() {
}

func retrieveClips(payload map[string]interface{}) ([]models.RecordingClip, error) {
	clips := []models.RecordingClip{}

	o, err := search[map[string]interface{}, map[string]interface{}](http.MethodGet, "_search", "", payload)
	if err != nil {
		fmt.Printf("😢 error retrieving alerted clips %v\n", err)
		return clips, err
	}

	// WARNING: Need to make sure we check all data elements as this can produce a panic
	outerHits, ok := o["hits"].(map[string]interface{})
	if !ok {
		fmt.Printf("😢 expected hits element to be a map of interface{}....found nothing %v\n", err)
		return clips, nil
	}

	_, ok = outerHits["hits"]
	if !ok {
		fmt.Printf("😢 expected hits element....found nothing %v\n", err)
		return clips, nil
	}

	// This is really strange, why is assertion failing on []map[string]interface{} but not on []interface{}?
	hits, ok := outerHits["hits"].([]interface{})
	if !ok {
		fmt.Printf("😢 expected hits element to be an array of map of interface{}....found nothing %v\n", err)
		return clips, nil
	}

	//hits := o["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		hitValue, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		_, ok = hitValue["_source"]
		if !ok {
			continue
		}

		source, ok := hitValue["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		// Convert the times to time.Time
		createTime, ok := source["createTime"].(string)
		if !ok {
			continue
		}
		createTimeTime, err := time.Parse(time.RFC3339, createTime)
		if err != nil {
			continue
		}

		recordingBeginTime, ok := source["recordingBeginTime"].(string)
		if !ok {
			continue
		}
		recordingBeginTimeTime, err := time.Parse(time.RFC3339, recordingBeginTime)
		if err != nil {
			continue
		}

		recordingEndTime, ok := source["recordingEndTime"].(string)
		if !ok {
			continue
		}
		recordingEndTimeTime, err := time.Parse(time.RFC3339, recordingEndTime)
		if err != nil {
			continue
		}

		publishTime, ok := source["publishTime"].(string)
		if !ok {
			continue
		}
		publishTimeTime, err := time.Parse(time.RFC3339, publishTime)
		if err != nil {
			continue
		}

		modelInvocationBeginTime, ok := source["modelInvocationBeginTime"].(string)
		if !ok {
			continue
		}
		modelInvocationBeginTimeTime, err := time.Parse(time.RFC3339, modelInvocationBeginTime)
		if err != nil {
			continue
		}

		modelInvocationEndTime, ok := source["modelInvocationEndTime"].(string)
		if !ok {
			continue
		}
		modelInvocationEndTimeTime, err := time.Parse(time.RFC3339, modelInvocationEndTime)
		if err != nil {
			continue
		}

		alertInvocationBeginTime, ok := source["alertInvocationBeginTime"].(string)
		if !ok {
			continue
		}
		alertInvocationBeginTimeTime, err := time.Parse(time.RFC3339, alertInvocationBeginTime)
		if err != nil {
			continue
		}

		alertInvocationEndTime, ok := source["alertInvocationEndTime"].(string)
		if !ok {
			continue
		}
		alertInvocationEndTimeTime, err := time.Parse(time.RFC3339, alertInvocationEndTime)
		if err != nil {
			continue
		}

		indexTime, ok := source["indexTime"].(string)
		if !ok {
			continue
		}
		indexTimeTime, err := time.Parse(time.RFC3339, indexTime)
		if err != nil {
			continue
		}

		// Convert tags
		tags, ok := source["tags"].([]interface{})
		if !ok {
			continue
		}
		tagsStr := []string{}
		for _, tag := range tags {
			tagsStr = append(tagsStr, tag.(string))
		}

		// Convert analytics
		analytics, ok := source["analytics"].([]interface{})
		if !ok {
			continue
		}
		analyticsStr := []string{}
		for _, tag := range analytics {
			analyticsStr = append(analyticsStr, tag.(string))
		}

		// Convert alert types
		alertTypes, ok := source["alertTypes"].([]interface{})
		if !ok {
			continue
		}
		alertTypesStr := []string{}
		for _, tag := range alertTypes {
			alertTypesStr = append(alertTypesStr, tag.(string))
		}

		// Convert media index types
		mediaIndexerTypes, ok := source["mediaIndexerTypes"].([]interface{})
		if !ok {
			continue
		}
		mediaIndexerTypesStr := []string{}
		for _, tag := range mediaIndexerTypes {
			mediaIndexerTypesStr = append(mediaIndexerTypesStr, tag.(string))
		}

		clip := models.RecordingClip{
			ID:                       source["id"].(string),
			CreateTime:               createTimeTime,
			LocalReference:           source["localReference"].(string),
			CloudReference:           source["cloudReference"].(string),
			AlertReference:           source["alertReference"].(string),
			StorageProvider:          source["storageProvider"].(string),
			Capturer:                 source["capturer"].(string),
			Camera:                   source["camera"].(string),
			CameraID:                 int(source["cameraId"].(float64)),
			Region:                   source["region"].(string),
			Location:                 source["location"].(string),
			Priority:                 source["priority"].(string),
			Frames:                   int(source["frames"].(float64)),
			RecordingBeginTime:       recordingBeginTimeTime,
			RecordingEndTime:         recordingEndTimeTime,
			PublishTime:              publishTimeTime,
			ModelInvocationBeginTime: modelInvocationBeginTimeTime,
			ModelInvocationEndTime:   modelInvocationEndTimeTime,
			AlertInvocationBeginTime: alertInvocationBeginTimeTime,
			AlertInvocationEndTime:   alertInvocationEndTimeTime,
			IndexTime:                indexTimeTime,
			PrevClip:                 source["prevClip"].(string),
			Analytics:                analyticsStr,
			AlertTypes:               alertTypesStr,
			MediaIndexerTypes:        mediaIndexerTypesStr,
			Tags:                     tagsStr,
			TagsCount:                int(source["tagsCount"].(float64)),
			AlertsCount:              int(source["alertsCount"].(float64)),
			ModelInvoker:             source["modelInvoker"].(string),
			ClipType:                 int(source["clipType"].(float64)),
			RecordingDuration:        int64(source["recordingDuration"].(float64)),
			ModelInvocationDuration:  int64(source["modelInvocationDuration"].(float64)),
			AlertInvocationDuration:  int64(source["alertInvocationDuration"].(float64)),
			CreateToIndexDuration:    int64(source["createToIndexDuration"].(float64)),
		}
		clips = append(clips, clip)
	}

	return clips, nil
}

func search[P any, O any](httpMethod, rest, identifier string, payload P) (O, error) {
	results := new(O)

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

	url := fmt.Sprintf("%s/%s/%s", os.Getenv("OPEN_SEARCH_DOMAIN_ENDPOINT"), os.Getenv("OPEN_SEARCH_INDEX_NAME"), rest)
	if identifier != "" {
		url = fmt.Sprintf("%s/%s", url, identifier)
	}

	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(&payload)
	if err != nil {
		fmt.Printf("😢 error encoding clip %v\n", err)
		return *results, err
	}

	fmt.Printf("Request body: %s\n", payloadBuf.String())

	req, err := http.NewRequest(httpMethod, url, payloadBuf)
	if err != nil {
		fmt.Printf("😢 error creating a new request %v\n", err)
		return *results, err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("😢 error executing request %v\n", err)
		return *results, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("😢 error reading response body %v\n", err)
		return *results, err
	}
	defer resp.Body.Close()

	fmt.Printf("Response body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Printf("😢 error executing request %v\n", err)
		return *results, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err = json.Unmarshal(body, &results)
	if err != nil {
		fmt.Printf("😢 error unmarshalling response %v\n", err)
		return *results, err
	}

	return *results, nil
}
