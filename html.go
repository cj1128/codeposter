/*
* @Author: CJ Ting
* @Date: 2017-02-14 22:09:13
* @Email: fatelovely1128@gmail.com
 */

package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"strings"
	"text/template"
)

var htmlTemplateString = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>CodePoster</title>
  <style type="text/css">
    .container {
      width: {{ .Width }}px;
      height: {{ .Height }}px;
      margin: auto;
      font-family: "{{ .FontFamily }}";
      font-size: 0;
      position: fixed;
      left: 0;
      right: 0;
      top: 0;
      bottom: 0;
    }
    .container > div {
      display: inline-block;
      font-size: {{ .FontSize }}px;
    }
  </style>
</head>
<body>
  <div class="container">
  {{ .Content }}
  </div>
</body>
</html>
`

var htmlTemplate *template.Template

func init() {
	htmlTemplate = template.Must(template.New("html").Parse(htmlTemplateString))
}

func htmlOutput(cp *CodePoster) string {
	var content []string
	for row := 0; row < cp.Rows(); row++ {
		for col := 0; col < cp.Cols(); col++ {
			char, color := cp.Get(row, col)
			content = append(
				content,
				fmt.Sprintf(
					`<div style="color: %s;">%c</div>`,
					colorToString(color),
					char,
				),
			)
		}
	}

	buffer := &bytes.Buffer{}

	err := htmlTemplate.Execute(buffer, &struct {
		Width      int
		Height     int
		FontFamily string
		FontSize   float64
		Content    string
	}{
		cp.width,
		cp.height,
		cp.fontFamily,
		cp.fontSize,
		strings.Join(content, "\n"),
	})
	if err != nil {
		log.Fatal(err)
	}
	return buffer.String()
}

func colorToString(c color.Color) string {
	r, g, b, a := c.RGBA()
	r_ := int((float64(r) / 0xffff) * 0xff)
	g_ := int((float64(g) / 0xffff) * 0xff)
	b_ := int((float64(b) / 0xffff) * 0xff)
	a_ := float64(a) / 0xffff
	return fmt.Sprintf("rgba(%d,%d,%d,%f)", r_, g_, b_, a_)
}
