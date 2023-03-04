package middleware

import "time"

func ServerLoop() {
	for {
		time.Sleep(time.Hour)
	}
}
