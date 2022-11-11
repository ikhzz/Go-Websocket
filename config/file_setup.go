package config

import (
	"os"
	"io"
	"log"
)

func InitFileSetup() io.Writer {
	//setup files and directories
	if _, err := os.Stat("./log"); os.IsNotExist(err) {
		os.Mkdir("./log", 0777)
	}
	if _, err := os.Stat("./files"); os.IsNotExist(err) {
		os.Mkdir("./files", 0777)
	}
	if _, err := os.Stat("./files/templates"); os.IsNotExist(err) {
		os.Mkdir("./files/templates", 0777)
	}
	if _, err := os.Stat("./files/assets"); os.IsNotExist(err) {
		os.Mkdir("./files/assets", 0777)
	}

	f, _ := os.Create("./log/service.log")

	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(f)
	
	return mw
}