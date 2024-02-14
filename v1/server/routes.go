package server

import (
	// fiber "github.com/gofiber/fiber/v2"
)

func ( s *Server ) SetupRoutes() {
	s.FiberApp.Get( "/" , public_limiter , s.Home )

	// Auth
	s.FiberApp.Get( "/login" , public_limiter , s.LoginGet )
	s.FiberApp.Post( "/login" , public_limiter , s.LoginPost )
	s.FiberApp.Get( "/logout" , public_limiter , s.Logout )

	// Posts
	s.FiberApp.Post( "/post" , private_limiter , validate_session_mw , s.Post )
	// s.FiberApp.Get( "/post/:seq_id" , private_limiter , validate_session_mw , s.PostGetViaSeqID )
	s.FiberApp.Get( "/post/get/all" , private_limiter , validate_session_mw , s.PostGetAll )
	s.FiberApp.Get( "/post/get/:ulid" , private_limiter , validate_session_mw , s.PostGetViaULID )
	// s.FiberApp.Get( "/post/get/previous/unix/:start/:total" , private_limiter , validate_session_mw , s.PostGetPreviousViaUNIX )
	// s.FiberApp.Get( "/post/get/previous/ulid/:start/:total" , private_limiter , validate_session_mw , s.PostGetPreviousViaULID )
	// s.FiberApp.Get( "/post/get/after/:total/:start" , private_limiter , validate_session_mw , s.PostGetRangeViaUNIX )
	s.FiberApp.Get( "/post/get/range/unix/:start/:stop" , private_limiter , validate_session_mw , s.PostGetRangeViaUNIX )
	s.FiberApp.Get( "/post/get/range/ulid/:start/:stop" , private_limiter , validate_session_mw , s.PostGetRangeViaULID )
	// s.FiberApp.Get( "/post/:uuid" , private_limiter , validate_session_mw , s.PostGetViaUUID )
	// s.FiberApp.Get( "/post/delete/:ulid" , private_limiter , validate_session_mw , s.PostDeleteViaULID )

	// Uploads
	s.FiberApp.Post( "/upload" , private_limiter , validate_session_mw , s.Upload )

	// Calendar
	//

	// Static / Custom Pages
	s.FiberApp.Get( "/page/add/wysiwyg" , private_limiter , validate_session_mw , s.PageAddGetWYSIWYG )
	s.FiberApp.Post( "/page/add" , private_limiter , validate_session_mw , s.PageAddPost )
	s.FiberApp.Get( "/page/get" , private_limiter , validate_session_mw , s.PageGet )
	// s.FiberApp.Get( "/page/delete/:uuid" , private_limiter , validate_session_mw , s.PageDelete )
	// s.FiberApp.Get( "/page/get/:uuid" , private_limiter , validate_session_mw , s.PageGet )
	// Serving Dynamic Static Routes and Pages
	// we are doing this LAST in a catch all.
	// other ways like using only 1 handler and then the look up
	s.FiberApp.Get( "/*" , public_limiter , s.PageHandler )

}