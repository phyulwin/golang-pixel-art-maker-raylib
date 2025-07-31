package main

import (
    "time"
    rl "github.com/gen2brain/raylib-go/raylib"
)

const (
    screenWidth  = 1024
    screenHeight = 768
    gridSize     = 32
    cellSize     = 20
    canvasOffset = 160
)

type Pixel struct {
    X, Y  int
    Color rl.Color
}

type Action struct {
    Pixels []Pixel
}

func App() {
    rl.InitWindow(screenWidth, screenHeight, "R-artBit🐰")
    rl.SetTargetFPS(60)

    var canvas [gridSize][gridSize]rl.Color
    var currentColor = rl.Black
    var showGrid = true
    var mouseDown bool
    var panX, panY float32
    var zoom float32 = 1.0

    var undoStack []Action
    var redoStack []Action

    palette := []rl.Color{rl.Black, rl.Red, rl.Green, rl.Blue, rl.Yellow, rl.Purple, rl.Orange, rl.White}
    paletteRects := make([]rl.Rectangle, len(palette))
    for i := range palette {
        paletteRects[i] = rl.NewRectangle(20, float32(20+i*30), 30, 30)
    }

    var red, green, blue uint8 = 0, 0, 0

    saveBtn := rl.NewRectangle(20, 280, 100, 30)
    undoBtn := rl.NewRectangle(20, 320, 100, 30)
    redoBtn := rl.NewRectangle(20, 360, 100, 30)
    gridChk := rl.NewRectangle(20, 400, 20, 20)
    newBtn := rl.NewRectangle(20, 440, 100, 30)

    for !rl.WindowShouldClose() {
        rl.BeginDrawing()
        rl.ClearBackground(rl.RayWhite)

        mouse := rl.GetMousePosition()

        rl.DrawText("Color", 20, 0, 20, rl.DarkGray)
        for i, rect := range paletteRects {
            rl.DrawRectangleRec(rect, palette[i])
            if rl.CheckCollisionPointRec(mouse, rect) && rl.IsMouseButtonPressed(0) {
                currentColor = palette[i]
            }
            if currentColor == palette[i] {
                rl.DrawRectangleLinesEx(rect, 2, rl.DarkGray)
            }
        }

        sliderX := float32(20)
        sliderW := float32(100)

        drawSlider := func(label string, y float32, val *uint8, color rl.Color) {
            rl.DrawText(label, int32(sliderX), int32(y), 10, rl.Black)
            rl.DrawRectangle(int32(sliderX), int32(y+12), int32(sliderW), 8, rl.LightGray)
            barWidth := float32(*val) / 255 * sliderW
            rl.DrawRectangle(int32(sliderX), int32(y+12), int32(barWidth), 8, color)

            if rl.CheckCollisionPointRec(mouse, rl.NewRectangle(sliderX, y+12, sliderW, 8)) &&
                rl.IsMouseButtonDown(0) {
                rel := mouse.X - sliderX
                if rel < 0 {
                    rel = 0
                } else if rel > sliderW {
                    rel = sliderW
                }
                *val = uint8((rel / sliderW) * 255)
            }
        }

        drawSlider("R", 520, &red, rl.Red)
        drawSlider("G", 550, &green, rl.Green)
        drawSlider("B", 580, &blue, rl.Blue)

        preview := rl.NewRectangle(20, 620, 30, 30)
        rl.DrawRectangleRec(preview, rl.NewColor(red, green, blue, 255))
        rl.DrawRectangleLinesEx(preview, 1, rl.Black)
        rl.DrawText("Pick", 60, 626, 14, rl.DarkGray)

        if rl.CheckCollisionPointRec(mouse, preview) && rl.IsMouseButtonPressed(0) {
            currentColor = rl.NewColor(red, green, blue, 255)
        }

        rl.DrawRectangleRec(saveBtn, rl.LightGray)
        rl.DrawText("Save", int32(saveBtn.X)+10, int32(saveBtn.Y)+8, 20, rl.Black)
        if rl.CheckCollisionPointRec(mouse, saveBtn) && rl.IsMouseButtonPressed(0) {
            SaveCanvasAsPNG(canvas, "pixel_art_"+time.Now().Format("150405")+".png")
        }

        rl.DrawRectangleRec(undoBtn, rl.LightGray)
        rl.DrawText("Undo", int32(undoBtn.X)+10, int32(undoBtn.Y)+8, 20, rl.Black)
        if rl.CheckCollisionPointRec(mouse, undoBtn) && rl.IsMouseButtonPressed(0) && len(undoStack) > 0 {
            last := undoStack[len(undoStack)-1]
            undoStack = undoStack[:len(undoStack)-1]
            redoStack = append(redoStack, CaptureCanvas(canvas))
            ApplyAction(&canvas, last)
        }

        rl.DrawRectangleRec(redoBtn, rl.LightGray)
        rl.DrawText("Redo", int32(redoBtn.X)+10, int32(redoBtn.Y)+8, 20, rl.Black)
        if rl.CheckCollisionPointRec(mouse, redoBtn) && rl.IsMouseButtonPressed(0) && len(redoStack) > 0 {
            last := redoStack[len(redoStack)-1]
            redoStack = redoStack[:len(redoStack)-1]
            undoStack = append(undoStack, CaptureCanvas(canvas))
            ApplyAction(&canvas, last)
        }

        rl.DrawText("Grid", 50, int32(gridChk.Y)-3, 20, rl.DarkGray)
        rl.DrawRectangleLinesEx(gridChk, 1, rl.DarkGray)
        if showGrid {
            rl.DrawLine(int32(gridChk.X), int32(gridChk.Y), int32(gridChk.X+20), int32(gridChk.Y+20), rl.DarkGray)
            rl.DrawLine(int32(gridChk.X+20), int32(gridChk.Y), int32(gridChk.X), int32(gridChk.Y+20), rl.DarkGray)
        }
        if rl.CheckCollisionPointRec(mouse, gridChk) && rl.IsMouseButtonPressed(0) {
            showGrid = !showGrid
        }

        // Draw "New" button
        rl.DrawRectangleRec(newBtn, rl.LightGray)
        rl.DrawText("New", int32(newBtn.X)+10, int32(newBtn.Y)+8, 20, rl.Black)
        if rl.CheckCollisionPointRec(mouse, newBtn) && rl.IsMouseButtonPressed(0) {
            for y := 0; y < gridSize; y++ {
                for x := 0; x < gridSize; x++ {
                    canvas[y][x] = rl.Blank
                }
            }
            undoStack = nil
            redoStack = nil
        }

        if rl.IsMouseButtonDown(1) && rl.IsKeyDown(rl.KeyLeftShift) {
            delta := rl.GetMouseDelta()
            panX += delta.X
            panY += delta.Y
        }

        wheel := rl.GetMouseWheelMove()
        if wheel != 0 {
            zoom += wheel * 0.1
            if zoom < 0.2 {
                zoom = 0.2
            } else if zoom > 5 {
                zoom = 5
            }
        }

        canvasX := int((mouse.X - canvasOffset - panX) / (cellSize * zoom))
        canvasY := int((mouse.Y - panY) / (cellSize * zoom))

        if rl.IsMouseButtonDown(0) &&
            mouse.X > canvasOffset && canvasX >= 0 && canvasX < gridSize && canvasY >= 0 && canvasY < gridSize {
            if !mouseDown {
                undoStack = append(undoStack, CaptureCanvas(canvas))
                redoStack = nil
                mouseDown = true
            }
            canvas[canvasY][canvasX] = currentColor
        } else {
            mouseDown = false
        }

        rl.BeginScissorMode(canvasOffset, 0, screenWidth-canvasOffset, screenHeight)
        for y := 0; y < gridSize; y++ {
            for x := 0; x < gridSize; x++ {
                rect := rl.NewRectangle(
                    canvasOffset+float32(x)*cellSize*zoom+panX,
                    float32(y)*cellSize*zoom+panY,
                    cellSize*zoom, cellSize*zoom,
                )
                rl.DrawRectangleRec(rect, canvas[y][x])
                if showGrid {
                    rl.DrawRectangleLinesEx(rect, 1, rl.LightGray)
                }
            }
        }
        rl.EndScissorMode()

        rl.EndDrawing()
    }

    rl.CloseWindow()
}