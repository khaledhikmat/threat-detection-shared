package models

import "time"

const (
	Layout                    = "2006-01-02 15:04:05"
	ThreatDetectionSecrets    = "threat-detection-secrets"
	ThreatDetectionStateStore = "threat-detection-statestore" // name must match config/redis-statestore.yaml
	ThreatDetectionPubSub     = "threat-detection-pubsub"     // name must match config/redis-pubsub.yaml
	RecordingsTopic           = "recordings-topic"
	AlertsTopic               = "alerts-topic"
	MetadataTopic             = "metadata-topic"
)

type RecordingClip struct {
	ID                string    `json:"id"`
	TimeStamp         time.Time `json:"timeStamp"`
	LocalReference    string    `json:"localReference"`
	CloudReference    string    `json:"cloudReference"`
	StorageProvider   string    `json:"storageProvider"`
	Capturer          string    `json:"capturer"`
	Camera            string    `json:"camera"`
	CameraID          int       `json:"cameraId"` // To give us a numeric field we can aggregate on
	Region            string    `json:"region"`   // Home Office
	Location          string    `json:"location"` // B03W-E1 or C03E-EXEC
	Priority          string    `json:"priority"` // ATM, Critical Infrastructure, Transportation
	Frames            int       `json:"frames"`
	BeginTime         string    `json:"beginTime"`
	EndTime           string    `json:"endTime"`
	PrevClip          string    `json:"prevClip"`          // The ID of the previous clip in the sequence
	Analytics         []string  `json:"analytics"`         // The analytics that were applied on this clip
	AlertTypes        []string  `json:"alertTypes"`        // The alert types that are required on this clip
	MediaIndexerTypes []string  `json:"mediaIndexerTypes"` // The media indexer types that are required for this clip
	Tags              []string  `json:"tags"`              // The tags that were detected on this clip
	TagsCount         int       `json:"tagsCount"`         // The number of tags that were detected on this clip
	AlertsCount       int       `json:"alertsCount"`       // Whether the clip is an alert or not
}

type ClipStats struct {
	Region  string `json:"region"`
	Cameras int    `json:"cameras"`
	Clips   int    `json:"clips"`
	Frames  int    `json:"frames"`
	Tags    int    `json:"tags"`
	Alerts  int    `json:"alerts"`
}
