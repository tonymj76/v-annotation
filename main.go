package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/tonymj76/video-annotator/router"
)

// init gets called before the main function
func init() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalln("No .env file found create")
	}
}

func main() {
	r, err := router.Router()
	if err != nil {
		logrus.Fatalln(err)
	}

	if err := r.Run(":9090"); err != nil {
		logrus.Fatalln(err)
	}

}
