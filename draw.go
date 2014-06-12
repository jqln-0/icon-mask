package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"path"
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

func (m *maskDrawer) ComposeSized(overlay string, outFile string, size uint) {
	// Make sure the parent directory exists.
	createPath(outFile)

	log.Printf("composing %v\n", outFile)
	cmd := exec.Command("composite", "-gravity", "center", overlay, m.baseImage,
		"-resize", fmt.Sprintf("%v,%v", size, size), outFile)
	cmd.Run()
}

func (m *maskDrawer) ComposeSVG(overlay string, outFile string) {
	// TODO.
}

func (m *maskDrawer) CreateIcons(icon GTKIconProperties, outdir string) {
	// Draw one icon for each available size.
	for _, size := range icon.Sizes {
		sizeStr := fmt.Sprintf("%vx%v", size, size)
		outPath := path.Join(outdir, sizeStr, "apps", icon.Name+".png")
		filePath := m.theme.GetIcon(icon.Name, int(size))
		m.ComposeSized(filePath, outPath, size)
	}

	// Draw a scaled icon if possible.
	if icon.Scalable {
		outPath := path.Join(outdir, "scalable", "apps", icon.Name+".svg")
		// TODO: Replace crazy constant with a flag.
		filePath := m.theme.GetIcon(icon.Name, 1024)
		m.ComposeSVG(filePath, outPath)
	}
}
