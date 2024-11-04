package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"

	"ugames/internal/config"
	"ugames/internal/models"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepository(cnf *config.Cnf) *Repo {
	pool, err := NewPgxPool(context.Background(), cnf)
	if err != nil {
		log.Error().Msg("[PGXPOOL]: " + err.Error())
	}

	return &Repo{db: pool}
}

func NewPgxPool(ctx context.Context, cnf *config.Cnf) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cnf.Db.User, cnf.Db.Pass, cnf.Db.Host, cnf.Db.Port, cnf.Db.Name)
	log.Printf(dsn)
	pgConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Error().Msg("[PGXPOOL]: " + err.Error())
	}

	if err != nil {
		log.Error().Msg("[PGXPOOL]: " + err.Error())
	}

	pool, err := pgxpool.ConnectConfig(ctx, pgConfig)
	if err != nil {
		log.Error().Msg("[PGXPOOL]: " + err.Error())
	}

	log.Info().Msg("Database connected!")

	return pool, nil
}

func (repo *Repo) GetKeyWordsList() ([]models.KeyWord, error) {
	sql := `SELECT id, key_word FROM ugames.key_words`

	var data []models.KeyWord
	rows, err := repo.db.Query(context.Background(), sql)
	if err != nil {
		log.Error().Msg("[PGXPOOL] Keywords select: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var g models.KeyWord
		err = rows.Scan(&g.Id, &g.KeyWord)
		if err != nil {
			log.Error().Msg("[PGXPOOL] Keywords rows scan: " + err.Error())
		}
		data = append(data, g)
	}

	return data, err
}

func (repo *Repo) GetUncheckedRepos() ([]models.Repos, error) {
	sql := `SELECT id, repo_name, homepage, is_checked FROM ugames.repos WHERE is_checked isnull OR is_checked = false  ORDER BY created_at DESC`
	var data []models.Repos
	rows, err := repo.db.Query(context.Background(), sql)
	if err != nil {
		log.Error().Msg("[PGXPOOL] GetUncheckedRepos select: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var g models.Repos
		err = rows.Scan(&g.Id, &g.RepoName, &g.Homepage, &g.IsChecked)
		if err != nil {
			log.Error().Msg("[PGXPOOL] GetUncheckedRepos rows scan: " + err.Error())
		}
		data = append(data, g)
	}

	return data, err
}

func (repo *Repo) GetCheckedRepos() ([]models.Repos, error) {
	sql := `SELECT id, key_word, repo_name, homepage, content, comment, created_at FROM ugames.repos WHERE is_checked = true AND homepage <> '' OR content <> '' ORDER BY created_at DESC`
	var data []models.Repos
	rows, err := repo.db.Query(context.Background(), sql)
	if err != nil {
		log.Error().Msg("[PGXPOOL] GetCheckedRepos select: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var g models.Repos
		err = rows.Scan(&g.Id, &g.KeyWord, &g.RepoName, &g.Homepage, &g.Content, &g.Comment, &g.CreatedAt)
		if err != nil {
			log.Error().Msg("[PGXPOOL] GetCheckedRepos rows scan: " + err.Error())
		}
		data = append(data, g)
	}

	return data, err
}

func (repo *Repo) UpdateCheckedRepo(repos models.Repos) error {
	_, err := repo.db.Exec(context.Background(), "UPDATE ugames.repos SET content=$1, is_checked = true WHERE id=$2", repos.Content, repos.Id)
	if err != nil {
		log.Error().Msg("[PGXPOOL] UpdateCheckedRepo update: " + err.Error())
		return err
	}
	return nil
}

func (repo *Repo) AddComment(comment models.ReqComment) error {
	_, err := repo.db.Exec(context.Background(), "UPDATE ugames.repos SET comment=$1 WHERE id=$2", comment.Comment, comment.Id)
	if err != nil {
		log.Error().Msg("[PGXPOOL] AddComment update: " + err.Error())
		return err
	}
	return nil
}

func (repo *Repo) InsertRepo(repoName, homePage, keyWord string) {
	sqlStatement := `INSERT INTO ugames.repos (repo_name, homePage, key_word) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := repo.db.Exec(context.Background(), sqlStatement, repoName, homePage, keyWord)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}
