package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"ugames/internal/models"

	"ugames/internal/config"
	"ugames/internal/repo"
)

type Handler struct {
	pool *repo.Repo
	cnf  *config.Cnf
}

func NewHandler(pool *repo.Repo, cnf *config.Cnf) *Handler {
	return &Handler{
		pool: pool,
		cnf:  cnf,
	}
}

func (h *Handler) GetKeyWordsList(c *fiber.Ctx) error {
	data, err := h.pool.GetKeyWordsList()
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return c.JSON(data)
}

func (h *Handler) GetCheckedReposList(c *fiber.Ctx) error {
	err := h.pool.DeallocateAll()
	if err != nil {
		log.Error().Msg(err.Error())
	}

	data, err := h.pool.GetCheckedRepos()
	if err != nil {
		log.Error().Msg(err.Error())
	}

	//go h.CheckReposFunc()
	return c.JSON(data)
}

func (h *Handler) CollectGitRepos(c *fiber.Ctx) error {
	var resp models.Resp
	var req models.KeyWordReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	githubResponse := h.GetGitRepos(req.KeyWord, 1)
	for _, item := range githubResponse.Items {
		h.pool.InsertRepo(item.FullName, item.Homepage, req.KeyWord)
	}

	pages := int(math.Ceil(float64(githubResponse.TotalCount) / 100))
	if pages > 1 {
		for i := 2; i <= pages; i++ {
			githubResponse = h.GetGitRepos(req.KeyWord, i)
			for _, item := range githubResponse.Items {
				h.pool.InsertRepo(item.FullName, item.Homepage, req.KeyWord)
			}
		}
	}

	resp.Status = "Success"
	resp.Message = "Напарсено " + strconv.Itoa(githubResponse.TotalCount) + " репозиториев, " + strconv.Itoa(pages) + " страниц API. Запущена проверка."
	resp.Data = nil

	go h.CheckReposFunc()
	return c.JSON(resp)
}

func (h *Handler) CheckRepos(c *fiber.Ctx) error {
	var resp models.Resp
	uncheckedList, err := h.pool.GetUncheckedRepos()
	if err != nil {
		log.Error().Msg(err.Error())
		resp.Status = "Error"
		resp.Message = err.Error()
		resp.Data = nil
		return c.JSON(resp)
	}

	for _, repo := range uncheckedList {
		h.CheckGetRepo(repo)
	}

	resp.Status = "Success"
	resp.Message = ""
	resp.Data = uncheckedList
	return c.JSON(resp)
}

func (h *Handler) CheckReposFunc() error {
	uncheckedList, err := h.pool.GetUncheckedRepos()
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	for _, repo := range uncheckedList {
		h.CheckGetRepo(repo)
	}

	return nil
}

func (h *Handler) GetGitRepos(keyWord string, page int) models.GitReposResp {
	encodedString := url.QueryEscape(keyWord)
	pageStr := strconv.Itoa(page)

	url := "https://api.github.com/search/repositories?q=" + encodedString + "&per_page=100&page=" + pageStr
	log.Printf(url)
	client := &http.Client{}
	reqGit, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	reqGit.Header.Set("Authorization", "Bearer "+h.cnf.GithubToken)
	respGit, err := client.Do(reqGit)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer respGit.Body.Close()

	var githubResponse models.GitReposResp
	if err := json.NewDecoder(respGit.Body).Decode(&githubResponse); err != nil {
		log.Error().Msg(err.Error())
	}

	return githubResponse
}

func (h *Handler) CheckGetRepo(repo models.Repos) error {
	var url string
	if repo.RepoName != nil {
		url = "https://api.github.com/repos/" + *repo.RepoName + "/contents/README.md"
	} else {
		return nil
	}

	client := &http.Client{}
	reqGit, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	reqGit.Header.Set("Authorization", "Bearer "+h.cnf.GithubToken)
	respGit, err := client.Do(reqGit)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer respGit.Body.Close()

	var gitReposReadme models.GitReposReadme
	if err := json.NewDecoder(respGit.Body).Decode(&gitReposReadme); err != nil {
		log.Error().Msg(err.Error())
	}

	decodedContent, err := base64.StdEncoding.DecodeString(gitReposReadme.Content)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	keywords := []string{
		"github.io",
		"itch.io",
		"youtube.com",
		"play.unity.com",
		"play.google.com",
		"gamejolt.com",
		"unityroom.com",
		"maikire.xyz",
		".gif",
		".png",
		".jpg",
		"simmer.io",
		"pubnub.com",
		"kongregate.com",
		"web.app",
		"sharemygame.com",
	}

	// Приводим текст к нижнему регистру для сравнения
	lowerText := strings.ToLower(string(decodedContent))
	found := []string{}

	for _, keyword := range keywords {
		if strings.Contains(lowerText, strings.ToLower(keyword)) {
			found = append(found, keyword)
		}
	}

	var content string
	if len(found) > 0 {
		for _, match := range found {
			content += match + " "
			//log.Printf(*repo.Content)
		}
	} else {
		log.Printf("Слова не найдены.")
	}

	repo.Content = &content

	err = h.pool.UpdateCheckedRepo(repo)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	return nil
}

