package server

import (
	"fmt"
	json "encoding/json"
	fiber "github.com/gofiber/fiber/v2"
	types "github.com/0187773933/FileServer/v1/types"
	// bolt "github.com/boltdb/bolt"
)

func ( s *Server ) Home( context *fiber.Ctx ) ( error ) {
	session := validate_session( context )
	if session == false {
		return context.JSON( fiber.Map{
			"route": "/" ,
			"source": "https://github.com/0187773933/Blogger" ,
		})
	}
	log.Debug( "Logged In User , Sending Home Page" )
	return context.SendFile( "./v1/server/html/home.html" )
}

// Adds { key: HTML_STRING-b64 } to static-routes
func ( s *Server ) PageAddPost( context *fiber.Ctx ) ( error ) {
	log.Debug( "PageAddPost()" )
	context_body := context.Body()
	var p types.Page
	json.Unmarshal( context_body , &p )
	// p.UUID = uuid.NewV4().String()
	log.Debug( fmt.Sprintf( "Storing Content for URL : %s" , p.URL ) )
	s.Set( "pages" , p.URL ,  p.HTMLB64 )
	fmt.Println( p.HTMLB64 )
	return context.JSON( fiber.Map{
		"route": "/page/add" ,
		// "uuid": p.UUID ,
		"url": p.URL ,
		"result": true ,
	})
}

func ( s *Server ) PageAddGetWYSIWYG( context *fiber.Ctx ) ( error ) {
	// log.Debug( "PageAddGet()" )
	return context.SendFile( "./v1/server/html/page_add_wysiwyg.html" )
}

func ( s *Server ) PageGet( context *fiber.Ctx ) ( error ) {
	log.Debug( "PageGet()" )
	// x_url := context.Params( "url" )
	x_url := context.Query( "url" )
	page_html_b64 := s.Get( "pages" , x_url )
	return context.JSON( fiber.Map{
		"route": "/page/get/:url" ,
		"url": x_url ,
		"html_b64": page_html_b64 ,
		"result": true ,
	})
}

// so like https://github.com/quilljs/quill
// then we save the quill.root.innerHTML ? idk
// we have to do 2 lookups then , theres no way.
// first lookup == confirm static page exists , render html parent
// second lookup == GET JSON request sent by html parent for actual content
// or we could just render parent template as catch all , and then do 1 lookup for if any content exists ?
func ( s *Server ) PageHandler( context *fiber.Ctx ) ( error ) {
	sent_path := context.Path()
	log.Debug( fmt.Sprintf( "PageHandler( %s )" , sent_path ) )
	// sent_queries := context.Queries()
	// page_html := s.Get( "pages" , sent_path )
	// if page_html == "" {
	// 	return context.JSON( fiber.Map{
	// 		"route": "/*" ,
	// 		"sent_path": sent_path ,
	// 		"sent_queries": sent_queries ,
	// 		"page_html": page_html ,
	// 		"result": false ,
	// 	})
	// }
	return context.SendFile( "./v1/server/html/page.html" )
}