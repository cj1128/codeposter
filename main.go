/*
* @Author: CJ Ting
* @Date: 2017-02-14 19:10:43
* @Email: fatelovely1128@gmail.com
 */

package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	fp "path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

// set by -ldflags
var appVersion string

var (
	sourcePath   string
	imgPath      string
	charWidth    int
	charHeight   int
	width        int
	height       int
	font         string
	fontSize     string
	bgColor      string
	outputFormat string
)

func parseFlags() {
	kingpin.Flag("font", "font family, please use monospace font,").
		Default("Hack").
		StringVar(&font)
	// font size must corresponding to char width and char height
	kingpin.Flag("fontsize", "font size, valid css font size, must corresponding to char width and char height").
		Default("11.65px").
		StringVar(&fontSize)
	kingpin.Flag("charwidth", "single character width in pixels").
		Default("7").
		IntVar(&charWidth)
	kingpin.Flag("charheight", "single character height in pixels").
		Default("14").
		IntVar(&charHeight)

	kingpin.Flag("width", "output poster width in pixels").
		Default("800").
		IntVar(&width)
	kingpin.Flag("height", "output poster height in pixels").
		Default("760").
		IntVar(&height)

	kingpin.Flag("bgcolor", "background color, valid css color").
		Default("#eee").
		StringVar(&bgColor)
	kingpin.Flag("output", "specify output format, [canvas | dom]").
		Default("canvas").
		EnumVar(&outputFormat, "canvas", "dom")
	kingpin.CommandLine.HelpFlag.Short('h')

	kingpin.Arg("source", "source code path").
		Required().
		StringVar(&sourcePath)
	kingpin.Arg("image", "image path").
		Required().
		StringVar(&imgPath)

	kingpin.Version(appVersion)
	kingpin.Parse()
}

func fatalln(args ...interface{}) {
	log.Println(args...)
	os.Exit(1)
}

func main() {
	parseFlags()

	source, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		fatalln("open source code error:", err)
	}

	imgFile, err := os.Open(imgPath)
	if err != nil {
		fatalln("open image error:", err)
	}
	defer imgFile.Close()

	var img image.Image
	ext := fp.Ext(imgPath)
	switch ext {
	case ".jpg":
		img, err = jpeg.Decode(imgFile)
	case ".png":
		img, err = png.Decode(imgFile)
	case ".gif":
		img, err = gif.Decode(imgFile)
	default:
		img, err = jpeg.Decode(imgFile)
	}

	if err != nil {
		fatalln("decode image error:", err)
	}

	codePoster := newCodePoster(
		string(source),
		img,
		font,
		fontSize,
		charWidth,
		charHeight,
		width,
		height,
	)
	var output string
	if outputFormat == "canvas" {
		output = canvasOutput(codePoster)
	}
	if outputFormat == "dom" {
		output = domOutput(codePoster)
	}
	fmt.Println(output)
}
