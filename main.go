package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/nfnt/resize"
	sdlimg "github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"gopkg.in/alecthomas/kingpin.v2"
)

// set by -ldflags
var appVersion string

const defaultFont = "Hack-Regular.ttf"
const defaultImage = "gopher.png"

var config struct {
	sourcePath string
	imgPath    string
	fontPath   string
	fontSize   int
	bgColor    color
	codeColor  color
	padding    padding
	width      int // in chars
	height     int // in chars
}

type color sdl.Color

func (c *color) Set(value string) error {
	if !colorReg.MatchString(value) {
		return errInvalidColor
	}

	var r, g, b, a string

	if len(value) == 4 {
		r = value[1:2] + value[1:2]
		g = value[2:3] + value[2:3]
		b = value[3:4] + value[3:4]
		a = "ff"
	} else {
		r = value[1:3]
		g = value[3:5]
		b = value[5:7]
		a = "ff"

		if len(value) == 9 {
			a = value[7:9]
		}
	}

	var parsed int64

	parsed, _ = strconv.ParseInt(r, 16, 32)
	c.R = uint8(parsed)

	parsed, _ = strconv.ParseInt(g, 16, 32)
	c.G = uint8(parsed)

	parsed, _ = strconv.ParseInt(b, 16, 32)
	c.B = uint8(parsed)

	parsed, _ = strconv.ParseInt(a, 16, 32)
	c.A = uint8(parsed)

	return nil
}

func (c *color) String() string {
	return fmt.Sprintf("#%2x%2x%2x%2x(rgba)", c.R, c.G, c.B, c.A)
}

type padding struct {
	horizontal int
	vertical   int
}

func (p *padding) Set(value string) error {
	parts := strings.Split(value, ",")
	nums := make([]int, len(parts))

	for i, p := range parts {
		num, err := strconv.Atoi(p)
		if err != nil {
			return errors.Wrap(err, "could not parse string to integer")
		}
		nums[i] = num
	}

	if len(nums) == 1 {
		p.vertical = nums[0]
		p.horizontal = nums[0]
	} else if len(nums) == 2 {
		p.vertical = nums[0]
		p.horizontal = nums[1]
	} else {
		return errors.New("invalid padding value, should be x or x,x")
	}

	return nil
}

func (p *padding) String() string {
	return fmt.Sprintf("[vertical: %d, horizontal: %d]", p.vertical, p.horizontal)
}

var colorReg = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)

var errInvalidColor = errors.New("color should be '#rgb' or '#rrggbb' or '#rrggbbaa'")

type fontTexture struct {
	texture *sdl.Texture
	w       int32
	h       int32
}

type sdlContext struct {
	win           *sdl.Window
	winSurface    *sdl.Surface
	renderer      *sdl.Renderer
	font          *ttf.Font
	padding       padding
	charWidth     int
	charHeight    int
	winWidth      int
	winHeight     int
	contentWidth  int // without padding
	contentHeight int // without padding
}

// char: ASCII character, 0x20 ~ 0x7e inclusive
func renderChar(font *ttf.Font, renderer *sdl.Renderer, char byte, color sdl.Color) (*fontTexture, error) {
	result := &fontTexture{}

	surface, err := font.RenderUTF8Blended(string(char), color)
	if err != nil {
		return nil, errors.Wrap(err, "could not render font")
	}

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, errors.Wrap(err, "could not create texture from surface")
	}

	result.texture = texture
	result.w = surface.W
	result.h = surface.H

	return result, nil
}

func initFlags() {
	kingpin.Flag("font", fmt.Sprintf("specify font file (default: %s bundled in binary)", defaultFont)).
		StringVar(&config.fontPath)

	kingpin.Flag("font-size", "font size").
		Default("12").
		IntVar(&config.fontSize)

	kingpin.Flag("width", "poster width in characters").
		Default("120").
		IntVar(&config.width)

	kingpin.Flag("height", "poster height in characters").
		Default("50").
		IntVar(&config.height)

	kingpin.Flag("code-color", "source code color, '#rgb' or '#rrggbb' or '#rrggbbaa'").
		Default("#e9e9e9").
		SetValue(&config.codeColor)

	kingpin.Flag("bg-color", "background color, '#rgb' or '#rrggbb' or '#rrggbbaa'").
		Default("#fff").
		SetValue(&config.bgColor)

	kingpin.Flag("img", fmt.Sprintf("image used to render poster (default: %s bundled in binary", defaultImage)).
		StringVar(&config.imgPath)

	kingpin.Flag("padding", "padding space in characters, e.g. 1,2").
		Default("1,2").
		SetValue(&config.padding)

	kingpin.Arg("source", "source code path").
		Required().
		StringVar(&config.sourcePath)

	kingpin.Version(appVersion)
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.CommandLine.VersionFlag.Short('v')
}

