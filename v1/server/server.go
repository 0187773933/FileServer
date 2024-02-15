package server

import (
	"fmt"
	bolt "github.com/boltdb/bolt"
	redis "github.com/redis/go-redis/v9"
	fiber "github.com/gofiber/fiber/v3"
	fiber_cookie "github.com/gofiber/fiber/v3/middleware/encryptcookie"
	fiber_cors "github.com/gofiber/fiber/v3/middleware/cors"
	// favicon "github.com/gofiber/fiber/v3/middleware/favicon"
	types "github.com/0187773933/FileServer/v1/types"
	logger "github.com/0187773933/FileServer/v1/logger"
)

var GlobalConfig *types.ConfigFile
var log = logger.GetLogger()

type Server struct {
	FiberApp *fiber.App `yaml:"fiber_app"`
	DB *bolt.DB `yaml:"-"`
	REDIS *redis.Client `yaml:"-"`
	Config types.ConfigFile `yaml:"config"`
	PublicLimiter fiber.Handler `yaml:"-"`
}

func request_logging_middleware( context fiber.Ctx ) ( error ) {
	ip_address := context.Get( "x-forwarded-for" )
	if ip_address == "" { ip_address = context.IP() }
	log_message := fmt.Sprintf( "%s === %s === %s" , ip_address , context.Method() , context.Path() );
	log.Println( log_message )
	return context.Next()
}

func New( db *bolt.DB , x_redis *redis.Client , config types.ConfigFile ) ( server Server ) {
	server.FiberApp = fiber.New()
	server.DB = db
	server.REDIS = x_redis
	server.Config = config
	GlobalConfig = &config
	server.FiberApp.Use( request_logging_middleware )
	// server.FiberApp.Use( favicon.New() )
	server.FiberApp.Use( func( c fiber.Ctx ) error {
		if c.Path() == "/favicon.ico" { return c.SendStatus( fiber.StatusNoContent ) }
		return c.Next()
	})
	server.FiberApp.Use( fiber_cookie.New( fiber_cookie.Config{
		Key: server.Config.ServerCookieSecret ,
	}))
	allow_origins_string := fmt.Sprintf( "%s, %s" , server.Config.ServerBaseUrl , server.Config.ServerLiveUrl )
	server.FiberApp.Use( fiber_cors.New( fiber_cors.Config{
		AllowOrigins: allow_origins_string ,
		AllowHeaders:  "Origin, Content-Type, Accept, key" ,
		AllowCredentials: true ,
	}))
	server.SetupLimiter()
	server.SetupRoutes()
	server.FiberApp.Get( "/*" , func( context fiber.Ctx ) ( error ) { return context.Redirect().To( "/" ) } )
	return
}

func ( s *Server ) Start() {
	fmt.Println( "\n" )
	fmt.Printf( "Listening on http://localhost:%s\n" , s.Config.ServerPort )
	fmt.Printf( "Admin Login @ http://localhost:%s/admin/login\n" , s.Config.ServerPort )
	fmt.Printf( "Admin Username === %s\n" , s.Config.AdminUsername )
	fmt.Printf( "Admin Password === %s\n" , s.Config.AdminPassword )
	fmt.Printf( "Admin API Key === %s\n" , s.Config.ServerAPIKey )
	s.FiberApp.Listen( fmt.Sprintf( ":%s" , s.Config.ServerPort ) )
}