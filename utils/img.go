package utils

import (
	"golang.org/x/image/draw"
	"image"
)

func ResizeImage(src image.Image, dstSize image.Point) *image.RGBA {
	srcRect := src.Bounds()
	dstRect := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: dstSize,
	}
	dst := image.NewRGBA(dstRect)
	draw.CatmullRom.Scale(dst, dstRect, src, srcRect, draw.Over, nil)
	return dst
}

// CalculatedImageSize of the image based on the given max size
// if the image is smaller than the max size no calculation will be made
// otherwise the newly size will be determined keeping the aspect ratio
// of the original image. The calculated indicates if a calculation of the size has been made.
func CalculatedImageSize(src image.Image, maxSize image.Point) (size image.Rectangle, calculated bool) {
	srcSize := src.Bounds().Size()

	if srcSize.X <= maxSize.X && srcSize.Y <= maxSize.Y {
		return image.Rectangle{}, false
	}

	originWhite := float64(srcSize.X)
	originHeight := float64(srcSize.Y)

	scaleX, scaleY := calcScaleFactors(
		uint(maxSize.X),
		uint(maxSize.X*srcSize.Y/srcSize.X),
		originWhite,
		originHeight,
	)

	newWidth := originWhite / scaleX
	newHeight := originHeight / scaleY

	dstRect := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: int(newWidth), Y: int(newHeight)},
	}

	return dstRect, true
}

func calcScaleFactors(width, height uint, oldWidth, oldHeight float64) (scaleX, scaleY float64) {
	if width == 0 {
		if height == 0 {
			scaleX = 1.0
			scaleY = 1.0
		} else {
			scaleY = oldHeight / float64(height)
			scaleX = scaleY
		}
	} else {
		scaleX = oldWidth / float64(width)
		if height == 0 {
			scaleY = scaleX
		} else {
			scaleY = oldHeight / float64(height)
		}
	}
	return
}
