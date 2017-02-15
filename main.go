/*
* @Author: CJ Ting
* @Date: 2017-02-14 19:10:43
* @Email: fatelovely1128@gmail.com
 */

package main

import (
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"

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
	kingpin.Flag("charwidth", "single character width in pixels").
		Default("7").
		IntVar(&charWidth)
	kingpin.Flag("charheight", "single character height in pixels").
		Default("14").
		IntVar(&charHeight)
	// font size must corresponding to char width and char height
	kingpin.Flag("fontsize", "font size, valid css unit, must corresponding to char width and char height").
		Default("11.65px").
		StringVar(&fontSize)

	kingpin.Flag("width", "output poster width in pixels").
		Default("800").
		IntVar(&width)
	kingpin.Flag("height", "output poster height in pixels").
		Default("760").
		IntVar(&height)

	kingpin.Flag("font", "font family, please use monospace font,").
		Default("Hack").
		StringVar(&font)
	kingpin.Flag("bgcolor", "background color, valid css unit").
		Default("#eee").
		StringVar(&bgColor)
	kingpin.Flag("output", "specify output format, [canvs | dom]").
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

	img, _, err := image.Decode(imgFile)
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
