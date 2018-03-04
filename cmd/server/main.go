package main

import (
	"log"

	"github.com/sir-wiggles/arc/pkg/webstore"
	"github.com/sir-wiggles/arc/pkg/webstore/mock"
)

func main() {

	var mus *mock.UserService
	var us *webstore.UserService
	log.Println("from api", mus, us)
}