func (h *Handler) AddComment(c *fiber.Ctx) error {
	var resp models.Resp
	var req models.ReqComment
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := h.pool.AddComment(req)
	if err != nil {
		log.Error().Msg(err.Error())
		resp.Status = "Error"
		resp.Message = "Ошибка при добавлении комментария!"
		return c.JSON(resp)
	}

	//go h.pool.DeallocateAll()

	resp.Status = "Success"
	resp.Message = "Комментарий добавлен успешно!"
	return c.JSON(resp)
}

func (h *Handler) FixDb(c *fiber.Ctx) error {
	err := h.pool.DeallocateAll()
	if err != nil {
		log.Error().Msg(err.Error())
	}

	var resp models.Resp
	resp.Status = "Success"
	resp.Message = "Выполнена волшебная SQL команда!"
	return c.JSON(resp)
}

func (h *Handler) CollectGMGames(c *fiber.Ctx) error {
	// URL API
	baseURL := "https://gamemonetize.com/feed.php?format=0&num=50&page="

	// Цикл по страницам от 1 до 666
	for page := 1; page <= 666; page++ {
		// Формируем URL для текущей страницы
		apiURL := fmt.Sprintf("%s%d", baseURL, page)
		// Запрос к API
		resp, err := http.Get(apiURL)
		if err != nil {
			fmt.Println("Ошибка запроса:", err)
		}
		defer resp.Body.Close()

		// Чтение ответа
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Ошибка чтения ответа:", err)
		}

		// Распарсить XML
		var games []models.GmGame
		err = json.Unmarshal(body, &games)
		if err != nil {
			fmt.Println("Ошибка парсинга json:", err)
		}

		// Вывод данных
		for _, game := range games {
			fmt.Printf("ID: %s\nTitle: %s\nDescription: %s\nURL: %s\n\n", game.ID, game.Title, game.Description, game.URL)
			h.pool.InsertGmGame(game)
		}
	}
	return nil
}

func (h *Handler) FindConstructGame(c *fiber.Ctx) error {
	for i := 0; i < 40; i++ {
		games, _ := h.pool.GetGmGames()
		log.Info().Msg("Получено " + strconv.Itoa(len(games)) + " игр для проверки.")

		for _, game := range games {
			url := game.URL

			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("Ошибка при запросе URL:", err)
				return nil
			}
			defer resp.Body.Close()

			// Проверяем статус ответа
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Ошибка: статус ответа %d\n", resp.StatusCode)
				return nil
			}

			// Читаем содержимое страницы
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Ошибка при чтении содержимого страницы:", err)
				return nil
			}

			// Преобразуем содержимое в строку
			pageContent := string(body)

			// Проверяем наличие строки "Construct" (игнорируя регистр)
			if strings.Contains(strings.ToLower(pageContent), "construct") {
				game.IsConstruct = "Y"
				h.pool.UpdateGmGame(game)
				log.Info().Msg(url)
			} else {
				game.IsConstruct = "N"
				h.pool.UpdateGmGame(game)
			}
		}
	}
	return nil
}

func (h *Handler) GetConstructGamesList(c *fiber.Ctx) error {
	filter := c.Params("filter")
	data, err := h.pool.GetConstructGames(filter)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return c.JSON(data)
}

func (h *Handler) AddCommentC3(c *fiber.Ctx) error {
	var resp models.Resp
	var req models.ReqCommentC3
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := h.pool.AddCommentC3(req)
	if err != nil {
		log.Error().Msg(err.Error())
		resp.Status = "Error"
		resp.Message = "Ошибка при добавлении комментария!"
		return c.JSON(resp)
	}

	//go h.pool.DeallocateAll()

	resp.Status = "Success"
	resp.Message = "Комментарий добавлен успешно!"
	return c.JSON(resp)
}

//https://supabase.com/dashboard/project/omdkwxnqhidnjisdfboh/settings/database
//https://docs.github.com/ru/rest/search/search?apiVersion=2022-11-28
//https://api.github.com/search/repositories?q=unity+game&per_page=100&page=2
//https://api.github.com/repos/connorwright1122/piratesoftware-jam-2024/contents/README.md
