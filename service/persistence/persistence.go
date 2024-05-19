package persistence

import (
	"fmt"

	"github.com/khaledhikmat/threat-detection-shared/equates"
	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

var providers = map[string]IService{
	"sqllite": newSqllite(),
}

func New(cfgsvc config.IService) IService {
	return &persistence{
		CfgSvc: cfgsvc,
	}
}

type persistence struct {
	CfgSvc config.IService
}

func (p *persistence) NewClip(clip equates.RecordingClip) error {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.NewClip(clip)
}

func (p *persistence) RetrieveClipCount(lastDays int) (int, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return 0, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveClipCount(lastDays)
}

func (p *persistence) RetrieveClipsStatsByRegion(lastDays int) ([]equates.ClipStats, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return []equates.ClipStats{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveClipsStatsByRegion(lastDays)
}

func (p *persistence) RetrieveClipsByRegion(region string, page, pageSize int) ([]equates.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return []equates.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveClipsByRegion(region, page, pageSize)
}

func (p *persistence) RetrieveTopCapturers(top int, lastDays int) ([]string, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return []string{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveTopCapturers(top, lastDays)
}

func (p *persistence) RetrieveTopCameras(top int, lastDays int) ([]string, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return []string{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveTopCameras(top, lastDays)
}

func (p *persistence) RetrieveTopRegions(top int, lastDays int) ([]string, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return []string{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveTopRegions(top, lastDays)
}

func (p *persistence) RetrieveTopLocations(top int, lastDays int) ([]string, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return []string{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveTopLocations(top, lastDays)
}

func (p *persistence) RetrieveTopTags(top int, lastDays int) ([]string, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return []string{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveTopTags(top, lastDays)
}

func (p *persistence) RetrieveTopAnalytics(top int, lastDays int) ([]string, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return []string{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveTopAnalytics(top, lastDays)
}

func (p *persistence) RetrieveTopAlertTypes(top int, lastDays int) ([]string, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return []string{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveTopAlertTypes(top, lastDays)
}

func (p *persistence) RetrieveClipsByTag(tag string, lastDays int) ([]equates.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return []equates.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveClipsByTag(tag, lastDays)
}

func (p *persistence) RetrieveClipsByTags(tags []string, lastDays int) ([]equates.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return []equates.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveClipsByTags(tags, lastDays)
}

func (p *persistence) RetrieveClipsByAnalytic(analytic string, lastDays int) ([]equates.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return []equates.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveClipsByAnalytic(analytic, lastDays)
}

func (p *persistence) RetrieveClipsByAnalytics(analytics []string, lastDays int) ([]equates.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return []equates.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveClipsByAnalytics(analytics, lastDays)
}

func (p *persistence) RetrieveClipsByAlertType(typ string, lastDays int) ([]equates.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return []equates.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveClipsByAlertType(typ, lastDays)
}

func (p *persistence) RetrieveClipsByAlertTypes(typs []string, lastDays int) ([]equates.RecordingClip, error) {
	r, ok := providers[p.CfgSvc.GetIndexProvider()]
	if !ok {
		return []equates.RecordingClip{}, fmt.Errorf("provider %s not supported", p.CfgSvc.GetIndexProvider())
	}

	return r.RetrieveClipsByAlertTypes(typs, lastDays)
}

func (p *persistence) Finalize() {
}
