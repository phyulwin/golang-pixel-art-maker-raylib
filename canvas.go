package main

import (
    rl "github.com/gen2brain/raylib-go/raylib"
)

func SaveCanvasAsPNG(canvas [gridSize][gridSize]rl.Color, filename string) {
    img := rl.GenImageColor(gridSize, gridSize, rl.Blank)
    for y := 0; y < gridSize; y++ {
        for x := 0; x < gridSize; x++ {
            rl.ImageDrawPixel(img, int32(x), int32(y), canvas[y][x])
        }
    }
    rl.ExportImage(*img, filename)
    rl.UnloadImage(img)
}

func CaptureCanvas(canvas [gridSize][gridSize]rl.Color) Action {
    var pixels []Pixel
    for y := 0; y < gridSize; y++ {
        for x := 0; x < gridSize; x++ {
            pixels = append(pixels, Pixel{
                X:     x,
                Y:     y,
                Color: canvas[y][x],
            })
        }
    }
    return Action{Pixels: pixels}
}

func ApplyAction(canvas *[gridSize][gridSize]rl.Color, action Action) {
    for _, p := range action.Pixels {
        canvas[p.Y][p.X] = p.Color
    }
}