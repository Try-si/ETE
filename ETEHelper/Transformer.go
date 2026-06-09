package ETEHelper

import (
	"encoding/json"
	"image"
	"image/draw"
)

// Json

func JsonToStruct[T any](jsonPath string) T {
	return StringToStruct[T](LoadJson(jsonPath))
}

func StringToStruct[T any](jsonString string) T {
	var result T
	json.Unmarshal([]byte(jsonString), &result)
	return result
}

func StructToJson[T any](data T) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

// img

func SliceImageByGrid(img image.Image, cellInPx int) []image.Image {
	bounds := img.Bounds()
	cols := bounds.Max.X / cellInPx
	rows := bounds.Max.Y / cellInPx

	result := make([]image.Image, 0, cols*rows)

	for iy := 0; iy < rows; iy++ {
		for ix := 0; ix < cols; ix++ {
			rect := image.Rect(
				ix*cellInPx, iy*cellInPx,
				(ix+1)*cellInPx, (iy+1)*cellInPx,
			)

			subImg := image.NewRGBA(rect)
			draw.Draw(subImg, rect, img, rect.Min, draw.Src)
			result = append(result, subImg)
		}
	}
	return result
}
