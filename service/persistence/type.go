package persistence

import "github.com/khaledhikmat/threat-detection-shared/equates"

type IService interface {
	NewClip(clip equates.RecordingClip) error

	RetrieveTopCapturers(top int, lastDays int) ([]string, error)
	RetrieveTopCameras(top int, lastDays int) ([]string, error)
	RetrieveTopRegions(top int, lastDays int) ([]string, error)
	RetrieveTopLocations(top int, lastDays int) ([]string, error)

	RetrieveTopTags(top int, lastDays int) ([]string, error)
	RetrieveTopAnalytics(top int, lastDays int) ([]string, error)
	RetrieveTopAlertTypes(top int, lastDays int) ([]string, error)

	RetrieveClipsByTag(tag string, lastDays int) ([]equates.RecordingClip, error)
	RetrieveClipsByTags(tags []string, lastDays int) ([]equates.RecordingClip, error)
	RetrieveClipsByAnalytic(analytic string, lastDays int) ([]equates.RecordingClip, error)
	RetrieveClipsByAnalytics(analytics []string, lastDays int) ([]equates.RecordingClip, error)
	RetrieveClipsByAlertType(typ string, lastDays int) ([]equates.RecordingClip, error)
	RetrieveClipsByAlertTypes(typs []string, lastDays int) ([]equates.RecordingClip, error)

	Finalize()
}
