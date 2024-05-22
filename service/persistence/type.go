package persistence

import "github.com/khaledhikmat/threat-detection-shared/models"

type IService interface {
	NewClip(clip models.RecordingClip) error

	RetrieveClipCount(lastDays int) (int, error)

	RetrieveClipsStatsByRegion(lastDays int) ([]models.ClipStats, error)
	RetrieveClipsByRegion(region string, page, pageSize int) ([]models.RecordingClip, error)

	RetrieveTopCapturers(top int, lastDays int) ([]string, error)
	RetrieveTopCameras(top int, lastDays int) ([]string, error)
	RetrieveTopRegions(top int, lastDays int) ([]string, error)
	RetrieveTopLocations(top int, lastDays int) ([]string, error)

	RetrieveTopTags(top int, lastDays int) ([]string, error)
	RetrieveTopAnalytics(top int, lastDays int) ([]string, error)
	RetrieveTopAlertTypes(top int, lastDays int) ([]string, error)

	RetrieveClipsByTag(tag string, lastDays int) ([]models.RecordingClip, error)
	RetrieveClipsByTags(tags []string, lastDays int) ([]models.RecordingClip, error)
	RetrieveClipsByAnalytic(analytic string, lastDays int) ([]models.RecordingClip, error)
	RetrieveClipsByAnalytics(analytics []string, lastDays int) ([]models.RecordingClip, error)
	RetrieveClipsByAlertType(typ string, lastDays int) ([]models.RecordingClip, error)
	RetrieveClipsByAlertTypes(typs []string, lastDays int) ([]models.RecordingClip, error)

	Finalize()
}
