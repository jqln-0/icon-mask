package main

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"path"
	"sync"

	"github.com/nfnt/resize"
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
	baseImage image.Image
	theme     GTKTheme
	sizeMutex sync.RWMutex
	sizes     map[uint]image.Image
	interp    resize.InterpolationFunction
	scale     float64
}

func CreateMaskDrawer(base string, theme GTKTheme, scale float64, interp resize.InterpolationFunction) *maskDrawer {
	var m maskDrawer
	img, err := loadImage(base)
	if err != nil {
		log.Fatal(err)
	}
	m.baseImage = img
	m.theme = theme
	m.sizes = make(map[uint]image.Image)
	m.interp = interp
	m.scale = scale
	return &m
}

func createPath(outFile string) {
	err := os.MkdirAll(path.Dir(outFile), 0775)
	if err != nil {
		log.Fatal(err)
	}
}

func (m *maskDrawer) getScaled(size uint) image.Image {
	// Check if a base image of this size has been cached.
	m.sizeMutex.RLock()
	img, ok := m.sizes[size]
	m.sizeMutex.RUnlock()
	if ok {
		return img
	}

	// Not image was cached. We need to generate one.
	defer m.sizeMutex.Unlock()
	m.sizeMutex.Lock()

	// We need to check again incase the image was cached between our locks.
	img, ok = m.sizes[size]
	if ok {
		return img
	}

	// Create the scaled image.
	img = resize.Resize(size, size, m.baseImage, m.interp)
	m.sizes[size] = img
	return img
}

func (m *maskDrawer) ComposeSized(overlay string, outFile string, size uint) {
	// Make sure the parent directory exists.
	createPath(outFile)
	log.Printf("composing %v\n", outFile)

	// Get the resized base image.
	img := m.getScaled(size)

	// Load the overlay image.
	overlayImg, err := loadImage(overlay)
	if err != nil {
		log.Printf("Failed to compose %v: %v\n", outFile, err)
		return
	}

	// Reduce the size of the overlay according to the scale.
	scaledSize := uint(float64(size) * m.scale)
	overlayImg = resize.Thumbnail(scaledSize, scaledSize, overlayImg, m.interp)

	// Compose the images.
	result := image.NewRGBA(img.Bounds())
	draw.Draw(result, img.Bounds(), img, image.ZP, draw.Src)

	// The overlay image should be drawn with an offset as it may be scaled down.
	offset := (img.Bounds().Dx() - overlayImg.Bounds().Dx()) / 2
	destPoint := image.Pt(offset, offset)
	r := image.Rectangle{destPoint, destPoint.Add(overlayImg.Bounds().Size())}
	draw.Draw(result, r, overlayImg, overlayImg.Bounds().Min, draw.Over)

	// Write out the result.
	fd, err := os.Create(outFile)
	if err != nil {
		log.Print(err)
		return
	}
	defer fd.Close()

	// TODO: Check using a buffered reader helps performance here.
	writer := bufio.NewWriter(fd)

	png.Encode(writer, result)
	writer.Flush()
}

func (m *maskDrawer) ComposeSVG(overlay string, outFile string) {
	// TODO.
	log.Printf("Failed to compose %v: %v\n", outFile, "SVG not supported")
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
