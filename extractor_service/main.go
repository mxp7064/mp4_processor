package main

import (
	ex "extractor-service/extractor"
	"os"
	"runtime"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func init() {
	loggingLevel, err := log.ParseLevel(os.Getenv("LOGGING_LEVEL"))
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(loggingLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {

	// construct a InitSegmentExtractorImplementation struct which implements the InitSegmentExtractor interface
	// this way we can easily swap implementations in the future
	var initSegmentExtractor ex.InitSegmentExtractor = ex.InitSegmentExtractorImplementation{}

	const natsURL = "nats://localhost:4222"
	const subject = "init-segment"

	natsOptions := nats.Options{
		Url:      natsURL,
		User:     os.Getenv("NATS_USER"),
		Password: os.Getenv("NATS_PASS"),
	}

	// connect to nats server
	conn, err := natsOptions.Connect()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Go subscriber connected to NATS server")

	// subscribe to the init-segment subject and process incoming file paths
	conn.Subscribe(subject, func(msg *nats.Msg) {
		mp4Path := string(msg.Data)
		log.Debugf("Got message: '%s\n", mp4Path+"'")
		var reply string

		// extracit initialization segment
		initSegmentFilePath, err := initSegmentExtractor.ExtractInitSegment(mp4Path)

		// construct and publish appropriate reply message - initialization segment file path or error message
		if err != nil {
			reply = "An error occurred during file proccessing: " + err.Error()
			log.Error(err)
		} else {
			reply = initSegmentFilePath
		}
		log.Debugf("Sending reply: '%s'\n", reply)
		conn.Publish(msg.Reply, []byte(reply))
	})
	log.Infof("Subscribed to subject: '%s'\n", subject)

	// Calling Goexit from the main goroutine terminates that goroutine without func main returning.
	// Since func main has not returned, the program continues execution of other goroutines
	// this way the program won't exit so we can continue receiving messages from nats
	runtime.Goexit()
}
