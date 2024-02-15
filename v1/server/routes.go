package server

import (
	// "io/fs"
	"os"
	filesystem "github.com/gofiber/fiber/v3/middleware/filesystem"
)

// https://github.com/gofiber/fiber/blob/main/middleware/filesystem/utils.go#L46
func ( s *Server ) SetupRoutes() {
	s.FiberApp.Use( "/" , s.PublicLimiter , filesystem.New( filesystem.Config{
		Root: os.DirFS( s.Config.ServeDirectory ) ,
		Browse: s.Config.ServeBrowsable ,
		Index: s.Config.ServeIndexFile ,
	}))
}