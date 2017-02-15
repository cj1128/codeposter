/*
* @Author: CJ Ting
* @Date: 2017-02-15 11:22:55
* @Email: fatelovely1128@gmail.com
 */

// output code poster in html using canvas
package main

import (
	"bytes"
	"fmt"
	"text/template"
)

var canvasTemplateString = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>CodePoster</title>
  <style type="text/css">
  #container {
    position: fixed;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    margin: auto;
  }
  </style>
</head>
<body>
  <canvas id="container" width="{{ .Width }}" height="{{ .Height }}">
  </canvas>

  <script type="text/javascript">
    var ctx = document.getElementById("container").getContext("2d")
    ctx.font = "{{ .FontSize }} {{ .FontFamily }}"
    var charWidth = {{ .CharWidth }}
    var charHeight = {{ .CharHeight }}
    var rows = {{ .Rows }}
    var cols = {{ .Cols }}
    var data = {{ .Data }}
    for(row = 0; row < rows; row++) {
      for(col = 0; col < cols; col++) {
        var x = col * charWidth
        var y = row * charHeight + charHeight
        var index = row * cols + col
        var item = data[index]
        ctx.fillStyle = item[1]
        ctx.fillText(item[0], x, y)
      }
    }
  </script>
</body>
</html>
`

var canvasTemplate *template.Template

func init() {
	canvasTemplate = template.Must(template.New("canvas").Parse(canvasTemplateString))
}

func canvasOutput(cp *CodePoster) string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	for row := 0; row < cp.rows; row++ {
		for col := 0; col < cp.cols; col++ {
			index := row*cp.cols + col
			data := cp.data[index]
			char := data.char
			str := string(char)
			// we need to escape double quote
			if char == '"' {
				str = `\"`
			}
			if char == '\\' {
				str = `\\`
			}
			buffer.WriteString(fmt.Sprintf(
				`["%s", "%s"],`,
				str,
				colorToString(data.color),
			))
		}
	}
	buffer.WriteString("]")
	var output bytes.Buffer

	err := canvasTemplate.Execute(&output, &struct {
		Width      int
		Height     int
		FontSize   string
		FontFamily string
		CharWidth  int
		CharHeight int
		Rows       int
		Cols       int
		Data       string
	}{
		cp.width,
		cp.height,
		cp.fontSize,
		cp.fontFamily,
		cp.charWidth,
		cp.charHeight,
		cp.rows,
		cp.cols,
		buffer.String(),
	})

	if err != nil {
		fatalln("compile template error:", err)
	}

	return output.String()
}
