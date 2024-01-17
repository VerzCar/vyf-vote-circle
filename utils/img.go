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

// CalculatedImageSize calculates the new size of an image based on the maximum size specified
// If the image is smaller or equal to the maximum size, it returns an empty rectangle and false
// Otherwise, it calculates the scale factors using the calcScaleFactors function and then calculates the new width and height
// Finally, it returns a rectangle with the new size and true
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

// calcScaleFactors calculates the scale factors for resizing an image
// based on the desired width, height, and the original width and height
// If the desired width is 0, the scaleX factor is set to 1.0.
// If the desired height is 0, the scaleY factor is set to 1.0.
// If both the desired width and height are 0, the scaleX and scaleY factors are both set to 1.0.
// If both the desired width and height are non-zero, the scaleX factor is calculated as the ratio of
// the original width to the desired width, and the scaleY factor is calculated as the ratio of
// the original height to the desired height.
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
