package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"extractor-service/extractor"
	"github.com/nats-io/nats.go"
)

func main() {

	var initSegmentExtractor exctractor.InitSegmentExtractor = exctractor.InitSegmentExtractorImplementation{}

	const natsURL = "nats://localhost:4222"
	const subject = "init-segment"

	opts := nats.Options{
		Url:      natsURL,
		User:     os.Getenv("NATS_USER"),
		Password: os.Getenv("NATS_PASS"),
	}

	conn, err := opts.Connect()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Go subscriber connected to NATS server")

	conn.Subscribe(subject, func(msg *nats.Msg) {
		mp4Path := string(msg.Data)
		fmt.Printf("Got message: '%s\n", mp4Path+"'")
		var reply string

		initSegmentFilePath, err := initSegmentExtractor.ExtractInitSegment(mp4Path)

		if err != nil {
			reply = "An error occurred during file proccessing: " + err.Error()
		} else {
			reply = initSegmentFilePath
		}
		fmt.Printf("Sending reply: '%s'\n", reply)
		conn.Publish(msg.Reply, []byte(reply))
	})
	fmt.Printf("Subscribed to subject: '%s'\n", subject)
	runtime.Goexit()
}
