package utils

import (
	"image"
	"testing"
)

func TestResizeImage(t *testing.T) {
	tests := []struct {
		name     string
		src      *image.RGBA
		dstSize  image.Point
		expected image.Point
	}{
		{
			name:     "Test Resize Image",
			src:      image.NewRGBA(image.Rect(0, 0, 10, 10)),
			dstSize:  image.Point{X: 20, Y: 20},
			expected: image.Point{X: 20, Y: 20},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				resized := ResizeImage(tt.src, tt.dstSize)
				if resized.Bounds().Max != tt.expected {
					t.Errorf("Expected size: %v, but got: %v", tt.expected, resized.Bounds().Max)
				}
			},
		)
	}
}

func TestCalculatedImageSize(t *testing.T) {
	tests := []struct {
		name     string
		src      *image.RGBA
		maxSize  image.Point
		expected image.Point
	}{
		{
			name:     "Test Calculated Image Size",
			src:      image.NewRGBA(image.Rect(0, 0, 10, 10)),
			maxSize:  image.Point{X: 5, Y: 5},
			expected: image.Point{X: 5, Y: 5},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				newSize, calculated := CalculatedImageSize(tt.src, tt.maxSize)
				if calculated {
					if newSize.Max != tt.expected {
						t.Errorf("Expected size: %v, but got: %v", tt.expected, newSize.Max)
					}
				} else {
					if newSize.Max != tt.src.Bounds().Max {
						t.Errorf("Expected original size: %v, but got: %v", tt.src.Bounds().Max, newSize.Max)
					}
				}
			},
		)
	}
}

func TestCalcScaleFactors(t *testing.T) {
	tests := []struct {
		name                 string
		width, height        uint
		oldWidth, oldHeight  float64
		expectedX, expectedY float64
	}{
		{
			name:      "Test Calc Scale Factors",
			width:     10,
			height:    10,
			oldWidth:  20,
			oldHeight: 20,
			expectedX: 2.0,
			expectedY: 2.0,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				scaleX, scaleY := calcScaleFactors(tt.width, tt.height, tt.oldWidth, tt.oldHeight)
				if scaleX != tt.expectedX || scaleY != tt.expectedY {
					t.Errorf(
						"Expected scale factors: %v, %v, but got: %v, %v",
						tt.expectedX,
						tt.expectedY,
						scaleX,
						scaleY,
					)
				}
			},
		)
	}
}
