package equates

import (
	"github.com/guregu/null"
)

const (
	ThreatDetectionSecrets    = "threat-detection-secrets"
	ThreatDetectionStateStore = "threat-detection-statestore" // name must match config/redis-statestore.yaml
	ThreatDetectionPubSub     = "threat-detection-pubsub"     // name must match config/redis-pubsub.yaml
	RecordingsTopic           = "recordings-topic"
	NotificationsTopic        = "notifications-topic"
	MetadataTopic             = "metadata-topic"
)

type RecordingClip struct {
	ID              string    `json:"id"`
	LocalReference  string    `json:"localReference"`
	CloudReference  string    `json:"cloudReference"`
	StorageProvider string    `json:"storageProvider"`
	Capturer        string    `json:"capturer"`
	Camera          string    `json:"camera"`
	Frames          int       `json:"frames"`
	BeginTime       null.Time `json:"beginTime"`
	EndTime         null.Time `json:"endTime"`
	PrevClip        string    `json:"prevClip"`  // The ID of the previous clip in the sequence
	Analytics       []string  `json:"analytics"` // The analytics that were applied on this clip
	Tags            []string  `json:"tags"`      // The tags that were detected on this clip
}
