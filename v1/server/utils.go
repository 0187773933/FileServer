package server

import (
	"fmt"
	"time"
	bolt "github.com/boltdb/bolt"
	fiber "github.com/gofiber/fiber/v3"
	rate_limiter "github.com/gofiber/fiber/v3/middleware/limiter"
	utils "github.com/0187773933/FileServer/v1/utils"
)

// weak attempt at sanitizing form input to build a "username"
func SanitizeUsername( first_name string , last_name string ) ( username string ) {
	if first_name == "" { first_name = "Not Provided" }
	if last_name == "" { last_name = "Not Provided" }
	sanitized_first_name := utils.SanitizeInputString( first_name )
	sanitized_last_name := utils.SanitizeInputString( last_name )
	username = fmt.Sprintf( "%s-%s" , sanitized_first_name , sanitized_last_name )
	return
}

func ServeLoginPage( context fiber.Ctx ) ( error ) {
	return context.SendFile( "./v1/server/html/login.html" )
}

// func ServeAuthenticatedPage( context fiber.Ctx ) ( error ) {
// 	if validate_admin_cookie( context ) == false { return serve_failed_attempt( context ) }
// 	x_path := context.Route().Path
// 	url_key := strings.Split( x_path , "/admin" )
// 	if len( url_key ) < 2 { return context.SendFile( "./v1/server/html/admin_login.html" ) }
// 	// fmt.Println( "Sending -->" , url_key[ 1 ] , x_path )
// 	return context.SendFile( ui_html_pages[ url_key[ 1 ] ] )
// }

var public_limiter = rate_limiter.New( rate_limiter.Config{
	Max: 1 ,
	Expiration: 1 * time.Second ,
	KeyGenerator: func( c fiber.Ctx ) string {
		return c.Get( "x-forwarded-for" )
	} ,
	LimitReached: func( c fiber.Ctx ) error {
		ip_address := c.IP()
		log_message := fmt.Sprintf( "%s === %s === %s === PUBLIC RATE LIMIT REACHED !!!" , ip_address , c.Method() , c.Path() );
		fmt.Println( log_message )
		c.Set( "Content-Type" , "text/html" )
		return c.SendString( "<html><h1>loading ...</h1><script>setTimeout(function(){ window.location.reload(1); }, 6);</script></html>" )
	} ,
})

var private_limiter = rate_limiter.New( rate_limiter.Config{
	Max: 3 ,
	Expiration: 1 * time.Second ,
	KeyGenerator: func( c fiber.Ctx ) string {
		return c.Get( "x-forwarded-for" )
	} ,
	LimitReached: func( c fiber.Ctx ) error {
		ip_address := c.IP()
		log_message := fmt.Sprintf( "%s === %s === %s === PUBLIC RATE LIMIT REACHED !!!" , ip_address , c.Method() , c.Path() );
		fmt.Println( log_message )
		c.Set( "Content-Type" , "text/html" )
		return c.SendString( "<html><h1>loading ...</h1><script>setTimeout(function(){ window.location.reload(1); }, 6);</script></html>" )
	} ,
})

func ( s *Server ) Set( bucket_name string , key string , value string ) {
	s.DB.Update( func( tx *bolt.Tx ) error {
		b , err := tx.CreateBucketIfNotExists( []byte( bucket_name ) )
		if err != nil { log.Debug( err ); return nil }
		err = b.Put( []byte( key ) , []byte( value ) )
		if err != nil { log.Debug( err ); return nil }
		return nil
	})
	return
}

func ( s *Server ) Get( bucket_name string , key string ) ( result string ) {
	s.DB.View( func( tx *bolt.Tx ) error {
		b := tx.Bucket( []byte( bucket_name ) )
		if b == nil { return nil }
		v := b.Get( []byte( key ) )
		if v == nil { return nil }
		result = string( v )
		return nil
	})
	return
}