package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"ugames/internal/handler"
)

func NewRoutes(h *handler.Handler) *fiber.App {
	app := fiber.New()
	app.Static("/", "./static")
	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		log.Info().Msg("Healthcheck is working!")
		return c.SendString("Alive!")
	})
	//app.Get("/keywords", h.GetKeyWordsList)
	app.Get("/api/repos", h.GetCheckedReposList)
	app.Post("/api/collectGitRepos", h.CollectGitRepos)
	app.Get("/api/checkRepos", h.CheckRepos)
	app.Put("/api/addComment", h.AddComment)
	app.Get("/api/dbfix", h.FixDb)
	return app
}
