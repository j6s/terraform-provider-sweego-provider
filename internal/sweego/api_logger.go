package sweego

import (
	"fmt"
	"log"
)

type SweegoApiLogger interface {
	Info(message string)
	Error(message string)
	Debug(message string)
}

type GolangLogger struct{}

func (l GolangLogger) Info(message string) {
	log.Println(fmt.Sprintf("[INFO] %s", message))
}
func (l GolangLogger) Error(message string) {
	log.Println(fmt.Sprintf("[ERROR] %s", message))
}
func (l GolangLogger) Debug(message string) {
	log.Println(fmt.Sprintf("[DEBUG] %s", message))
}
