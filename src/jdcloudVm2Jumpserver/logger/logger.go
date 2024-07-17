package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// SetupLogger 设置日志记录器
func SetupLogger() (*os.File, error) {
	// 确保日志目录存在
	logDir := "log"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// 创建当天的日志文件
	logFileName := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// 设置日志输出到文件和控制台
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	// 删除旧日志文件
	go func() {
		for {
			deleteOldLogs(logDir)
			time.Sleep(24 * time.Hour)
		}
	}()

	return logFile, nil
}

// deleteOldLogs 删除超过3天的日志文件
func deleteOldLogs(logDir string) {
	files, err := os.ReadDir(logDir)
	if err != nil {
		log.Printf("failed to read log directory: %v", err)
		return
	}

	cutoff := time.Now().AddDate(0, 0, -3)
	for _, file := range files {
		if !file.IsDir() {
			fileInfo, err := file.Info()
			if err != nil {
				log.Printf("failed to get file info: %v", err)
				continue
			}

			if fileInfo.ModTime().Before(cutoff) {
				err := os.Remove(filepath.Join(logDir, file.Name()))
				if err != nil {
					log.Printf("failed to delete old log file: %v", err)
				} else {
					log.Printf("deleted old log file: %s", file.Name())
				}
			}
		}
	}
}
