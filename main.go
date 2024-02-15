package main

import (
	"fmt"
	"os"
	"time"
	"os/signal"
	"context"
	"syscall"
	"path/filepath"
	types "github.com/0187773933/FileServer/v1/types"
	utils "github.com/0187773933/FileServer/v1/utils"
	server "github.com/0187773933/FileServer/v1/server"
	bolt "github.com/boltdb/bolt"
	redis "github.com/redis/go-redis/v9"
	logger "github.com/0187773933/FileServer/v1/logger"
)

var s server.Server
var DB *bolt.DB
var REDIS *redis.Client

func SetupCloseHandler() {
	c := make( chan os.Signal )
	signal.Notify( c , os.Interrupt , syscall.SIGTERM , syscall.SIGINT )
	go func() {
		<-c
		fmt.Println( "\r- Ctrl+C pressed in Terminal" )
		fmt.Println( "Shutting Down File Server" )
		DB.Close()
		REDIS.Close()
		s.FiberApp.Shutdown()
		os.Exit( 0 )
	}()
}

func SetupRedis( config *types.ConfigFile ) {
	REDIS = redis.NewClient( &redis.Options{
		Addr: config.RedisAddress ,
		Password: config.RedisPassword ,
		DB: config.RedisDBNumber ,
	})
	var ctx = context.Background()
	ping_result , err := REDIS.Ping( ctx ).Result()
	logger.Log.Printf( "REDIS Connected : PING = %s" , ping_result )
	if err != nil { panic( err ) }
}

func main() {

	config_file_path := "./config.yaml"
	if len( os.Args ) > 1 { config_file_path , _ = filepath.Abs( os.Args[ 1 ] ) }
	config := utils.ParseConfig( config_file_path )
	// fmt.Printf( "Loaded Config File From : %s\n" , config_file_path )

	logger.Init()
	logger.Log.Printf( "Loaded Config File From : %s" , config_file_path )

	DB , _ = bolt.Open( config.BoltDBPath , 0600 , &bolt.Options{ Timeout: ( 3 * time.Second ) } )
	SetupRedis( &config )

	SetupCloseHandler()

	// utils.GenerateNewKeys()
	s = server.New( DB , REDIS , config )
	s.Start()

}
