package ETEHelper

import (
	"image"
	_ "image/png"
	"os"
)

// Json

func LoadJson(jsonPath string) string {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// Image

func LoadImage(imagePath string) image.Image {
	data, err := os.OpenFile(imagePath, os.O_RDONLY, 0o644)
	if err != nil {
		panic(err)
	}
	defer data.Close()

	img, _, err := image.Decode(data)
	if err != nil {
		panic(err)
	}
	return img
}
