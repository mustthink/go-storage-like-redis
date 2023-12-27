package main

import (
	"github.com/sirupsen/logrus"

	"github.com/mustthink/go-storage-like-redis/client"
)

func main() {
	log := logrus.New()
	storageClient := client.New("localhost", "8081")

	if err := storageClient.SetTimeless("", "1", 1); err != nil {
		log.Fatalf("client couldn't set object w err: %s", err.Error())
	}

	var obj int
	if err := storageClient.Get("", "1", &obj); err != nil {
		log.Fatalf("client couldn't get object w err: %s", err.Error())
	}

	if err := storageClient.Delete("", "1"); err != nil {
		log.Fatalf("client couldn't delete object w err: %s", err.Error())
	}
}
