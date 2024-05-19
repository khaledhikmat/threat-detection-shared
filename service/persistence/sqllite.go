package persistence

import (
	"database/sql"
	_ "embed"
	"os"
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
	db, err := sql.Open("sqlite3", os.Getenv("SQLLITE_FILE_PATH"))
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
	VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`,
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
		strings.Join(clip.Tags, ","),
		clip.TagsCount,
		clip.AlertsCount,
	)
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

func (p *sqllite) RetrieveClipCount(lastDays int) (int, error) {
	q := `SELECT COUNT(*) as clips 
	FROM clips`

	var row *sql.Row
	if lastDays == -1 {
		q += ";"
		row = p.db.QueryRow(q)
	} else {
		q += ` WHERE julianday('now') - julianday(beginTime) <= ?;`
		row = p.db.QueryRow(q, lastDays)
	}

	var err error
	var count = 0
	if err = row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (p *sqllite) RetrieveClipsStatsByRegion(lastDays int) ([]equates.ClipStats, error) {
	q := `SELECT region, COUNT(DISTINCT camera) as cameras, COUNT(*) as clips, 
	SUM(frames) as frames, SUM(tagsCount) as tags, SUM(alertsCount) as alerts
	FROM clips
	WHERE julianday('now') - julianday(beginTime) <= ?
	GROUP BY region;`
	rows, err := p.db.Query(q, lastDays)
	if err != nil {
		return []equates.ClipStats{}, err
	}
	defer rows.Close()

	stats := []equates.ClipStats{}
	for rows.Next() {
		var region string
		var cameras, clips, frames, tags, alerts int
		if err := rows.Scan(&region, &cameras, &clips, &frames, &tags, &alerts); err != nil {
			return []equates.ClipStats{}, err
		}
		stats = append(stats, equates.ClipStats{
			Region:  region,
			Cameras: cameras,
			Clips:   clips,
			Frames:  frames,
			Tags:    tags,
			Alerts:  alerts,
		})
	}

	return stats, nil
}

func (p *sqllite) RetrieveClipsByRegion(region string, page, pageSize int) ([]equates.RecordingClip, error) {
	q := `SELECT id, cloudReference, storageProvider, capturer, camera, region, location, priority, frames, tagsCount, alertsCount 
	FROM clips 
	WHERE region = ? LIMIT ? OFFSET ?;`
	rows, err := p.db.Query(q, region, pageSize, page*pageSize)
	if err != nil {
		return []equates.RecordingClip{}, err
	}
	defer rows.Close()

	clips := []equates.RecordingClip{}
	for rows.Next() {
		var clip equates.RecordingClip
		if err := rows.Scan(&clip.ID, &clip.CloudReference, &clip.StorageProvider, &clip.Capturer,
			&clip.Camera, &clip.Region, &clip.Location, &clip.Priority,
			&clip.Frames, &clip.TagsCount, &clip.AlertsCount); err != nil {
			return []equates.RecordingClip{}, err
		}
		clips = append(clips, clip)
	}

	return clips, nil
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
