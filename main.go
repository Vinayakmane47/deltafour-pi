package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	serial "github.com/tarm/goserial"
)

var count int
var file_count int

func read_int32_big(data []byte) (ret uint8) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func upload(filess chan string, uploader *s3manager.Uploader) {
	for {
		select {
		case filename, ok := <-filess:
			if !ok {
				log.Fatal("error in channel")
				return
			}

			upFile, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
				return
			}
			defer upFile.Close()

			upFileInfo, _ := upFile.Stat()
			var fileSize int64 = upFileInfo.Size()
			fileBuffer := make([]byte, fileSize)
			upFile.Read(fileBuffer)

			input := &s3manager.UploadInput{
				Bucket:      aws.String("isolation-point-images"), // bucket's name
				Key:         aws.String("devtesting/" + filename), // files destination location
				Body:        bytes.NewReader(fileBuffer),          // content of the file
				ContentType: aws.String("text"),                   // content type
			}
			_, errs := uploader.UploadWithContext(context.Background(), input)
			if errs != nil {
				log.Fatal(errs)
				return
			} else {
				e := os.Remove(filename)
				if e != nil {
					log.Fatal(e)
					return
				}
			}
		}
	}
}

func readDataCh(ints chan uint8, val string, pi_files chan string) {

	for {
		select {
		case d, ok := <-ints:
			if !ok {
				log.Fatal("error in channel")
				return
			}
			val = val + strconv.Itoa(int(d)) + ","
			fmt.Println(strconv.Itoa(int(d)))
			if count%10000 == 0 {
				f, err := os.Create("data" + strconv.Itoa(file_count) + ".txt")

				if err != nil {
					log.Fatal(err)
				}

				data := []byte(val)

				_, err2 := f.Write(data)

				if err2 != nil {
					log.Fatal(err2)
				}
				f.Close()
				val = ""
				pi_files <- "data" + strconv.Itoa(file_count) + ".txt"
				file_count = file_count + 1
			}
			count = count + 1
		}
	}
}

func main() {

	s3Config := &aws.Config{
		Region:      aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials("AKIAYQ5GABTRPHZY32DP", "9meuXNrV+xHcQAc1M8VbEnyGP0IdTi7OgQbnrnr2", ""),
	}
	s3Session := session.New(s3Config)

	uploader := s3manager.NewUploader(s3Session)

	count = 1
	file_count = 0
	pi_channel := make(chan uint8)
	pi_files := make(chan string)

	val := ""
	go upload(pi_files, uploader)
	go readDataCh(pi_channel, val, pi_files)

	c := &serial.Config{Name: "/dev/ttyACM0", Baud: 115200, ReadTimeout: time.Millisecond * 5000}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	// defer s.Close()

	go func() {
		for {
			buf := make([]byte, 1)
			n, err := s.Read(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
			if n <= 0 {
				fmt.Println("got 0")
			}
			fmt.Println(uint8(buf[0]))
			pi_channel <- read_int32_big(buf[:n])
		}
	}()

	go func() {
		for {
			buf := make([]byte, 1)
			n, err := s.Read(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
			if n <= 0 {
				fmt.Println("got 0")
			}
			fmt.Println(uint8(buf[0]))
			pi_channel <- read_int32_big(buf[:n])
		}
	}()
}
