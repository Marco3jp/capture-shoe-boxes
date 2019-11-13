package main

import (
	"database/sql"
	"fmt"
	"gopkg.in/gographics/imagick.v3/imagick"
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

type Config struct {
	voidShoeBoxPath string
	threshold       uint // uint16とか？
	metricType      imagick.MetricType
}

type diffImageResult struct {
	difference float64
}

var config = Config{
	voidShoeBoxPath: "./test/void.jpg",
	threshold:       100,
	metricType:      imagick.METRIC_MEAN_ABSOLUTE_ERROR,
}

func main() {
	// db, err := connectDb()
	// latestImagePath := getLatestImagePath(db)
	latestImagePath := "./test/exist_1.jpg"
	// latestImage := getImage(latestImagePath)
	// grayImage := convertToGrayscale(latestImage)
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

	_, result := currentBox.CompareImages(voidBox, config.metricType)
	return diffImageResult{difference: result}
}
