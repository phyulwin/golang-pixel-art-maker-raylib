# R-artBit🐰 Golang Pixel Art Maker

A lightweight pixel art editor built with Go using [raylib-go](https://github.com/gen2brain/raylib-go). No external UI dependencies. Designed for quick sketching, color experimentation with basic single-layer canvas.

![Pixel Art Editor Empty Canvas](Screenshots/Screenshot%202025-07-30%20204033.png)
---

## Features
- **Interactive Canvas**  
    - 32×32 grid, real-time painting
    - Pan with Shift + Right Click
    - Zoom with mouse wheel

- **Color Selection**  
    - Fixed color palette  
    - RGB sliders for custom colors  
    - Live preview and selection

- **Tools**  
    - Undo / Redo  
    - Toggle Grid  
    - New Canvas (clear all)  
    - Save to PNG

### Live Drawing in Canvas View:

![Pixel Art Example 1](Screenshots/Screenshot%202025-07-30%20204427.png)

### Artwork is exported as PNG file:

![Pixel Art Example 2](Screenshots/Screenshot%202025-07-30%20204510.png)
---

| Action                  | Input                      |
|-------------------------|----------------------------|
| Draw                   | Left Click                 |
| Erase                  | Paint with blank color     |
| Pan                   | Shift + Right Click & Drag |
| Zoom                   | Mouse Wheel                |
| Color Pick (Custom)    | Adjust RGB sliders + "Pick" button |

---

## Build & Run

```go get github.com/gen2brain/raylib-go/raylib```

```go run .```
