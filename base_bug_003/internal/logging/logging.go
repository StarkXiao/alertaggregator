package logging

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Logger struct{ out *log.Logger }

func New() *Logger { return &Logger{log.New(os.Stdout, "", 0)} }
func (l *Logger) Info(message string) {
	b, _ := json.Marshal(map[string]any{"time": time.Now(), "level": "info", "message": message})
	l.out.Print(string(b))
}
