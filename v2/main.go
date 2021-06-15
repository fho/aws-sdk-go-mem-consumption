package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

var parallelCount int
var parallelLoops int

var s3Uploader *manager.Uploader

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
	cfg, err := config.LoadDefaultConfig(context.Background())
	exitOnErr(err)
	clt := s3.NewFromConfig(cfg)

	s3Uploader = manager.NewUploader(clt, func(u *manager.Uploader) {
		//u.BufferProvider = manager.NewBufferedReadSeekerWriteToPool(512 * 1024)
	})
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
	data := []byte("1")

	if withReadSeeker {
		reader = bytes.NewReader(data)
	} else {
		reader = bytes.NewBuffer(data)
	}

	_, err := s3Uploader.Upload(context.Background(), &s3.PutObjectInput{
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
