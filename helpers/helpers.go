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
