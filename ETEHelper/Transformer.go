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

// transforme une list en une grille de taille W x H avec la formule i = y * W + x
func ListToGridYWX[T any](list []T, W, H int, center [2]int) map[[2]int]T {
	result := make(map[[2]int]T)
	for y := center[1] - H; y < H-center[1]; y++ {
		for x := center[0] - W; x < W-center[0]; x++ {
			result[[2]int{x, y}] = list[y*W+x]
		}
	}
	return result
}
