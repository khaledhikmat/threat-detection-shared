package soicat

import (
	"fmt"
)

var cameras = []Camera{
	{
		ID:                 "100",
		Name:               "Camera1",
		Region:             "Home Office",
		Location:           "B03W-E1",
		Priority:           "ATM",
		RtspURL:            "rtsp://admin:gooze_bumbs@192.168.1.206:554/cam/realmonitor?channel=1&subtype=0",
		Analytics:          []string{"weapon", "fire"},
		AlertTypes:         []string{"ccure", "snow", "pers", "slack"},
		MediaIndexerTypes:  []string{"database", "elastic"},
		Capturer:           "dead-capturer",
		RetentionDays:      1,
		DeepRetentionDays:  30,
		CaptureWidth:       500,
		CaptureHeight:      500,
		MaxLengthRecording: 3,
		Timezone:           "",
	},
}

func New() IService {
	return &soicat{}
}

type soicat struct {
}

func (s *soicat) UpdateCamera(c Camera) error {
	for idx, camera := range cameras {
		if camera.Name == c.Name {
			ptr := &cameras[idx]
			(*ptr).Capturer = c.Capturer
			return nil
		}
	}

	return fmt.Errorf("could not find the camera %s", c.Name)
}

func (s *soicat) Cameras() ([]Camera, error) {
	return cameras, nil
}

func (s *soicat) Finalize() {
}
