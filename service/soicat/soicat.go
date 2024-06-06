package soicat

import (
	"fmt"

	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

var localCameras = []Camera{
	{
		ID:                 "100",
		Name:               "Camera1",
		Region:             "Home_Office",
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

var awsCameras = []Camera{
	{
		ID:                 "100",
		Name:               "Camera1",
		Region:             "Home_Office",
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
	{
		ID:                 "200",
		Name:               "PHX AB JTRON QUAD CAM 1 (8240)",
		Region:             "Phoenix",
		Location:           "B01W-C1",
		Priority:           "Normal",
		RtspURL:            "rtsp://admin:gooze_bumbs@192.168.1.206:554/cam/realmonitor?channel=1&subtype=0",
		Analytics:          []string{"weapon"},
		AlertTypes:         []string{"ccure", "snow", "pers"},
		MediaIndexerTypes:  []string{"elastic"},
		Capturer:           "dead-capturer",
		RetentionDays:      1,
		DeepRetentionDays:  30,
		CaptureWidth:       500,
		CaptureHeight:      500,
		MaxLengthRecording: 3,
		Timezone:           "",
	},
	{
		ID:                 "300",
		Name:               "PLANO A ROOF NE PTZ (2160)",
		Region:             "Plano",
		Location:           "E03-FLOOR",
		Priority:           "Normal",
		RtspURL:            "rtsp://admin:gooze_bumbs@192.168.1.206:554/cam/realmonitor?channel=1&subtype=0",
		Analytics:          []string{"weapon"},
		AlertTypes:         []string{"ccure", "snow", "pers"},
		MediaIndexerTypes:  []string{"elastic"},
		Capturer:           "dead-capturer",
		RetentionDays:      1,
		DeepRetentionDays:  30,
		CaptureWidth:       500,
		CaptureHeight:      500,
		MaxLengthRecording: 3,
		Timezone:           "",
	},
	{
		ID:                 "400",
		Name:               "CENTS 10TH FLR E STAIR ENTRANCE (7528)",
		Region:             "Charlotte",
		Location:           "C-10W",
		Priority:           "Normal",
		RtspURL:            "rtsp://admin:gooze_bumbs@192.168.1.206:554/cam/realmonitor?channel=1&subtype=0",
		Analytics:          []string{"weapon"},
		AlertTypes:         []string{"ccure", "snow", "pers"},
		MediaIndexerTypes:  []string{"elastic"},
		Capturer:           "dead-capturer",
		RetentionDays:      1,
		DeepRetentionDays:  30,
		CaptureWidth:       500,
		CaptureHeight:      500,
		MaxLengthRecording: 3,
		Timezone:           "",
	},
}

func New(cfgsvc config.IService) IService {
	return &soicat{
		CfgSvc: cfgsvc,
	}
}

type soicat struct {
	CfgSvc config.IService
}

func (s *soicat) UpdateCamera(c Camera) error {
	existingCameras, err := s.Cameras()
	if err != nil {
		return err
	}

	for idx, camera := range existingCameras {
		if camera.Name == c.Name {
			ptr := &existingCameras[idx]
			(*ptr).Capturer = c.Capturer
			return nil
		}
	}

	return fmt.Errorf("could not find the camera %s", c.Name)
}

func (s *soicat) Cameras() ([]Camera, error) {
	if s.CfgSvc.GetRuntimeEnv() == "local" {
		return localCameras, nil
	}
	return awsCameras, nil
}

func (s *soicat) Finalize() {
}
