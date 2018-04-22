package main

import (
	"fmt"
	"go-dropbox/service"

	"github.com/labstack/gommon/log"
)

func main() {
	dropbox := godropbox.NewDropbox()

	//get user information
	log.Info("get user information")
	if user, err := dropbox.User().Get(); err != nil {
		log.Error(err.Error())
	} else {
		fmt.Printf("\n\nUSER: %+v \n\n", user)
	}

	// upload a file
	log.Info("upload a file")
	if response, err := dropbox.File().Upload("/teste.txt", []byte("teste")); err != nil {
		log.Error(err.Error())
	} else {
		fmt.Printf("\n\nUPLOADED: %+v \n\n", response)
	}

	// download the uploaded file
	log.Info("download the uploaded file")
	if response, err := dropbox.File().Download("/teste.txt"); err != nil {
		log.Error(err.Error())
	} else {
		fmt.Printf("\n\nDOWNLOADED: %s \n\n", string(response))
	}

	// create folder
	log.Info("listing folder")
	if response, err := dropbox.Folder().Create("/bananas"); err != nil {
		log.Error(err.Error())
	} else {
		fmt.Printf("\n\nCREATED FOLDER: %+v \n\n", response)
	}

	// listing folder
	log.Info("listing folder")
	if response, err := dropbox.Folder().List("/"); err != nil {
		log.Error(err.Error())
	} else {
		fmt.Printf("\n\nLIST FOLDER: %+v \n\n", response)
	}

	// deleting the uploaded file
	log.Info("deleting the uploaded file")
	if response, err := dropbox.File().Delete("/teste.txt"); err != nil {
		log.Error(err.Error())
	} else {
		fmt.Printf("\n\nDELETED FILE: %+v \n\n", response)
	}

	// deleting the created folder
	log.Info("deleting the created folder")
	if response, err := dropbox.Folder().DeleteFolder("/bananas"); err != nil {
		log.Error(err.Error())
	} else {
		fmt.Printf("\n\nDELETED FOLDER: %+v \n\n", response)
	}
}
