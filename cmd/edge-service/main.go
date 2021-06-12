package main

import "github.com/fcproto/prototype/pkg/logger"

func main() {
	log := logger.New("edge")
	log.Info("starting edge service...")
}
