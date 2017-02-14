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
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: code-poster [source] [image]")
		os.Exit(1)
	}
	sourcePath := os.Args[1]
	imgPath := os.Args[2]

	source, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		log.Fatal(err)
	}

	imgFile, err := os.Open(imgPath)
	if err != nil {
		log.Fatal(err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
	}

	codePoster := newCodePoster(
		string(source),
		img,
		"Hack",
		8.3,
		5,
		10,
		800,
		800,
	)
	output := htmlOutput(codePoster)
	fmt.Println(output)
}
