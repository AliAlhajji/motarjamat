package sqlitedb

import (
	"errors"

	"github.com/AliAlhajji/Motarjamat/models"
	"github.com/jmoiron/sqlx"
)

type settingsRepo struct {
	db *sqlx.DB
}

func InitSettingsRepo() (*settingsRepo, error) {
	if db == nil {
		return nil, errors.New("database is nil")
	}
	return &settingsRepo{
		db: db,
	}, nil
}

func (r *settingsRepo) GetSettings() (*models.SiteSettings, error) {
	q := `SELECT * FROM site_settings`
	var settings models.SiteSettings

	err := r.db.Get(&settings, q)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func (r *settingsRepo) Update(settings *models.SiteSettings) error {
	q := `UPDATE site_settings SET
	 title = :title, 
	is_running = :is_running, 
	announcement = :announcement, 
	results_per_page = :results_per_page
	`

	_, err := r.db.NamedExec(q, settings)
	if err != nil {
		return err
	}

	return nil
}
