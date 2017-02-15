/*
* @Author: CJ Ting
* @Date: 2017-02-14 22:09:13
* @Email: fatelovely1128@gmail.com
 */

// output code poster in html using multiple divs
package main

import (
	"bytes"
	"fmt"
	"text/template"
)

var domTemplateString = `
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
      font-size: {{ .FontSize }};
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

var domTemplate *template.Template

func init() {
	domTemplate = template.Must(template.New("dom").Parse(domTemplateString))
}

func domOutput(cp *CodePoster) string {
	var buffer bytes.Buffer
	for row := 0; row < cp.rows; row++ {
		for col := 0; col < cp.cols; col++ {
			index := row*cp.cols + col
			data := cp.data[index]
			buffer.WriteString(fmt.Sprintf(
				`<div style="color: %s;">%c</div>`,
				colorToString(data.color),
				data.char,
			))
		}
	}

	var output bytes.Buffer

	err := domTemplate.Execute(&output, &struct {
		Width      int
		Height     int
		FontFamily string
		FontSize   string
		Content    string
	}{
		cp.width,
		cp.height,
		cp.fontFamily,
		cp.fontSize,
		buffer.String(),
	})
	if err != nil {
		fatalln("compile template error:", err)
	}
	return buffer.String()
}
