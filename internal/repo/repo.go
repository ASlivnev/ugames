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

func (repo *Repo) DeallocateAll() error {
	sql := `DEALLOCATE PREPARE ALL;`
	query, err := repo.db.Query(context.Background(), sql)
	if err != nil {
		log.Error().Msg("[PGXPOOL] DeallocateAll : " + err.Error())
	}
	defer query.Close()

	return err
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
	sql := `SELECT id, key_word, repo_name, homepage, content, comment, created_at FROM ugames.repos WHERE is_checked = true AND (homepage <> '' OR content <> '')  AND (comment <> '-' OR comment IS NULL) ORDER BY created_at DESC;`
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

func (repo *Repo) InsertGmGame(game models.GmGame) {
	sqlStatement := `INSERT INTO ugames.gm_games (category, description, tags, thumb, title, url, instructions) VALUES ($1,$2,$3,$4, $5,$6,$7) ON CONFLICT DO NOTHING`
	_, err := repo.db.Exec(context.Background(), sqlStatement, game.Category, game.Description, game.Tags, game.Thumb, game.Title, game.URL, game.Instructions)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func (repo *Repo) GetGmGames() ([]models.GmGame, error) {
	sql := `SELECT id,url FROM ugames.gm_games WHERE is_construct is null LIMIT 1000;`
	var data []models.GmGame
	rows, err := repo.db.Query(context.Background(), sql)
	if err != nil {
		log.Error().Msg("[PGXPOOL] GetGmGames select: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var g models.GmGame
		err = rows.Scan(&g.ID, &g.URL)
		if err != nil {
			log.Error().Msg("[PGXPOOL] GetGmGames rows scan: " + err.Error())
		}
		data = append(data, g)
	}

	return data, nil
}

func (repo *Repo) UpdateGmGame(game models.GmGame) error {
	_, err := repo.db.Exec(context.Background(), "UPDATE ugames.gm_games SET is_construct=$2 WHERE id=$1", game.ID, game.IsConstruct)
	if err != nil {
		log.Error().Msg("[PGXPOOL] UpdateGmGame update: " + err.Error())
		return err
	}
	return nil
}

func (repo *Repo) GetConstructGames(filter string) ([]models.GmGame, error) {
	var sql string
	if filter == "all" {
		sql = `SELECT * FROM ugames.gm_games WHERE is_construct='Y' AND list is null ORDER BY updated_at DESC LIMIT 50;`
	}

	if filter == "published" {
		sql = `SELECT * FROM ugames.gm_games WHERE is_construct='Y' AND list = 'published' ORDER BY updated_at DESC ;`
	}

	if filter == "white" {
		sql = `SELECT * FROM ugames.gm_games WHERE is_construct='Y' AND list = 'white' ORDER BY updated_at DESC ;`
	}

	if filter == "grey" {
		sql = `SELECT * FROM ugames.gm_games WHERE is_construct='Y' AND list = 'grey' ORDER BY updated_at DESC ;`
	}

	if filter == "black" {
		sql = `SELECT * FROM ugames.gm_games WHERE is_construct='Y' AND list = 'black' ORDER BY updated_at DESC ;`
	}

	var data []models.GmGame
	rows, err := repo.db.Query(context.Background(), sql)
	if err != nil {
		log.Error().Msg("[PGXPOOL] GetGmGames select: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var g models.GmGame
		err = rows.Scan(&g.ID, &g.Title, &g.Description, &g.Instructions, &g.URL, &g.Category, &g.Tags, &g.Thumb, &g.IsConstruct, &g.List, &g.Comment, &g.UpdatedAt)
		if err != nil {
			log.Error().Msg("[PGXPOOL] GetConstructGames rows scan: " + err.Error())
		}
		data = append(data, g)
	}

	return data, nil
}

func (repo *Repo) AddCommentC3(comment models.ReqCommentC3) error {
	_, err := repo.db.Exec(context.Background(), "UPDATE ugames.gm_games SET comment=$1, list = $3, updated_at = NOW() WHERE id=$2", comment.Comment, comment.Id, comment.List)
	if err != nil {
		log.Error().Msg("[PGXPOOL] AddCommentC3 update: " + err.Error())
		return err
	}
	return nil
}

func (repo *Repo) GetC3GamesSearch(searchReq string) ([]models.GmGame, error) {
	sql := `SELECT id, title, description, url, thumb FROM ugames.gm_games 
            WHERE title ILIKE '%' || $1 || '%' 
              AND is_construct = 'Y' 
            ORDER BY updated_at DESC 
            LIMIT 500;`

	var data []models.GmGame
	rows, err := repo.db.Query(context.Background(), sql, searchReq)
	if err != nil {
		log.Error().Msg("[PGXPOOL] GetC3GamesSearch select: " + err.Error())
		return nil, err // Возвращаем ошибку
	}
	defer rows.Close()

	for rows.Next() {
		var g models.GmGame
		err = rows.Scan(&g.ID, &g.Title, &g.Description, &g.URL, &g.Thumb)
		if err != nil {
			log.Error().Msg("[PGXPOOL] GetC3GamesSearch rows scan: " + err.Error())
			return nil, err // Возвращаем ошибку при сканировании
		}
		data = append(data, g)
	}

	// Проверяем ошибки после итерации
	if rows.Err() != nil {
		log.Error().Msg("[PGXPOOL] GetGamesSearch rows iteration: " + rows.Err().Error())
		return nil, rows.Err()
	}

	return data, nil
}