func fatalln(args ...interface{}) {
	log.Println(args...)
	os.Exit(1)
}

func main() {
	initFlags()

	if len(os.Args) == 1 {
		kingpin.Usage()
		os.Exit(0)
	}

	kingpin.Parse()

	// print config info
	{
		imgPath := config.imgPath
		if imgPath == "" {
			imgPath = fmt.Sprintf("builtin %s", defaultFont)
		}

		fontPath := config.fontPath
		if fontPath == "" {
			imgPath = fmt.Sprintf("builtin %s", defaultImage)
		}

		log.Printf(`Config:
  source path: %s
  img path: %s
  font path: %s
  font size: %d
  background color: %s
  code color: %s
  width in characters: %d
  height in characters: %d
  padding in characters: %s
`, config.sourcePath,
			imgPath,
			fontPath,
			config.fontSize,
			config.bgColor.String(),
			config.codeColor.String(),
			config.width,
			config.height,
			config.padding.String(),
		)
	}

	if err := run(); err != nil {
		log.Fatalln(err)
	} else {
		log.Println("All done 🎉")
	}
}

func readCode() ([]byte, error) {
	content, err := ioutil.ReadFile(config.sourcePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not read source code")
	}

	// only keep ascii characters and remove whitespaces
	var buf bytes.Buffer

	for _, b := range content {
		if b >= 0x20 && b <= 0x7e && b != ' ' && b != '\t' && b != '\n' {
			buf.WriteByte(b)
		}
	}

	return buf.Bytes(), nil
}

func openAndResizeImage(contentWidth, contentHeight int) (image.Image, error) {
	var imgReader io.Reader

	if config.imgPath == "" {
		buf, _ := Asset(defaultImage)
		imgReader = bytes.NewReader(buf)
	} else {
		imgFile, err := os.Open(config.imgPath)
		if err != nil {
			return nil, errors.Wrap(err, "could not open image file")
		}
		defer imgFile.Close()

		imgReader = imgFile
	}

	img, _, err := image.Decode(imgReader)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode image")
	}

	// resize image if necessary
	imgWidth := img.Bounds().Max.X
	imgHeight := img.Bounds().Max.Y

	if imgWidth > contentWidth || imgHeight > contentHeight {
		imgRatio := (float64)(imgWidth) / (float64)(imgHeight)
		winRatio := (float64)(contentWidth) / (float64)(contentHeight)
		if imgRatio > winRatio {
			img = resize.Resize(uint(contentWidth), 0, img, resize.Lanczos3)
		} else {
			img = resize.Resize(0, uint(contentHeight), img, resize.Lanczos3)
		}
	}

	return img, nil
}

func initSDLAndTTF() (*sdlContext, error) {
	// init sdl
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, errors.Wrap(err, "could not init sdl")
	}

	// init ttf
	if err := ttf.Init(); err != nil {
		return nil, errors.Wrap(err, "could not init sdl ttf")
	}

	// load default font
	if config.fontPath == "" {
		tmpFile, err := ioutil.TempFile("", "codeposter")
		if err != nil {
			return nil, errors.Wrap(err, "couldn't create temporary file")
		}

		config.fontPath = tmpFile.Name()
		buf, _ := Asset(defaultFont)
		if _, err := tmpFile.Write(buf); err != nil {
			return nil, errors.Wrap(err, "could not write to temporary file")
		}

		tmpFile.Close()
	}

	// open font
	font, err := ttf.OpenFont(config.fontPath, config.fontSize)
	if err != nil {
		return nil, errors.Wrap(err, "could not open font")
	}

	charWidth, charHeight, err := font.SizeUTF8("a")
	if err != nil {
		return nil, errors.Wrap(err, "could not get size of character")
	}

	contentWidth := charWidth * config.width
	contentHeight := charHeight * config.height

	winWidth := contentWidth + config.padding.horizontal*2*charWidth
	winHeight := contentHeight + config.padding.vertical*2*charHeight

	win, err := sdl.CreateWindow("", 0, 0, int32(winWidth), int32(winHeight), sdl.WINDOW_HIDDEN)
	if err != nil {
		return nil, errors.Wrap(err, "could not create sdl window")
	}

	winSurface, err := win.GetSurface()
	if err != nil {
		return nil, errors.Wrap(err, "could not get surface of window")
	}

	renderer, err := sdl.CreateSoftwareRenderer(winSurface)
	if err != nil {
		return nil, errors.Wrap(err, "could not create renderer from window surface")
	}

	return &sdlContext{
		win:           win,
		renderer:      renderer,
		charWidth:     charWidth,
		charHeight:    charHeight,
		winWidth:      winWidth,
		winHeight:     winHeight,
		contentWidth:  contentWidth,
		contentHeight: contentHeight,
		winSurface:    winSurface,
		font:          font,
	}, nil
}

