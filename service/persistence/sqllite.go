package persistence

import (
	"database/sql"
	_ "embed"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/khaledhikmat/threat-detection-shared/equates"
)

//go:embed sql/createclipstable4lite.sql
var createClipTable string

//go:embed sql/createcliptagstable4lite.sql
var createClipTagsTable string

//go:embed sql/createclipanalyticstable4lite.sql
var createClipAnalyticsTable string

//go:embed sql/createclipalerttypestable4lite.sql
var createClipAlertTypesTable string

//go:embed sql/createclipindextypestable4lite.sql
var createClipIndexTypesTable string

func newSqllite() IService {
	db, err := sql.Open("sqlite3", "./clips.db")
	if err != nil {
		panic(err)
	}

	if _, err := db.Exec(createClipTable); err != nil {
		panic(err)
	}

	if _, err := db.Exec(createClipTagsTable); err != nil {
		panic(err)
	}

	if _, err := db.Exec(createClipAnalyticsTable); err != nil {
		panic(err)
	}

	if _, err := db.Exec(createClipAlertTypesTable); err != nil {
		panic(err)
	}

	if _, err := db.Exec(createClipIndexTypesTable); err != nil {
		panic(err)
	}

	return &sqllite{
		db: db,
	}
}

type sqllite struct {
	db *sql.DB
}

func (p *sqllite) NewClip(clip equates.RecordingClip) error {
	bts, err := time.Parse(equates.Layout, clip.BeginTime)
	if err != nil {
		return err
	}

	ets, err := time.Parse(equates.Layout, clip.EndTime)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(`INSERT INTO clips 
	VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`,
		clip.ID,
		clip.LocalReference,
		clip.CloudReference,
		clip.StorageProvider,
		clip.Capturer,
		clip.Camera,
		clip.Region,
		clip.Location,
		clip.Priority,
		clip.Frames,
		bts,
		ets,
		clip.PrevClip,
		strings.Join(clip.Analytics, ","),
		strings.Join(clip.AlertTypes, ","),
		strings.Join(clip.MediaIndexerTypes, ","),
		strings.Join(clip.Tags, ","))
	if err != nil {
		return err
	}

	for _, tag := range clip.Tags {
		_, err = p.db.Exec(`INSERT INTO cliptags 
		VALUES(NULL,?,?,?);`,
			clip.ID,
			tag,
			time.Now())
		if err != nil {
			return err
		}
	}

	for _, analytic := range clip.Analytics {
		_, err = p.db.Exec(`INSERT INTO clipanalytics 
		VALUES(NULL,?,?,?);`,
			clip.ID,
			analytic,
			time.Now())
		if err != nil {
			return err
		}
	}

	for _, alert := range clip.AlertTypes {
		_, err = p.db.Exec(`INSERT INTO clipalerttypes 
		VALUES(NULL,?,?,?);`,
			clip.ID,
			alert,
			time.Now())
		if err != nil {
			return err
		}
	}

	for _, idx := range clip.MediaIndexerTypes {
		_, err = p.db.Exec(`INSERT INTO clipindextypes 
		VALUES(NULL,?,?,?);`,
			clip.ID,
			idx,
			time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *sqllite) RetrieveTopCapturers(top int, lastDays int) ([]string, error) {
	return []string{}, nil
}

func (p *sqllite) RetrieveTopCameras(top int, lastDays int) ([]string, error) {
	return []string{}, nil
}

func (p *sqllite) RetrieveTopRegions(top int, lastDays int) ([]string, error) {
	return []string{}, nil
}

func (p *sqllite) RetrieveTopLocations(top int, lastDays int) ([]string, error) {
	return []string{}, nil
}

func (p *sqllite) RetrieveTopTags(top int, lastDays int) ([]string, error) {
	return []string{}, nil
}

func (p *sqllite) RetrieveTopAnalytics(top int, lastDays int) ([]string, error) {
	return []string{}, nil
}

func (p *sqllite) RetrieveTopAlertTypes(top int, lastDays int) ([]string, error) {
	return []string{}, nil
}

func (p *sqllite) RetrieveClipsByTag(tag string, lastDays int) ([]equates.RecordingClip, error) {
	return []equates.RecordingClip{}, nil
}

func (p *sqllite) RetrieveClipsByTags(tags []string, lastDays int) ([]equates.RecordingClip, error) {
	return []equates.RecordingClip{}, nil
}

func (p *sqllite) RetrieveClipsByAnalytic(analytic string, lastDays int) ([]equates.RecordingClip, error) {
	return []equates.RecordingClip{}, nil
}

func (p *sqllite) RetrieveClipsByAnalytics(analytics []string, lastDays int) ([]equates.RecordingClip, error) {
	return []equates.RecordingClip{}, nil
}

func (p *sqllite) RetrieveClipsByAlertType(typ string, lastDays int) ([]equates.RecordingClip, error) {
	return []equates.RecordingClip{}, nil
}

func (p *sqllite) RetrieveClipsByAlertTypes(typs []string, lastDays int) ([]equates.RecordingClip, error) {
	return []equates.RecordingClip{}, nil
}

func (p *sqllite) Finalize() {
}
