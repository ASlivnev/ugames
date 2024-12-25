package models

import "time"

type KeyWord struct {
	Id      int    `json:"id"`
	KeyWord string `json:"key_word"`
}

type KeyWordReq struct {
	KeyWord string `json:"key_word"`
}

type Resp struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type GitReposRespItem struct {
	FullName string `json:"full_name"`
	Homepage string `json:"homepage"`
}

type GitReposResp struct {
	TotalCount int                `json:"total_count"`
	Items      []GitReposRespItem `json:"items"`
}

type GitReposReadme struct {
	Content string `json:"content"`
}

type Repos struct {
	Id        int        `json:"id"`
	KeyWord   *string    `json:"key_word"`
	RepoName  *string    `json:"repo_name"`
	Homepage  *string    `json:"homepage"`
	Content   *string    `json:"content"`
	Comment   *string    `json:"comment"`
	IsChecked *bool      `json:"is_checked"`
	CreatedAt *time.Time `json:"created_at"`
}

type ReqComment struct {
	Id      int    `json:"id"`
	Comment string `json:"comment"`
}

type ReqCommentC3 struct {
	Id        int       `json:"id"`
	List      string    `json:"list"`
	Comment   string    `json:"comment"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GmGame struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Instructions string     `json:"instructions"`
	URL          string     `json:"url"`
	Category     string     `json:"category"`
	Tags         string     `json:"tags"`
	Thumb        string     `json:"thumb"`
	IsConstruct  string     `json:"is_construct"`
	List         *string    `json:"list"`
	Comment      *string    `json:"comment"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
