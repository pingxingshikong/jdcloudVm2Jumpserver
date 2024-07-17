package main

import (
	"jdcloudVm2Jumpserver/src/jdcloudVm2Jumpserver"
	"jdcloudVm2Jumpserver/src/jdcloudVm2Jumpserver/logger"
	"log"
)

func main() {
	// 设置日志记录器
	logFile, err := logger.SetupLogger()
	if err != nil {
		log.Fatalf("Error setting up logger: %v", err)
	}
	defer logFile.Close()

	jdcloudVm2Jumpserver.Init()
}
