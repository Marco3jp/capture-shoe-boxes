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
	threshold       float64 // 存在とするしきい値
	metricType      imagick.MetricType
	limit           float64 // 計測不能とするしきい値
	imageRoot       string
}

var config = Config{
	voidShoeBoxPath: "",
	threshold:       2300,
	metricType:      imagick.METRIC_MEAN_ABSOLUTE_ERROR,
	limit:           8500,
	imageRoot:       "./test/",
}

type DiffImageResult struct {
	comparedScore float64
	isExist       bool
	livingTimes   uint8
	row           uint8
	column        uint8
	captureId     uint
}

// for debug //
type DebugState struct {
	debug       bool
	debugConfig DebugStruct
}

type DebugStruct struct {
	voidImageFileName string
	targetImageId     uint
}

var debugState = DebugState{
	debug: false,
	debugConfig: DebugStruct{
		voidImageFileName: "cropped_1574072094.jpg",
		targetImageId:     1,
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

func connectDb() (*sql.DB, error) {
	dataSourceName := "diff_shoe_boxes:" + os.Getenv("DIFF_SHOE_BOX_SQL_PASSWORD") + "@/is_exist_researcher"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	return db, err
}

// 直前の撮影データのパスを受け取る
func getLatestImageName(db *sql.DB) (captureId uint, latestImageName string) {
	var row *sql.Row
	if debugState.debug {
		row = db.QueryRow("select id, file_name from capture where id=?", debugState.debugConfig.targetImageId)
	} else {
		row = db.QueryRow("select id, file_name from capture order by id desc limit 1 ")
	}

	err := row.Scan(&captureId, &latestImageName)
	if err != nil {
		panic(err)
	}

	return captureId, latestImageName
}

// 直前のデータでその靴箱が何連続生存だったか受け取る
func getLatestLivingTimes(db *sql.DB, tableRow uint8, column uint8) (livingTimes uint8) {
	row := db.QueryRow("select living_times from shoe_box where row=? and `column`=? order by id desc limit 1", tableRow, column)

	err := row.Scan(&livingTimes)
	if err != nil {
		panic(err)
	}

	return livingTimes
}

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

// imagickの初期化を行い、MagickWandを生成、画像を読み込んで返す
func setupImagick(currentImagePath string) (voidBox *imagick.MagickWand, currentBox *imagick.MagickWand) {
	voidBox = imagick.NewMagickWand()

	currentBox = imagick.NewMagickWand()

	err := voidBox.ReadImage(config.voidShoeBoxPath)
	if err != nil {
		panic(err)
	}
	err = currentBox.ReadImage(currentImagePath)
	if err != nil {
		panic(err)
	}

	return voidBox, currentBox
}

func setColorspace(target *imagick.MagickWand, colorspace imagick.ColorspaceType) {
	err := target.SetImageColorspace(colorspace)
	if err != nil {
		panic(err)
	}
}

func diffImage(voidBox *imagick.MagickWand, currentBox *imagick.MagickWand) float64 {
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

	_, compareScore := currentBox.CompareImages(voidBox, config.metricType)
	return compareScore
}
