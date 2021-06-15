package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var parallelCount int
var parallelLoops int

var s3Uploader *s3manager.Uploader

var bucketName = "sisu-file-service"
var objectKeyPrefix = "workspace/uploadTest"

var withReadSeeker bool

func main() {
	flag.IntVar(&parallelCount, "parallel-uploads", 5, "number of parallel uploads")
	flag.IntVar(&parallelLoops, "repeat", 10, "how many times to repeat the parallel uploads")
	flag.BoolVar(&withReadSeeker, "with-read-seeker", false, "use a read-seeker instead of a reader buffer")

	flag.Parse()
	log.Printf("WithReadSeeker: %+v\n", withReadSeeker)

	initS3Uploader()

	runParallel(s3lib)
}

func initS3Uploader() {
	awsSession, err := session.NewSession(aws.NewConfig())
	if err != nil {
		log.Fatalln(err)
	}

	s3Uploader = s3manager.NewUploader(awsSession)
}

func runParallel(fn func()) {
	for j := 0; j < parallelLoops; j++ {
		var wg sync.WaitGroup

		for i := 0; i < parallelCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				fn()
			}()
		}

		wg.Wait()
	}

	log.Println("runParallel DONE")
}

func s3lib() {
	var reader io.Reader
	data := []byte("1234567890")

	if withReadSeeker {
		reader = bytes.NewReader(data)
	} else {
		reader = bytes.NewBuffer(data)
	}

	_, err := s3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKeyPrefix + fmt.Sprint(time.Now().UnixNano())),
		Body:   reader,
	})
	exitOnErr(err)

	log.Println("upload finished")
}

func exitOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
