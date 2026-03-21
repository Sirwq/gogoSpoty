package helpers

import "log"

func CheckErr(ok bool, msg string) {
	if !ok {
		log.Fatal(msg, "\nRead manual")
	}
}
