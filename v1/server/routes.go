package server

import (
	// "io/fs"
	"os"
	filesystem "github.com/gofiber/fiber/v3/middleware/filesystem"
)

func ( s *Server ) SetupRoutes() {
	s.FiberApp.Use( "/" , filesystem.New( filesystem.Config{
        Root: os.DirFS( s.Config.ServeDirectory ) ,
        Browse: s.Config.ServeBrowsable ,
        Index: s.Config.ServeIndexFile ,
	}))
}