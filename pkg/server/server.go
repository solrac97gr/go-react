package server

import (
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	uiFS fs.FS
}

func NewServer(ui fs.FS) *Server {
	return &Server{
		uiFS: ui,
	}
}

func (s *Server) StartServer(port string) error {
	app := fiber.New()
	app.Get("/", s.handleStatic)

	if err := app.Listen(port); err != nil {
		return err
	}
	return nil
}

func (s *Server) handleStatic(c *fiber.Ctx) error {
	if c.Method() != http.MethodGet {
		return c.SendStatus(http.StatusMethodNotAllowed)
	}

	path := c.Path()
	if path == "/" { // Add other paths that you route on the UI side here
		path = "index.html"
	}
	path = strings.TrimPrefix(path, "/")

	file, err := s.uiFS.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))
	c.Set("Content-Type", contentType)
	if strings.HasPrefix(path, "static/") {
		c.Set("Cache-Control", "public, max-age=31536000")
	}
	stat, err := file.Stat()
	if err == nil && stat.Size() > 0 {
		c.Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	}

	_, err = io.Copy(c.Response().BodyWriter(), file)
	return err
}
