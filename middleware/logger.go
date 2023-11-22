package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/natefinch/lumberjack"
)

var fileLogger *log.Logger

func init() {
	logFile := &lumberjack.Logger{
		Filename:   "E:\\Telechargement\\estiam-main\\logfile.txt",
		MaxSize:    5, 
		MaxBackups: 3,
		MaxAge:     28, 
	}

	fileLogger = log.New(logFile, "", log.LstdFlags)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(startTime)

		log.Printf("[%s] [%s] %s %s %s\n",
			r.RemoteAddr, // Adresse IP du client
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			duration.String(), 
		)

		
		fileLogger.Printf("[%s] [%s] %s %s %s\n",
			r.RemoteAddr, 
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			duration.String(), 
		)
	})
}
