package server

import (
	"fmt"
	"time"
	"context"
	fiber "github.com/gofiber/fiber/v3"
	fiber_limiter "github.com/gofiber/fiber/v3/middleware/limiter"
)

// https://github.com/gofiber/fiber/blob/main/middleware/filesystem/utils.go#L46
func ( s *Server ) SetupLimiter() {

	s.FiberApp.Use( func( c fiber.Ctx ) ( error ) {
		var ctx = context.Background()
		var ip = c.IP()
		is_banned , _ := s.REDIS.SIsMember( ctx , "b_ips" , ip ).Result()
		if is_banned == true {
			log.Warning( fmt.Sprintf( "%s === banned" , ip ) )
			c.Set( "Content-Type" , "text/html" )
			return c.SendString( "<html><h1>why ...</h1></html>" )
		}
		log.Debug( fmt.Sprintf( "%s === not banned" , ip ) )
		is_limited , _ := s.REDIS.Exists( ctx , fmt.Sprintf( "l_ips.%s" , ip ) ).Result()
		if is_limited == 1 {
			log.Warning( fmt.Sprintf( "%s === limited" , ip ) )
			c.Set( "Content-Type" , "text/html" )
			html_string := fmt.Sprintf( "<html><h1>loading ...</h1><script>setTimeout(function(){window.location.reload(1);},%d);</script></html>" , ( ( s.Config.PublicLimiterSeconds + 1 ) * 1000 ) )
			return c.SendString( html_string )
		}
		log.Debug( fmt.Sprintf( "%s === not limited" , ip ) )
		return c.Next()
	})

	s.PublicLimiter = fiber_limiter.New( fiber_limiter.Config{
		Max: s.Config.PublicLimiterMax ,
		Expiration: ( time.Duration( s.Config.PublicLimiterSeconds) * time.Second ) ,
		KeyGenerator: func( c fiber.Ctx ) string {
			return c.Get( "x-forwarded-for" )
		} ,
		LimitReached: func( c fiber.Ctx ) error {
			ip_address := c.IP()
			log_message := fmt.Sprintf( "%s === %s === %s === PUBLIC RATE LIMIT REACHED !!!" , ip_address , c.Method() , c.Path() );
			log.Info( log_message )
			var ctx = context.Background()
			var key = fmt.Sprintf( "l_ips.%s" , ip_address )
			var key_count = fmt.Sprintf( "l_ips.count.%s" , ip_address )
			current_count , _ := s.REDIS.Get( ctx , key_count ).Int64()
			if current_count > s.Config.PublicLimiterMaxLimitCount {
				s.REDIS.SAdd( ctx , "b_ips" , ip_address )
				s.REDIS.Del( ctx , key )
				s.REDIS.Del( ctx , key_count )
				log.Warning( fmt.Sprintf( "%s === banned" , ip_address ) )
				c.Set( "Content-Type" , "text/html" )
				return c.SendString( "<html><h1>why ...</h1></html>" )
			}
			s.REDIS.Set( ctx , key , 1 , ( time.Duration( s.Config.PublicLimiterSeconds) * time.Second ) )
			s.REDIS.Incr( ctx , key_count )
			c.Set( "Content-Type" , "text/html" )
			html_string := fmt.Sprintf( "<html><h1>loading ...</h1><script>setTimeout(function(){window.location.reload(1);},%d);</script></html>" , ( ( s.Config.PublicLimiterSeconds + 1 ) * 1000 ) )
			return c.SendString( html_string )
		} ,
		LimiterMiddleware: fiber_limiter.SlidingWindow{} ,
	})

}