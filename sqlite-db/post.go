package sqlitedb

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/AliAlhajji/Motarjamat/models"
)

const pageLimit int = 10

type postRepo struct {
	db *sqlx.DB
}

func InitPostRepo() (*postRepo, error) {
	if db == nil {
		return nil, errors.New("database is nil")
	}

	return &postRepo{
		db: db,
	}, nil
}

func (r *postRepo) CreatePost(post *models.Post) (int64, error) {
	q := `INSERT INTO post(title, body, user_id, link, date) VALUES (
		:title, :body, :user_id, :link, :date
	) RETURNING id`

	post.Date = time.Now()

	result, err := r.db.NamedQuery(q, post)
	if err != nil {
		return 0, err
	}

	var id int64

	if result.Next() {
		err = result.Scan(&id)
		if err != nil {
			return 0, err
		}
	}
	defer result.Close()

	return id, nil
}

func (r *postRepo) GetPost(postID int64) (*models.Post, error) {
	q := `SELECT user.uuid user_id, post.id id, title, body, link, username, name 
	FROM post 
	INNER JOIN user ON post.user_id = user.uuid 
	WHERE post.id=$1`
	var post models.Post

	err := r.db.Get(&post, q, postID)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *postRepo) GetAllPostsPaged(page int) ([]*models.Post, error) {
	queryLowerBound := (page - 1) * pageLimit
	queryUpperBound := pageLimit

	q := fmt.Sprintf(`SELECT user.uuid user_id, post.id id, title, body, link, username, name 
	FROM post 
	INNER JOIN user ON post.user_id = user.uuid 
	LIMIT %d,%d`, queryLowerBound, queryUpperBound)
	var posts []*models.Post

	err := r.db.Select(&posts, q)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepo) EditPost(post *models.Post) error {
	q := `UPDATE post SET title = :title , body = :body , link = :link WHERE id = :id AND user_id = :user_id`

	_, err := r.db.NamedExec(q, post)
	if err != nil {
		return err
	}

	return nil
}

func (r *postRepo) DeletePost(postID int64) error {
	q := `DELETE FROM post WHERE id = $1`

	_, err := r.db.Exec(q, postID)
	if err != nil {
		return err
	}

	return nil
}
