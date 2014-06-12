package main

import (
	"fmt"
	"image"
	"log"
	"path"

	"os"
)

func loadImage(filename string) (image.Image, error) {
	// Open the image file for reading.
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	img, _, err := image.Decode(fd)
	if err != nil {
		return nil, err
	}

	return img, nil
}

type maskDrawer struct {
	baseImage string
	theme     GTKTheme
}

func CreateMaskDrawer(base string, theme GTKTheme) *maskDrawer {
	var m maskDrawer
	m.baseImage = base
	m.theme = theme
	return &m
}

func createPath(outFile string) {
	err := os.MkdirAll(path.Dir(outFile), 0775)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *maskDrawer) ComposeSized(overlay string, outFile string) {
	// Make sure the parent directory exists.
	createPath(outFile)

	fmt.Printf("composing %v\n", outFile)
}

func (m *maskDrawer) ComposeSVG(overlay string, outFile string) {
	// Make sure the parent directory exists.
	createPath(outFile)

	fmt.Printf("composing %v\n", outFile)
}

func (m *maskDrawer) CreateIcons(icon GTKIconProperties, outdir string) {
	// Draw one icon for each available size.
	for _, size := range icon.Sizes {
		sizeStr := fmt.Sprintf("%vx%v", size, size)
		outPath := path.Join(outdir, sizeStr, "masked", icon.Name+".png")
		filePath := m.theme.GetIcon(icon.Name, int(size))
		m.ComposeSized(filePath, outPath)
	}

	// Draw a scaled icon if possible.
	if icon.Scalable {
		outPath := path.Join(outdir, "scalable", "masked", icon.Name+".svg")
		// TODO: Replace crazy constant with a flag.
		filePath := m.theme.GetIcon(icon.Name, 1024)
		m.ComposeSVG(filePath, outPath)
	}
}
