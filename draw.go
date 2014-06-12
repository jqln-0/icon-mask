package main

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
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
	baseImage    image.Image
	scaledImages map[uint]image.Image
	scaledMutex  sync.RWMutex
	interp       resize.InterpolationFunction
}

func CreateMaskDrawer(base string, interp resize.InterpolationFunction) *maskDrawer {
	var m maskDrawer
	img, err := loadImage(base)
	if err != nil {
		log.Fatalf("Failed to load base image: %v\n", err)
	}
	m.baseImage = img
	m.scaledImages = make(map[uint]image.Image)
	m.interp = interp
	return &m
}

func (m *maskDrawer) getScaled(size uint) image.Image {
	m.scaledMutex.RLock()
	img, ok := m.scaledImages[size]
	m.scaledMutex.RUnlock()
	if ok {
		return img
	}

	// We haven't seen this size before; scale and cache the new image.
	defer m.scaledMutex.Unlock()
	m.scaledMutex.Lock()
	img = resize.Resize(size, size, m.baseImage, m.interp)
	m.scaledImages[size] = img
	return img
}

func (m *maskDrawer) CreateMask(icon string, size int) (io.Reader, error) {
	// TODO: All this loading and resizing could be done concurrently.
	img, err := loadImage(icon)
	if err != nil {
		return nil, err
	}

	// Make sure the icon image is the correct size.
	if img.Bounds().Dx() != size || img.Bounds().Dy() != size {
		// TODO: Make sure uint conversion doesn't overflow.
		img = resize.Resize(uint(size), uint(size), img, m.interp)
	}
	baseImg := m.getScaled(uint(size))

	// Perform the composition.
	final := image.NewRGBA(baseImg.Bounds())
	draw.Draw(final, baseImg.Bounds(), baseImg, image.ZP, draw.Src)
	draw.Draw(final, img.Bounds(), img, image.ZP, draw.Src)

	// Create and return the output reader.
	var buf bytes.Buffer
	png.Encode(&buf, final)
	return &buf, nil
}
