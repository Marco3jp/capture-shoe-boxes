package main

import (
	"database/sql"
	"fmt"
	"github.com/blackjack/webcam"
	"os"
	"strconv"
	"time"
)

/*
	Copyright (c) 2019 Marco
	MIT License
    https://github.com/Marco3jp/capture-shoe-boxes [WIP]
*/

type CaptureSize struct {
	height uint32
	width  uint32
}

type Config struct {
	captureRoot  string      // Set directory to save captured image.
	captureSize  CaptureSize // {width}x{height} [notice: "x" character for splitting]
	cameraDevice string      // like /dev/videoN, N is "0", "1"...etc(it is auto detected by udev.).
	timeOut      uint32      // this unit is unknown...it is mentioned next code. https://github.com/blackjack/webcam/blob/master/v4l2.go#L445)
}

var config = Config{
	captureRoot: "/tmp/",
	captureSize: CaptureSize{
		height: 720,
		width:  1280,
	},
	cameraDevice: "/dev/video0",
	timeOut:      1000000,
}

var cam, err = webcam.Open(config.cameraDevice)

/*
	This program require Linux4Video2 API.
	If you don't implementation api, cannot move on Windows family and other UNIX Family.
	Use Linux, or future Linux like OS.
*/
func main() {
	// db := connectDb()
	if err != nil {
		panic(err)
	}
	defer cam.Close()

	setupCamera()

	fileName := strconv.FormatInt(time.Now().Unix(), 10) + ".jpg"
	frame := takeCapture(cam, config)
	saveFrame(frame, fileName)

	err = cam.StopStreaming()

	if err != nil {
		panic(err.Error())
	}

	// db.Prepare("INSERT INTO capture(fileName, createdAt, updatedAt) VALUES(?,?,?)")
}

func connectDb() *sql.DB {
	db, err := sql.Open("mysql", "test:passwd@/example")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	return db
}

func setupCamera() {
	controlMap := cam.GetControls()
	for f, s := range controlMap {
		if s.Name == "White Balance Temperature, Auto" {
			err := cam.SetControl(f, 0)
			if err != nil {
				fmt.Println("cannnot set White Balance Auto Control")
			}
		}

		if s.Name == "White Balance Temperature" {
			err := cam.SetControl(f, 3400)
			if err != nil {
				fmt.Println("cannnot set White Balance Temperature Control")
			}
		}
	}

	var format webcam.PixelFormat
	formatMap := cam.GetSupportedFormats()
	for f, s := range formatMap {
		if s == "Motion-JPEG" {
			format = f
			break
		}
	}

	if format == 0 {
		panic("cannnot use format")
	}

	_, _, _, err := cam.SetImageFormat(format, config.captureSize.width, config.captureSize.height)

	if err != nil {
		panic(err)
	}
}

func takeCapture(cam *webcam.Webcam, config Config) []byte {
	err := cam.StartStreaming()

	if err != nil {
		panic(err.Error())
	}

	err = cam.WaitForFrame(config.timeOut)

	switch err.(type) {
	case nil:
	case *webcam.Timeout:
		panic(err.Error())
	default:
		panic(err.Error())
	}

	frame, err := cam.ReadFrame()

	if err != nil {
		panic(err.Error())
	}

	if len(frame) != 0 {
		return frame
	} else {
		panic("frame no data")
	}
}

func saveFrame(frame []byte, fileName string) {
	file, err := os.Create(config.captureRoot + fileName)

	if err != nil {
		panic(err)
	}

	_, err = file.Write(frame)
	if err != nil {
		panic(err)
	}
}
