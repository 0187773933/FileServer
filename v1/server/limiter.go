package server

import (
	"fmt"
	"time"
	"net"
	"context"
	"strings"
	_ "embed"
	// "strconv"
	fiber "github.com/gofiber/fiber/v3"
	fiber_limiter "github.com/gofiber/fiber/v3/middleware/limiter"
	// fiber_redis "github.com/gofiber/storage/redis/v3"
	// fiberpow "github.com/witer33/fiberpow"
)

//go:embed vps/ipv4/1.txt
var VPS_IPV4_1 string // custom

//go:embed vps/ipv4/2.txt
var VPS_IPV4_2 string // https://github.com/lord-alfred/ipranges/blob/main/all/ipv4_merged.txt

//go:embed vps/ipv4/3.txt
var VPS_IPV4_3 string // https://github.com/jhassine/server-ip-addresses/blob/master/data/datacenters.txt

//go:embed vps/ipv4/4.txt
var VPS_IPV4_4 string // https://github.com/blacklanternsecurity/cloudcheck/blob/master/cloud_providers.json
// jq -r '[.[] | .cidrs? // empty] | add | .[]' cloud_providers.json > 4.txt

//go:embed vps/ipv6/1.txt
var VPS_IPV6_1 string // custom

var VPS_IPV_CIDRS = []string{}
var VPS_IPV_NETWORKS = []*net.IPNet{}

// store in bolt instead of embed ?
func prep_ipv_files() {
	// ipv4
	lines := strings.Split( VPS_IPV4_1 , "\n" )
	lines = append( lines , strings.Split( VPS_IPV4_2 , "\n" )... )
	lines = append( lines , strings.Split( VPS_IPV4_3 , "\n" )... )
	lines = append( lines , strings.Split( VPS_IPV4_4 , "\n" )... )
	// ipv6
	lines = append( lines , strings.Split( VPS_IPV6_1 , "\n" )... )
	for _ , line := range lines {
		trimmed_line := strings.TrimSpace( line )
		_ , network , err := net.ParseCIDR( trimmed_line )
		if err != nil {
			continue
		}
		VPS_IPV_CIDRS = append( VPS_IPV_CIDRS , trimmed_line )
		VPS_IPV_NETWORKS = append( VPS_IPV_NETWORKS , network )
	}
}

func is_ip_a_vps( ip_address string ) ( result bool ) {
	result = false
	x_ip := net.ParseIP( ip_address )
	for _ , network := range VPS_IPV_NETWORKS {
		if network.Contains( x_ip ) {
			result = true
			return
		}
	}
	return
}

// const (
// 	SlidingWindowAlgorithm = go_limiter.SlidingWindowAlgorithm
// 	GCRAAlgorithm          = go_limiter.GCRAAlgorithm
// 	DefaultKeyPrefix       = "fiber_limiter"
// )

// https://github.com/gofiber/fiber/blob/main/middleware/filesystem/utils.go#L46
// https://github.com/Shareed2k/go_limiter/blob/master/gcra_lua.go
func ( s *Server ) SetupLimiter() {

	// Global Banned IPs
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

	// TODO = Basic UserAgent Test
	// https://github.com/stephenafamo/isbot/blob/main/crawler-user-agents.json
	// s.FiberApp.Use( func( c fiber.Ctx ) ( error ) {

	// }

	// TODO = VPS IP Address Range
	// https://github.com/lord-alfred/ipranges
	// https://github.com/jhassine/server-ip-addresses
	// https://github.com/blacklanternsecurity/cloudcheck
	// https://github.com/SecOps-Institute/Tor-IP-Addresses
	// https://github.com/borestad/0blocklist-abuseipdb
	prep_ipv_files()
	s.FiberApp.Use( func( c fiber.Ctx ) ( error ) {
		x_ip := c.IP()
		// x_ip := "216.244.66.246"
		is_vps := is_ip_a_vps( x_ip )
		if is_vps == true {
			log.Debug( fmt.Sprintf( "%s === vps" , x_ip ) )
			c.Set( "Content-Type" , "text/html" )
			return c.SendString( "<html><h1>why ...</h1></html>" )
		}
		log.Debug( fmt.Sprintf( "%s === not vps" , x_ip ) )
		return c.Next()
	})

	// TODO = Crawler Behavior
	// s.FiberApp.Use( func( c fiber.Ctx ) ( error ) {

	// }

	// // Javascript Sha256 Challenge
	// port_number , _ := strconv.Atoi( s.Config.RedisPort )
	// js_challenge_store := fiber_redis.New(fiber_redis.Config{
	// 	Host: s.Config.RedisHost ,
	// 	Port: port_number ,
	// 	Password: s.Config.RedisPassword ,
	// 	Database: s.Config.RedisDBNumber ,
	// 	Reset: false ,
	// })
	// s.FiberApp.Use( s.NewJSChallenge( JSChallengeConfig{
	// 	PowInterval: ( 10 * time.Minute ) ,
	// 	Difficulty: 60000 ,
	// 	// Filter: func( c fiber.Ctx ) bool {
	// 	// 	return c.IP() == "127.0.0.1"
	// 	// } ,
	// 	Storage: js_challenge_store ,
	// }))

	// Rate Limiter
	s.FiberApp.Use( fiber_limiter.New( fiber_limiter.Config{
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
	}))

}