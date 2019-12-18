package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gographics/imagick.v3/imagick"
	"image"
	"image/jpeg"
	"os"
	"strconv"
	"time"
)

/*
靴箱は4*6、思ったよりもカメラの画角が広いので、そのために切り出す必要があるかなぁという感じ。
*/

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
	db, err := connectDb()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// TODO: ROWとCOLUMNの対応
	var diffImageResult = DiffImageResult{
		comparedScore: 0,
		isExist:       false,
		livingTimes:   0,
		row:           0,
		column:        0,
		captureId:     0,
	}

	if len(os.Args) >= 2 && os.Args[1] == "debug" {
		fmt.Println("Debug mode")
		debugState.debug = true
		config.voidShoeBoxPath = config.imageRoot + debugState.debugConfig.voidImageFileName
	}

	// TODO: 何故か := による定義だとエラーを吐く
	var latestImageName string
	diffImageResult.captureId, latestImageName = getLatestImageName(db)

	latestImagePath := config.imageRoot + latestImageName

	imagick.Initialize()
	defer imagick.Terminate()
	voidBox, currentBox := setupImagick(latestImagePath)
	defer voidBox.Destroy()
	defer currentBox.Destroy()

	// グレースケール化
	setColorspace(voidBox, imagick.COLORSPACE_GRAY)
	setColorspace(currentBox, imagick.COLORSPACE_GRAY)

	// 比較処理
	diffImageResult.comparedScore = diffImage(voidBox, currentBox)
	diffImageResult.isExist = isExist(diffImageResult.comparedScore)

	// 何回生きてたかのチェック。初期値0なので非存在の場合はそのまま。
	if diffImageResult.isExist {
		diffImageResult.livingTimes = getLatestLivingTimes(db, diffImageResult.row, diffImageResult.column) + 1
	}

	// DBにDiff結果を挿入
	insertDiffResult(db, diffImageResult)
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

// 比較の結果をDBに挿入
func insertDiffResult(db *sql.DB, result DiffImageResult) {
	ins, err := db.Prepare("INSERT INTO shoe_box(is_exist, living_times, row, `column`, compared_score, compared_metric, exist_threshold,capture_id) VALUES(?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(err)
	}

	_, err = ins.Exec(result.isExist, result.livingTimes, result.row, result.column, result.comparedScore, config.metricType, config.threshold, result.captureId)
	if err != nil {
		panic(err)
	}
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

func isExist(comparedScore float64) bool {
	return comparedScore > config.threshold
}
