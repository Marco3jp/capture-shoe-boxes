package main

import (
	"database/sql"
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"strconv"
	"time"
)

type Config struct {
	voidShoeBoxPath string
	threshold       uint // uint16とか？
	metricType      imagick.MetricType
}

var config = Config{
	voidShoeBoxPath: "./test/void.jpg",
	threshold:       100,
	metricType:      imagick.METRIC_FUZZ_ERROR,
}

type DiffImageResult struct {
	difference float64
}

// for debug //
type DebugState struct {
	debug       bool
	debugConfig DebugStruct
}

type DebugStruct struct {
	oneImagePath string
	twoImagePath string
}

var debugState = DebugState{
	debug: false,
	debugConfig: DebugStruct{
		oneImagePath: "/tmp/1573657900.jpg",
		twoImagePath: "/tmp/1573657906.jpg",
	},
}

func main() {
	// db, err := connectDb()
	// latestImagePath := getLatestImagePath(db)
	latestImagePath := "./test/exist_1.jpg"

	if len(os.Args) >= 2 && os.Args[1] == "debug" {
		fmt.Println("Debug mode")
		debugState.debug = true
		config.voidShoeBoxPath = debugState.debugConfig.oneImagePath
		latestImagePath = debugState.debugConfig.twoImagePath
	}

	result := diffImage(latestImagePath)
	fmt.Printf("%#v", result)
}

func connectDb() *sql.DB {
	db, err := sql.Open("mysql", "test:passwd@/example")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	return db
}

func getLatestImagePath(db *sql.DB) string {
	var latestImagePath string

	row := db.QueryRow("select fileName from capture order by id desc limit 1 ")

	err := row.Scan(&latestImagePath)
	if err != nil {
		panic(err)
	}

	return latestImagePath
}

func getImage(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	photoImage, err := jpeg.Decode(file)
	if err != nil {
		panic(err)
	}

	return photoImage
}

func convertToGrayscale(colorImage image.Image) (grayImage *image.Gray16) {
	bounds := colorImage.Bounds()
	grayImage = image.NewGray16(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.Gray16Model.Convert(colorImage.At(x, y))
			gray, _ := c.(color.Gray16)
			grayImage.Set(x, y, gray)
		}
	}
	return grayImage
}

func diffImage(currentImagePath string) diffImageResult {
	imagick.Initialize()
	defer imagick.Terminate()

	voidBox := imagick.NewMagickWand()
	defer voidBox.Destroy()
	err := voidBox.SetColorspace(imagick.COLORSPACE_GRAY)
	if err != nil {
		panic(err)
	}

	currentBox := imagick.NewMagickWand()
	defer currentBox.Destroy()
	err = currentBox.SetColorspace(imagick.COLORSPACE_GRAY)
	if err != nil {
		panic(err)
	}

	err = voidBox.ReadImage(config.voidShoeBoxPath)
	if err != nil {
		panic(err)
	}
	err = currentBox.ReadImage(currentImagePath)
	if err != nil {
		panic(err)
	}

	if debugState.debug {
		fileA, err := os.Create("/tmp/" + strconv.FormatInt(time.Now().Unix(), 10) + "_void.jpg")
		if err != nil {
			panic(err)
		}
		fileB, err := os.Create("/tmp/" + strconv.FormatInt(time.Now().Unix(), 10) + "_fill.jpg")
		if err != nil {
			panic(err)
		}
		voidBox.WriteImageFile(fileA)
		currentBox.WriteImageFile(fileB)
	}

	_, result := currentBox.CompareImages(voidBox, config.metricType)
	return diffImageResult{difference: result}
}
