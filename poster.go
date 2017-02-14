/*
* @Author: CJ Ting
* @Date: 2017-02-14 19:14:09
* @Email: fatelovely1128@gmail.com
 */

package main

import (
	"image"
	"image/color"
	"regexp"

	"github.com/nfnt/resize"
)

type CodePoster struct {
	source     []rune
	img        image.Image
	charWidth  int
	charHeight int
	width      int
	height     int
	fontFamily string
	fontSize   float64
}

var bgColor = color.RGBA{
	R: 0xee,
	G: 0xee,
	B: 0xee,
	A: 0xee,
}

// get character and color of position at row and col
func (cp *CodePoster) Get(row, col int) (rune, color.Color) {
	index := row*cp.Cols() + col
	return cp.source[index], cp.getColor(row, col)
}

func (cp *CodePoster) getColor(row, col int) color.Color {
	// pixel coordinate
	x := col*cp.charWidth + cp.charWidth/2
	y := row*cp.charHeight + cp.charHeight/2
	imgX := x - (cp.width-cp.img.Bounds().Max.X)/2
	imgY := y - (cp.height-cp.img.Bounds().Max.Y)/2
	if imgX < 0 || imgX > cp.img.Bounds().Max.X {
		return bgColor
	}
	if imgY < 0 || imgY > cp.img.Bounds().Max.Y {
		return bgColor
	}
	result := cp.img.At(imgX, imgY)
	_, _, _, a := result.RGBA()
	if a == 0 {
		return bgColor
	}
	return result
}

func (cp *CodePoster) Rows() int {
	return cp.height / cp.charHeight
}

func (cp *CodePoster) Cols() int {
	return cp.width / cp.charWidth
}

func newCodePoster(
	source string, // source code
	img image.Image, // poster image
	fontFamily string,
	fontSize float64,
	charWidth, charHeight int, // single character width and height, in pixel
	width, height int, // the whole poster width and height, in pixel
) *CodePoster {
	regex := regexp.MustCompile(`\s*`)

	// we need to scale image
	imgWidth := img.Bounds().Max.X
	imgHeight := img.Bounds().Max.Y

	if imgWidth > width || imgHeight > height {
		imgRatio := (float64)(imgWidth) / (float64)(imgHeight)
		ratio := (float64)(width) / (float64)(height)
		if imgRatio > ratio {
			img = resize.Resize(uint(width), 0, img, resize.Lanczos3)
		} else {
			img = resize.Resize(0, uint(height), img, resize.Lanczos3)
		}
	}

	return &CodePoster{
		source:     []rune(regex.ReplaceAllString(source, "")),
		img:        img,
		charWidth:  charWidth,
		charHeight: charHeight,
		width:      width,
		height:     height,
		fontFamily: fontFamily,
		fontSize:   fontSize,
	}
}