// x, y are in pixels
func getColor(img image.Image, winWidth, winHeight, x, y int) sdl.Color {
	// coordinates relative to image
	imgX := x - (winWidth-img.Bounds().Max.X)/2
	imgY := y - (winHeight-img.Bounds().Max.Y)/2

	codeColor := sdl.Color(config.codeColor)

	// outside of image
	if imgX < 0 || imgX >= img.Bounds().Max.X {
		return codeColor
	}

	if imgY < 0 || imgY >= img.Bounds().Max.Y {
		return codeColor
	}

	color := img.At(imgX, imgY)
	r, g, b, a := color.RGBA()

	// full transparent
	if a == 0 {
		return codeColor
	}

	result := sdl.Color{
		R: uint8((float32)(r) / 0xffff * 0xff),
		G: uint8(float32(g) / 0xffff * 0xff),
		B: uint8(float32(b) / 0xffff * 0xff),
		A: uint8(float32(a) / 0xffff * 0xff),
	}

	if result == sdl.Color(config.bgColor) {
		result = codeColor
	}

	return result
}

func run() error {
	// init sdl
	sdlContext, err := initSDLAndTTF()
	if err != nil {
		return err
	}

	// read code
	code, err := readCode()
	if err != nil {
		return err
	}
	if len(code) == 0 {
		return errors.New("there is no valid characters in the source code (visible ascii characters)")
	}

	// open and resize image
	img, err := openAndResizeImage(int(sdlContext.contentWidth), int(sdlContext.contentHeight))
	if err != nil {
		return err
	}

	// poll sdl events
	for sdl.PollEvent() != nil {
	}

	// render
	sdlContext.renderer.SetDrawColor(config.bgColor.R, config.bgColor.G, config.bgColor.B, config.bgColor.A)
	sdlContext.renderer.Clear()
	var dstRect sdl.Rect
	for cy := 0; cy < config.height; cy++ {
		for cx := 0; cx < config.width; cx++ {
			index := (cy*config.width + cx) % len(code)
			char := code[index]
			x := (cx + config.padding.horizontal) * int(sdlContext.charWidth)
			y := (cy + config.padding.vertical) * int(sdlContext.charHeight)

			centerX := x + (int)(sdlContext.charWidth)/2
			centerY := y + (int)(sdlContext.charHeight)/2
			color := getColor(img, int(sdlContext.winWidth), int(sdlContext.winHeight), centerX, centerY)

			t, err := renderChar(sdlContext.font, sdlContext.renderer, char, color)

			if err != nil {
				return errors.Wrap(err, "could not render string in sdl ttf")
			}

			dstRect.X = int32(x)
			dstRect.Y = int32(y)
			dstRect.W = t.w
			dstRect.H = t.h

			if err := sdlContext.renderer.Copy(t.texture, nil, &dstRect); err != nil {
				return errors.Wrap(err, "sdl renderer failed")
			}
		}
	}

	// output
	sourceBase := path.Base(config.sourcePath)
	outputName := sourceBase + ".png"
	for i := 1; fileExists(outputName); i++ {
		outputName = fmt.Sprintf("%s.%d.png", sourceBase, i)
	}

	if err := sdlimg.SavePNG(sdlContext.winSurface, outputName); err != nil {
		return errors.Wrap(err, "could not save png of sdl surface")
	}

	log.Printf("code poster generated: %s\n", outputName)

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
