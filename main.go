package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/hajimehoshi/ebiten/v2/text"
	"io/ioutil"
	"log"
	"strings"
  "os"
  "path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 1500
	screenHeight = 750
  FONT_SIZE = 11

  mapWidth = 24
  mapHeight = 24

	textViewWidth = int(500)
	textViewHeight = screenHeight
	textViewX     = int(screenWidth - textViewWidth)
	textViewY     = 0
	scrollSpeed   = 3

)
type Game struct {
	textLines []string
	scroll    int // Tracks the scrolling position
	pressedKeys []ebiten.Key
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var (
	prevUpdateTime    = time.Now()
  timeDelta    = 0.0

worldMap = [][]int{
  {0,4,4,4,4,4,4,4,4,4,4,4,4,4,4,4,7,7,7,7,7,7,7,0},
  {4,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,7,0,0,0,0,0,0,7},
  {4,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,7},
  {4,0,2,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,7},
  {4,0,3,0,0,0,0,0,0,0,0,0,0,0,0,0,7,0,0,0,0,0,0,7},
  {4,0,4,0,0,0,0,5,5,5,5,5,5,5,5,5,7,7,0,7,7,7,7,7},
  {4,0,5,0,0,0,0,5,0,5,0,5,0,5,0,5,7,0,0,0,7,7,7,1},
  {4,0,6,0,0,0,0,5,0,0,0,0,0,0,0,5,7,0,0,0,0,0,0,8},
  {4,0,7,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,7,7,7,1},
  {4,0,8,0,0,0,0,5,0,0,0,0,0,0,0,5,7,0,0,0,0,0,0,8},
  {4,0,0,0,0,0,0,5,0,0,0,0,0,0,0,5,7,0,0,0,7,7,7,1},
  {4,0,0,0,0,0,0,5,5,5,5,0,5,5,5,5,7,7,7,7,7,7,7,1},
  {6,6,6,6,6,6,6,6,6,6,6,0,6,6,6,6,6,6,6,6,6,6,6,6},
  {8,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,4},
  {6,6,6,6,6,6,0,6,6,6,6,0,6,6,6,6,6,6,6,6,6,6,6,6},
  {4,4,4,4,4,4,0,4,4,4,6,0,6,2,2,2,2,2,2,2,3,3,3,3},
  {4,0,0,0,0,0,0,0,0,4,6,0,6,2,0,0,0,0,0,2,0,0,0,2},
  {4,0,0,0,0,0,0,0,0,0,0,0,6,2,0,0,5,0,0,2,0,0,0,2},
  {4,0,0,0,0,0,0,0,0,4,6,0,6,2,0,0,0,0,0,2,2,0,2,2},
  {4,0,6,0,6,0,0,0,0,4,6,0,0,0,0,0,5,0,0,0,0,0,0,2},
  {4,0,0,5,0,0,0,0,0,4,6,0,6,2,0,0,0,0,0,2,2,0,2,2},
  {4,0,6,0,6,0,0,0,0,4,6,0,6,2,0,0,5,0,0,2,0,0,0,2},
  {4,0,0,0,0,0,0,0,0,4,6,0,6,2,0,0,0,0,0,2,0,0,0,2},
  {0,4,4,4,4,4,4,4,4,4,1,1,1,2,2,2,2,2,2,3,3,3,3,0}}

  posX = 5.0
  posY = 5.0
  dirX = 1.0
  dirY = 0.0
  planeX = 0.0
  planeY = 0.66

  wallCoords = [][]int{}
  wallSide = [screenWidth]int{}
  perpDists = [screenWidth]float64{}

  currentViewedBlockX = -1
  currentViewedBlockY = -1
  currentViewedBlockStart = -1
  currentViewedBlockEnd = -1

  book_names = []string{}
  curr_book_scroll = []int{}
  currentBookIndex = -1

)
var textFont font.Face = LoadFont("fonts/CinzelDecorative-Regular.ttf", 11)


func (g *Game) Update() error {
  currentViewedBlockStart = -1
	timeDelta = float64(time.Since(prevUpdateTime))/1e9
	prevUpdateTime = time.Now()
  
  for x:=0; x < screenWidth; x++{ // loop through each strip of screenWidth
    cameraX := 2.0 * float64(x) / float64(screenWidth) - 1 // in range [-1, 1] left to right
    rayDirX := dirX + planeX * cameraX
    rayDirY := dirY + planeY * cameraX


      //current box
    mapX := int(posX)
    mapY := int(posY)

      //length to next X or Y box from my pos
    sideDistX := 0.0
    sideDistY := 0.0

     //gen length to next X, Y
    deltaDistX := 0.0
    deltaDistY := 0.0

    if rayDirX == 0{
      deltaDistX = math.Inf(1)
    }  else{
      deltaDistX = math.Abs(1.0 / float64(rayDirX))
    }

    if rayDirY == 0{
      deltaDistY = math.Inf(1)
    }  else{
      deltaDistY = math.Abs(1.0 / float64(rayDirY))
    }
    perpWallDist := 0.0

    stepX := 0.0
    stepY:= 0.0

    hit := 0
    side := 0



      if rayDirX < 0{
        stepX = -1
        sideDistX = (posX - float64(mapX)) * deltaDistX
      } else
      {
        stepX = 1
        sideDistX = (float64(mapX) + 1.0 - posX) * deltaDistX
      }
      if rayDirY < 0 {
        stepY = -1
        sideDistY = (posY - float64(mapY)) * deltaDistY
      } else
      {
        stepY = 1
        sideDistY = (float64(mapY) + 1.0 - posY) * deltaDistY;
      }

      for hit == 0 {
        if sideDistX < sideDistY {
          sideDistX += deltaDistX
          mapX += int(stepX)
          side = 0
        } else
        {
          sideDistY += deltaDistY;
          mapY += int(stepY)
          side = 1
        }
        if worldMap[mapX][mapY] > 0{ 
          if (rayDirX == dirX) && (rayDirY == dirY) {
            currentViewedBlockX = mapX
            currentViewedBlockY = mapY
          }
          if (mapX == currentViewedBlockX) && (mapY == currentViewedBlockY){
            if (currentViewedBlockStart == -1){
              currentViewedBlockStart = x
            }else{
              currentViewedBlockEnd = x
            }

        }
            
          hit = 1
        }
      } 
      if side == 0 { 
        perpWallDist = (sideDistX - deltaDistX)
      }else {
        perpWallDist = (sideDistY - deltaDistY)
      }
    perpDists[x] = perpWallDist

    lineHeight := int((float64(screenHeight) / perpWallDist))

    drawStart := -lineHeight/2 + screenHeight/2

      if drawStart < 0 {
        drawStart = 0
      }

    drawEnd := lineHeight / 2 + screenHeight/ 2
      if drawEnd >= screenHeight {
        drawEnd = screenHeight - 1
      }
    wallStrip :=[]int{drawStart, drawEnd}
    wallCoords[x] = wallStrip
    wallSide[x] = side

  }

  moveSpeed := (timeDelta) * 2.0;
  rotSpeed := timeDelta * 2.0;


	g.pressedKeys = inpututil.AppendPressedKeys(g.pressedKeys[:0])

	for _, key := range g.pressedKeys {
		switch key.String() {
		case "Space":
		}
	}

	for _, key := range g.pressedKeys {
		switch key.String() {
		case "ArrowDown":
      x_index := int(posX - dirX * moveSpeed)
      y_index := int(posY)
      if worldMap[x_index][y_index] == 0 {posX -= dirX * moveSpeed}
      if worldMap[int(posX)][int(posY - dirY * moveSpeed)] == 0 {posY -= dirY * moveSpeed}
		case "ArrowUp":
      if worldMap[int(posX + dirX * moveSpeed)][int(posY)] == 0 { posX += dirX * moveSpeed}
      if worldMap[int(posX)][int(posY + dirY * moveSpeed)] == 0 {posY += dirY * moveSpeed}
		case "ArrowLeft":
      oldDirX := dirX
      dirX = dirX * math.Cos(-rotSpeed) - dirY * math.Sin(-rotSpeed)
      dirY = oldDirX * math.Sin(-rotSpeed) + dirY * math.Cos(-rotSpeed)
      oldPlaneX := planeX
      planeX = planeX * math.Cos(-rotSpeed) - planeY * math.Sin(-rotSpeed)
      planeY = oldPlaneX * math.Sin(-rotSpeed) + planeY * math.Cos(-rotSpeed)
		case "ArrowRight":
      oldDirX := dirX;
      dirX = dirX * math.Cos(rotSpeed) - dirY * math.Sin(rotSpeed);
      dirY = oldDirX * math.Sin(rotSpeed) + dirY * math.Cos(rotSpeed);
      oldPlaneX := planeX;
      planeX = planeX * math.Cos(rotSpeed) - planeY * math.Sin(rotSpeed);
      planeY = oldPlaneX * math.Sin(rotSpeed) + planeY * math.Cos(rotSpeed);
		}
	}

  handleScrolling(g)

	return nil
}

func (g *Game) DrawMiniMap(screen *ebiten.Image) {
  var alpha uint8 = 175
  whiteClr := color.RGBA{255, 255, 255, alpha} 
  blackClr := color.RGBA{0, 0, 0, alpha} 
  greenClr := color.RGBA{0, 255, 0, alpha} 
  focusClr := color.RGBA{196, 164, 132, alpha} 
	yellowClr := color.RGBA{255,255,0, 255}
  miniMapSize := 200
  blockSize := miniMapSize / mapWidth
  blockDrawClr := whiteClr

  vector.DrawFilledRect(screen, 0, 0, float32(miniMapSize), float32(miniMapSize), blackClr, false)

  for row := range mapWidth {
    for col := range mapHeight {
      val := worldMap[row][col]
      if val > 0 {
        if (row == currentViewedBlockX) && (col == currentViewedBlockY){
          blockDrawClr = focusClr
        }else {
          blockDrawClr = whiteClr
        }
        x0 := blockSize*row
        y0 := blockSize*col
          vector.DrawFilledRect(screen, float32(x0), float32(y0), float32(blockSize), float32(blockSize), blockDrawClr, false)
        }
      }
    }

  playerRadius := 5
  relPlayerX := float32(posX*float64(blockSize))
  relPlayerY := float32(posY*float64(blockSize))
  vector.DrawFilledCircle(screen, relPlayerX, relPlayerY, float32(playerRadius), greenClr, false)

  for x:=0; x < screenWidth; x++{ // loop through each strip of screenWidth
    cameraX := 2.0 * float64(x) / float64(screenWidth) - 1 // in range [-1, 1] left to right
    rayDirX := float32(dirX + planeX * cameraX)
    rayDirY := float32(dirY + planeY * cameraX)
    MAP_SCALE_FACTOR := float32(8.0)
    vector.StrokeLine(screen, relPlayerX, relPlayerY, relPlayerX + MAP_SCALE_FACTOR*float32(perpDists[x])*rayDirX*float32(1),relPlayerY + MAP_SCALE_FACTOR*float32(perpDists[x])*rayDirY*float32(1), 1, yellowClr, false)
  }
}


func (g *Game) Draw(screen *ebiten.Image) {
  skyBlueClr := color.RGBA{137, 196, 244, 255} 
  screen.Fill(skyBlueClr)


	whiteClr := color.RGBA{255, 255, 255, 255}
	grayClr := color.RGBA{200, 200, 200, 255}
	goldClr := color.RGBA{212,175,55, 255}
  focusClr := color.RGBA{196, 164, 132, 255} 

  for x:=0; x < screenWidth; x++ {
    x0 := float32(x)
    y0 := float32(wallCoords[x][0])
    x1 := x0+1
    y1 := float32(wallCoords[x][1])
    clr := whiteClr
    if (x >= currentViewedBlockStart) && (x <= currentViewedBlockEnd){
      clr = focusClr
    }else if wallSide[x] == 1 {
      clr = grayClr
    }else{
      clr = whiteClr
    }
    vector.StrokeLine(screen, x0, y0, x1, y1, 1, clr, false)
    if y1 < screenHeight{
      vector.StrokeLine(screen, x0, y1, x1, screenHeight, 1, goldClr, false)
    }
  }
  
  g.DrawMiniMap(screen)
  // FPS := 1/timeDelta

  bookIndex := getBookIndex(currentViewedBlockX, currentViewedBlockY, mapWidth, mapHeight)
  if len(book_names) > bookIndex {
    text.Draw(screen, book_names[bookIndex], textFont, screenWidth/2, 15, color.Black)
    if currentBookIndex != bookIndex {
      if currentBookIndex >= 0{
        curr_book_scroll[currentBookIndex] = g.scroll
      }

      lines := loadTextFile(book_names[bookIndex])
      currentBookIndex = bookIndex
      g.textLines = lines
      g.scroll = curr_book_scroll[currentBookIndex]
    }
    renderText(g, screen, textViewX, textViewY, textViewWidth, textViewHeight)
  }
}

func getBookIndex(mapX int, mapY int, mapWidth int, mapHeight int) int {
  apparentIndex := (mapY*mapWidth + mapX)

  dist_to_book := (perpDists[int(len(perpDists)/2)])
  if dist_to_book > 2 {
    return len(book_names) + 1
  }
  return apparentIndex % (len(book_names)-1)
}

// NewGame initializes the game with the text file
func NewGame(filePath string) *Game {
	lines := loadTextFile(filePath)
	return &Game{
		textLines: lines,
		scroll:    0,
	}
}

func handleScrolling(g *Game) {
	_, wheelY := ebiten.Wheel()
	if wheelY != 0 {
		g.scroll -= int(wheelY) * scrollSpeed
		if g.scroll < 0 {
			g.scroll = 0
		}
		if g.scroll > len(g.textLines)-1 {
			g.scroll = len(g.textLines) - 1
		}
	}
}

func renderText(g *Game, screen *ebiten.Image, x, y, width, height int) {
	fontFace := textFont
	lineHeight := fontFace.Metrics().Height.Ceil()
	maxLines := height / lineHeight
	startLine := g.scroll
	endLine := startLine + maxLines
	if endLine > len(g.textLines) {
		endLine = len(g.textLines)
	}

	rectangle := ebiten.NewImage(width, height)
  papyrusClr := color.RGBA{201, 152, 104, 255}
	rectangle.Fill(papyrusClr) // Text viewer background color
  op := &ebiten.DrawImageOptions{}
  op.GeoM.Translate(float64(x), float64(y))

  screen.DrawImage(rectangle, op)

	for i, line := range g.textLines[startLine:endLine] {
		textX := x + 10
		textY := y + (i * lineHeight) + lineHeight
		if textY > y+height {
			break
		}
		text.Draw(screen, line, textFont, textX, textY, color.Black)
	}
}

func LoadFont(filePath string, size float64) font.Face {
	fontBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read font file: %v", err)
	}
	tt, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("Failed to parse font: %v", err)
	}
	fontFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("Failed to create font face: %v", err)
	}
	return fontFace
}

func loadTextFile(filePath string) []string {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	return strings.Split(string(content), "\n")
}

func getBookNames(dir string) []string {
    entries, err := os.ReadDir(dir)
    if err != nil {
       log.Fatal(err)
    }

    var files []string
    for _, v := range entries {
       if v.IsDir() {
          continue
       }
       files = append(files, filepath.Join(dir, v.Name()))
       curr_book_scroll = append(curr_book_scroll, 0)
    }

    return files
  }

func main() {

  book_names = getBookNames("books")

  for i:=0; i < screenWidth; i++ {
    wallCoords = append(wallCoords, []int{})
  }

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("The Library of Alexandria")
	fmt.Println("The Library of Alexandria")

  game := NewGame(book_names[0])
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}

}
