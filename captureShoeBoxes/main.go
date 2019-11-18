package main

import (
	"database/sql"
	"fmt"
	"github.com/blackjack/webcam"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"time"
)

/*
	Copyright (c) 2019 Marco
	MIT License
    https://github.com/Marco3jp/is-exist-researcher
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
	db := connectDb()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	defer cam.Close()

	setupCamera()

	fileName := strconv.FormatInt(time.Now().Unix(), 10) + ".jpg"
	frame := takeCapture(cam, config)
	saveFrame(frame, fileName)

	// TODO: 参照の問題なのかFrameを保存する前にStopするとPanicで落ちるので暫定的にここで止めている
	err = cam.StopStreaming()

	if err != nil {
		panic(err.Error())
	}

	insertDb(db, fileName)
}

func connectDb() *sql.DB {
	dataSourceName := "capture_shoe_boxes:" + os.Getenv("CAPTURE_SHOE_BOX_SQL_PASSWORD") + "@/is_exist_researcher"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func insertDb(db *sql.DB, fileName string) (result sql.Result) {
	ins, err := db.Prepare("INSERT INTO capture(file_name) VALUES(?)")
	if err != nil {
		panic(err)
	}

	result, err = ins.Exec(fileName)
	if err != nil {
		panic(err)
	}

	return result
}

// 自動補正による画像の差異を減らすために、カメラのホワイトバランスなどをマニュアルでセット
// また、画像のフォーマットにMotion-JPEGを指定
func setupCamera() {
	controlMap := cam.GetControls()
	for f, s := range controlMap {
		fmt.Printf("%#v: %#v\n", f, s)
		switch s.Name {
		case "White Balance Temperature, Auto":
			err := cam.SetControl(f, 0)
			if err != nil {
				fmt.Println("cannnot set White Balance Auto Control")
			} else {
				fmt.Println("success setting White Balance Auto")
			}
			break
		case "White Balance Temperature":
			err := cam.SetControl(f, 3300)
			if err != nil {
				fmt.Println("cannnot set White Balance Temperature Control")
			} else {
				fmt.Println("success setting White Balance Temperature")
			}
			break
		case "Brightness":
			err := cam.SetControl(f, 64)
			if err != nil {
				fmt.Println("cannnot set Brightness Control")
			} else {
				fmt.Println("success setting Brightness")
			}
			break
		case "Sharpness":
			err := cam.SetControl(f, 7)
			if err != nil {
				fmt.Println("cannnot set Sharpness Control")
			} else {
				fmt.Println("success setting Sharpness")
			}
			break
		case "Exposure, Auto":
			err := cam.SetControl(f, 0)
			if err != nil {
				fmt.Println("cannnot set Exposure Auto Control")
			} else {
				fmt.Println("success setting Exposure Auto")
			}
			break
		case "Exposure (Absolute)":
			err := cam.SetControl(f, 300)
			if err != nil {
				fmt.Println("cannnot set Exposure (Absolute) Control")
			} else {
				fmt.Println("success setting Exposure (Absolute)")
			}
			break
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
