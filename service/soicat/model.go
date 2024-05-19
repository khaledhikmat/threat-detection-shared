package soicat

type Camera struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Region             string   `json:"region"`   // Home Office
	Location           string   `json:"location"` // B03W-E1 or C03E-EXEC
	Priority           string   `json:"priority"` // ATM, Critical Infrastructure, Transportation
	RtspURL            string   `json:"rtspUrl"`
	Analytics          []string `json:"analytics"`
	AlertTypes         []string `json:"alertTypes"`
	MediaIndexerTypes  []string `json:"mediaIndexerTypes"`
	Capturer           string   `json:"capturer"`
	RecordingState     string   `json:"recordingState"`    // started, stopped, paused
	RetentionDays      int      `json:"retentionDays"`     // days to keep initial recording clips
	DeepRetentionDays  int      `json:"deepRetentionDays"` // days to keep subsequent recording clips
	CaptureWidth       int      `json:"captureWidth"`
	CaptureHeight      int      `json:"captureHeight"`
	PreRecording       int64    `json:"preRecording"`
	MaxLengthRecording int64    `json:"maxLengthRecording"`
	Timezone           string   `json:"timezone"`
}
