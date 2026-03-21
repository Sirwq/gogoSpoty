package helpers

import (
	"log"
	"log/slog"
)

func CheckErrFatal(ok bool, msg string) {
	if !ok {
		log.Fatal(msg, "\nRead manual")
	}
}

func LogErr(msg string) {
	slog.Error(msg)
}

func Log(msg string) {
	log.Println(msg)
}
