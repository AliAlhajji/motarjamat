package sqlitedb

import (
	"database/sql"
	"errors"

	"github.com/AliAlhajji/Motarjamat/models"
	"github.com/jmoiron/sqlx"
)

type categoryRepo struct {
	db *sqlx.DB
}

func InitCategoryRepo() (*categoryRepo, error) {
	if db == nil {
		return nil, errors.New("database is nil")
	}

	return &categoryRepo{
		db: db,
	}, nil
}

func (r *categoryRepo) Add(title string) error {
	q := `INSERT INTO category (title) VALUES ($1)`

	_, err := r.db.Exec(q, title)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepo) UpdateCategory(categoryID int64, updatedCategory *models.Category) error {
	if updatedCategory.Title == "" {
		return errors.New("title cannot be empty")
	}

	q := `UPDATE category SET title = $1 WHERE id = $2`
	_, err := r.db.Exec(q, updatedCategory.Title, categoryID)
	if err != nil {
		return err
	}
	return nil
}

func (r *categoryRepo) DeleteCategory(categoryID int64) error {
	q := `DELETE FROM category WHERE id = $1`

	_, err := r.db.Exec(q, categoryID)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepo) GetAll() ([]*models.Category, error) {
	q := `SELECT id, title FROM category`

	var categories []*models.Category

	err := r.db.Select(&categories, q)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepo) GetUsedCategories() ([]*models.Category, error) {
	q := `SELECT id, title FROM category WHERE id in (SELECT DISTINCT category_id FROM post_category)`

	var categories []*models.Category

	err := r.db.Select(&categories, q)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepo) GetCategory(id int64) (*models.Category, error) {
	q := `SELECT id, title FROM category WHERE id = $1`

	var category models.Category

	err := r.db.Get(&category, q, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &category, nil
}
