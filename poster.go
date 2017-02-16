/*
* @Author: CJ Ting
* @Date: 2017-02-14 19:14:09
* @Email: fatelovely1128@gmail.com
 */

package main

import (
	"fmt"
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
	fontSize   string
	// derived fields
	rows int
	cols int
	data []*PostData
}

type PostData struct {
	char  rune
	color color.Color
}

// we can make this concurrent, but it's not necessary
// cause this is fast enough
func (cp *CodePoster) render() []*PostData {
	var data []*PostData
	for row := 0; row < cp.rows; row++ {
		for col := 0; col < cp.cols; col++ {
			index := (row*cp.cols + col) % len(cp.source)
			data = append(data, &PostData{
				cp.source[index],
				cp.getColor(row, col),
			})
		}
	}
	return data
}

func (cp *CodePoster) getColor(row, col int) color.Color {
	// pixel coordinates
	x := col*cp.charWidth + cp.charWidth/2
	y := row*cp.charHeight + cp.charHeight/2

	// coordinates relative to image
	imgX := x - (cp.width-cp.img.Bounds().Max.X)/2
	imgY := y - (cp.height-cp.img.Bounds().Max.Y)/2

	// outside of image
	if imgX < 0 || imgX > cp.img.Bounds().Max.X {
		return nil
	}
	if imgY < 0 || imgY > cp.img.Bounds().Max.Y {
		return nil
	}

	result := cp.img.At(imgX, imgY)
	_, _, _, a := result.RGBA()

	// full transparent
	if a == 0 {
		return nil
	}

	return result
}

func newCodePoster(
	source string, // source code
	img image.Image, // poster image
	fontFamily string,
	fontSize string,
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

	cp := &CodePoster{
		source:     []rune(regex.ReplaceAllString(source, "")),
		img:        img,
		charWidth:  charWidth,
		charHeight: charHeight,
		width:      width,
		height:     height,
		fontFamily: fontFamily,
		fontSize:   fontSize,
		rows:       height / charHeight,
		cols:       width / charWidth,
	}

	cp.data = cp.render()

	return cp
}

// return bgColor if color is nil
func colorToString(c color.Color) string {
	if c == nil {
		return bgColor
	}
	r, g, b, a := c.RGBA()
	r_ := int((float64(r) / 0xffff) * 0xff)
	g_ := int((float64(g) / 0xffff) * 0xff)
	b_ := int((float64(b) / 0xffff) * 0xff)
	a_ := float64(a) / 0xffff
	return fmt.Sprintf("rgba(%d,%d,%d,%f)", r_, g_, b_, a_)
}
