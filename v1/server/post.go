package server

import (
	"fmt"
	"bytes"
	// "strconv"
	"encoding/json"
	ulid "github.com/oklog/ulid/v2"
	uuid "github.com/satori/go.uuid"
	bolt "github.com/boltdb/bolt"
	types "github.com/0187773933/FileServer/v1/types"
	fiber "github.com/gofiber/fiber/v2"
	utils "github.com/0187773933/FileServer/v1/utils"
)

func ( s *Server ) PostGetRange( min string , max string ) ( result []types.Post ) {
	s.DB.View( func( tx *bolt.Tx ) error {
		c := tx.Bucket( []byte( "posts" ) ).Cursor()
		for k , v := c.Seek( []byte( min ) ); k != nil && bytes.Compare( k , []byte( max ) ) <= 0; k , v = c.Next() {
			var p types.Post
			json.Unmarshal( v , &p )
			result = append( result , p )
		}
		return nil
	})
	return
}

func ( s *Server ) PostGetAll( context *fiber.Ctx ) ( error ) {
	var posts []types.Post
	s.DB.View( func( tx *bolt.Tx ) error {
		c := tx.Bucket( []byte( "posts" ) ).Cursor()
		for k , v := c.First(); k != nil; k , v = c.Next() {
			var p types.Post
			json.Unmarshal( v , &p )
			posts = append( posts , p )
			t , _ := ulid.Parse( p.ULID )
			fmt.Println( t.Time() )
		}
		return nil
	})
	return context.JSON( fiber.Map{
		"url": "/post/get/all" ,
		"method": "GET" ,
		"posts": posts ,
		"result": true ,
	})
}

func ( s *Server ) PostGetViaULID( context *fiber.Ctx ) error {
	x_ulid := context.Params( "ulid" )
	var p types.Post
	s.DB.View( func( tx *bolt.Tx ) error {
		b := tx.Bucket( []byte( "posts" ) )
		v := b.Get( []byte( x_ulid ) )
		json.Unmarshal( v , &p )
		return nil
	})
	return context.JSON( fiber.Map{
		"url": "/post/get/:ulid" ,
		"ulid": x_ulid ,
		"method": "GET" ,
		"post": p ,
		"result": true ,
	})
}

func ( s *Server ) PostGetRangeViaUNIX( context *fiber.Ctx ) error {
	start := utils.UnixToULID( context.Params( "start" ) )
	stop := utils.UnixToULID( context.Params( "stop" ) )
	posts := s.PostGetRange( start , stop )
	return context.JSON( fiber.Map{
		"url": "/post/get/range" ,
		"method": "GET" ,
		"posts": posts ,
		"result": true ,
	})
}

func ( s *Server ) PostGetRangeViaULID( context *fiber.Ctx ) error {
	posts := s.PostGetRange( context.Params( "start" ) , context.Params( "stop" ) )
	return context.JSON( fiber.Map{
		"url": "/post/get/range" ,
		"method": "GET" ,
		"posts": posts ,
		"result": true ,
	})
}

func ( s *Server ) Post( context *fiber.Ctx ) ( error ) {
	context_body := context.Body()
	var p types.Post
	json.Unmarshal( context_body , &p )
	p.Date = utils.GetFormattedTimeString()
	p.UUID = uuid.NewV4().String()
	p.ULID = ulid.Make().String()
	post_json_bytes , _ := json.Marshal( p )
	s.DB.Update( func( tx *bolt.Tx ) error {
		posts_bucket , _ := tx.CreateBucketIfNotExists( []byte( "posts" ) )
		posts_bucket.Put( []byte( p.ULID ) , post_json_bytes )
		return nil
	})
	return context.JSON( fiber.Map{
		"url": "/post" ,
		"method": "POST" ,
		"post": p ,
		"result": true ,
	})
}

// func ( s *Server ) PostGetViaSeqID( context *fiber.Ctx ) ( error ) {
// 	// var p types.Post
// 	// json.Unmarshal( context_body , &p )
// 	seq_id := context.Params( "seq_id" )
// 	seq_id_int , _ := strconv.Atoi( seq_id )
// 	// seq_id_index_int := ( seq_id_int - 1 )
// 	// seq_id_index_string := strconv.Itoa( seq_id_index_int )
// 	// var post_string string
// 	var t_post types.Post
// 	fmt.Println( "Seq ID ===" , seq_id )
// 	s.DB.View( func( tx *bolt.Tx ) error {
// 		// c := tx.Bucket( []byte( "posts" ) ).Cursor()
// 		// k , v := c.Seek( []byte( seq_id ) )
// 		// fmt.Println( k , v )
// 		b := tx.Bucket( []byte( "posts" ) )
// 		c := b.Cursor()
// 		_ , v := c.Seek( []byte( seq_id ) )
// 		fmt.Println( string( v ) )
// 		if v == nil {
// 			_ , v := c.Prev() // you have to do it this way. there is no c.Curr(). and c.Seek always goes +1 somehow
// 			json.Unmarshal( v , &t_post )
// 		}

// 		// post_string = string( v )
// 		// fmt.Printf( "key=%s, value=%s\n" , k , v )
// 		// for k, v := c.First(); k != nil; k, v = c.Next() {
// 		// 	fmt.Printf("key=%s, value=%s\n", k, v)
// 		// }
// 		return nil
// 	})
// 	if t_post.SeqID != seq_id_int {
// 		return context.JSON( fiber.Map{
// 			"url": "/post/:seq_id" ,
// 			"method": "GET" ,
// 			"result": false ,
// 		})
// 	}
// 	return context.JSON( fiber.Map{
// 		"url": "/post/:seq_id" ,
// 		"method": "GET" ,
// 		"post": t_post ,
// 		"result": true ,
// 	})
// }

// func ( s *Server ) PostGetViaUUID( context *fiber.Ctx ) ( error ) {
// 	var p types.Post
// 	json.Unmarshal( context_body , &p )
// 	s.DB.View( func( tx *bolt.Tx ) error {

// 		c := tx.Bucket( []byte( "posts" ) ).Cursor()
// 		k , v := c.Seek( min )

// 		posts_bucket , _ := tx.CreateBucketIfNotExists( []byte( "posts" ) )
// 		post_id , _ := posts_bucket.NextSequence()
// 		posts_bucket.Put( utils.IToB( post_id ) , context_body )
// 		return nil
// 	})
// 	return context.JSON( fiber.Map{
// 		"url": "/post/:uuid" ,
// 		"method": "GET" ,
// 		"result": true ,
// 	})
// }