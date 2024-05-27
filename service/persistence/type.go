package persistence

import "github.com/khaledhikmat/threat-detection-shared/models"

type IService interface {
	NewClip(clip models.RecordingClip) error

	RetrieveClipCount(lastPeriods int) (int, error)

	RetrieveClipsStatsByRegion(lastPeriods int) ([]models.ClipStats, error)
	RetrieveAlertedClips(top, lastPeriods int) ([]models.RecordingClip, error)
	RetrieveClipsByRegion(region string, lastPeriods, page, pageSize int) ([]models.RecordingClip, error)

	RetrieveTopCapturers(top int, lastPeriods int) ([]string, error)
	RetrieveTopCameras(top int, lastPeriods int) ([]string, error)
	RetrieveTopRegions(top int, lastPeriods int) ([]string, error)
	RetrieveTopLocations(top int, lastPeriods int) ([]string, error)

	RetrieveTopTags(top int, lastPeriods int) ([]string, error)
	RetrieveTopAnalytics(top int, lastPeriods int) ([]string, error)
	RetrieveTopAlertTypes(top int, lastPeriods int) ([]string, error)

	RetrieveClipsByTag(tag string, lastPeriods int) ([]models.RecordingClip, error)
	RetrieveClipsByTags(tags []string, lastPeriods int) ([]models.RecordingClip, error)
	RetrieveClipsByAnalytic(analytic string, lastPeriods int) ([]models.RecordingClip, error)
	RetrieveClipsByAnalytics(analytics []string, lastPeriods int) ([]models.RecordingClip, error)
	RetrieveClipsByAlertType(typ string, lastPeriods int) ([]models.RecordingClip, error)
	RetrieveClipsByAlertTypes(typs []string, lastPeriods int) ([]models.RecordingClip, error)

	Finalize()
}
