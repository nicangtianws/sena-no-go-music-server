package mylog

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var LOG *zerolog.Logger

func InitLog(basedir *string) {
	// 初始化日志目录写入器
	logDir := path.Join(*basedir, "logs")
	_, err := os.Stat(logDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(logDir, 0755)
		if err != nil {
			panic("Error create log dir: " + logDir)
		}
	}
	logWriter := NewDateDirWriter(logDir)
	fmt.Printf("Log dir in: %s", logDir)

	// 创建 zerolog 实例
	logger := zerolog.New(logWriter).
		With().
		Timestamp().
		Logger()

	LOG = &logger
}

// DateDirWriter 实现按日期分割目录的 io.Writer
type DateDirWriter struct {
	baseDir     string
	currentDate string
	file        *os.File
	mu          sync.Mutex
}

func NewDateDirWriter(baseDir string) *DateDirWriter {
	return &DateDirWriter{
		baseDir: baseDir,
	}
}

func (w *DateDirWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 获取当前日期（按天分割）
	currentDate := time.Now().Format("2006-01-02")

	// 日期变化时创建新文件
	if currentDate != w.currentDate {
		if w.file != nil {
			w.file.Close()
		}

		// 创建日期目录
		dateDir := filepath.Join(w.baseDir, currentDate)
		if err := os.MkdirAll(dateDir, 0755); err != nil {
			return 0, fmt.Errorf("failed to create log directory: %v", err)
		}

		// 创建日志文件
		logFile := filepath.Join(dateDir, "app.log")
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return 0, fmt.Errorf("failed to open log file: %v", err)
		}

		w.file = f
		w.currentDate = currentDate
	}

	return w.file.Write(p)
}
