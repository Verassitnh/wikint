package main

import (
	"github.com/google/uuid"
)

func main() {
	usrScrape(User{
		name: "Sam Harkness",
		id:   uuid.NewString(),
		urls: []string{"https://www.facebook.com/profile.php?id=100085268747599&sk=friends"},
		occupation: Ocupation{
			name: "Mechanic",
			id:   uuid.NewString(),
		},
	})
}
