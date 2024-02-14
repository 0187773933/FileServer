package server

import (
	"fmt"
	"io/ioutil"
	// uuid "github.com/satori/go.uuid"
	fiber "github.com/gofiber/fiber/v2"
	types "github.com/0187773933/FileServer/v1/types"
)

func ( s *Server ) Upload( context *fiber.Ctx ) ( error ) {
	form , err := context.MultipartForm()
	if err != nil {
		return context.Status( fiber.StatusBadRequest ).SendString( "Error Parsing Form" )
	}
	var file_data types.FileData
	bytes , exists := form.File[ "bytes" ]
	if exists && len( bytes ) > 0 {
		file , err := bytes[ 0 ].Open()
		defer file.Close()
		if err == nil {
			data , err := ioutil.ReadAll( file )
			if err == nil {
				file_data.FileName = bytes[ 0 ].Filename
				file_data.Data = data
			}
		}
	}
	if file_data.Data != nil {
		fmt.Println( "Uploaded Bytes from :" , file_data.FileName )
		fmt.Println( file_data.Data )
	}
	return context.JSON( fiber.Map{
		"url": "/upload" ,
		"method": "POST" ,
		"file": file_data ,
		"result": true ,
	})
}