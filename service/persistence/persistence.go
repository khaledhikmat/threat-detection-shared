package persistence

import (
	"fmt"

	"github.com/khaledhikmat/threat-detection-shared/models"
	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

var providers map[string]IService

func New(cfgsvc config.IService) IService {
	providers = map[string]IService{
		"sqllite":    newSqllite(),
		"opensearch": newOpenSearch(),
	}

	return &persistence{
		CfgSvc: cfgsvc,
	}
}

type persistence struct {
	CfgSvc config.IService
}

func (p *persistence) NewClip(clip models.RecordingClip) error {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.NewClip(clip)
}

func (p *persistence) RetrieveClipCount(lastDays int) (int, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return 0, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveClipCount(lastDays)
}

func (p *persistence) RetrieveClipsStatsByRegion(lastDays int) ([]models.ClipStats, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []models.ClipStats{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveClipsStatsByRegion(lastDays)
}

func (p *persistence) RetrieveClipByID(id string) (models.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return models.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveClipByID(id)
}

func (p *persistence) RetrieveAlertedClips(top, lastPeriods int) ([]models.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []models.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveAlertedClips(top, lastPeriods)
}

func (p *persistence) RetrieveClipsByRegion(region string, lastPeriods, page, pageSize int) ([]models.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []models.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveClipsByRegion(region, lastPeriods, page, pageSize)
}

func (p *persistence) RetrieveTopCapturers(top int, lastDays int) ([]string, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []string{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveTopCapturers(top, lastDays)
}

func (p *persistence) RetrieveTopCameras(top int, lastDays int) ([]string, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []string{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveTopCameras(top, lastDays)
}

func (p *persistence) RetrieveTopRegions(top int, lastDays int) ([]string, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []string{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveTopRegions(top, lastDays)
}

func (p *persistence) RetrieveTopLocations(top int, lastDays int) ([]string, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []string{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveTopLocations(top, lastDays)
}

func (p *persistence) RetrieveTopTags(top int, lastDays int) ([]string, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []string{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveTopTags(top, lastDays)
}

func (p *persistence) RetrieveTopAnalytics(top int, lastDays int) ([]string, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []string{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveTopAnalytics(top, lastDays)
}

func (p *persistence) RetrieveTopAlertTypes(top int, lastDays int) ([]string, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []string{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveTopAlertTypes(top, lastDays)
}

func (p *persistence) RetrieveClipsByTag(tag string, lastDays int) ([]models.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []models.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveClipsByTag(tag, lastDays)
}

func (p *persistence) RetrieveClipsByTags(tags []string, lastDays int) ([]models.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []models.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveClipsByTags(tags, lastDays)
}

func (p *persistence) RetrieveClipsByAnalytic(analytic string, lastDays int) ([]models.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []models.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveClipsByAnalytic(analytic, lastDays)
}

func (p *persistence) RetrieveClipsByAnalytics(analytics []string, lastDays int) ([]models.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []models.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveClipsByAnalytics(analytics, lastDays)
}

func (p *persistence) RetrieveClipsByAlertType(typ string, lastDays int) ([]models.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []models.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveClipsByAlertType(typ, lastDays)
}

func (p *persistence) RetrieveClipsByAlertTypes(typs []string, lastDays int) ([]models.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexerType()]
	if !ok {
		return []models.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexerType())
	}

	return r.RetrieveClipsByAlertTypes(typs, lastDays)
}

func (p *persistence) Finalize() {
}
