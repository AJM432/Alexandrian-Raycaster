package main

import (
	"image/color"
	"math"
	"time"
  "fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 1200
	screenHeight = 1200

  mapWidth = 24
  mapHeight = 24

)

type Game struct {
	pressedKeys []ebiten.Key
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var (
	prevUpdateTime    = time.Now()

worldMap = [][]int{
  {1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1},
  {1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,0,0,0,0,0,2,2,2,2,2,0,0,0,0,3,0,3,0,3,0,0,0,1},
  {1,0,0,0,0,0,2,0,0,0,2,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,0,0,0,0,0,2,0,0,0,2,0,0,0,0,3,0,0,0,3,0,0,0,1},
  {1,0,0,0,0,0,2,0,0,0,2,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,0,0,0,0,0,2,2,0,2,2,0,0,0,0,3,0,3,0,3,0,0,0,1},
  {1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,4,4,4,4,4,4,4,4,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,4,0,4,0,0,0,0,4,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,4,0,0,0,0,5,0,4,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,4,0,4,0,0,0,0,4,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,4,0,4,4,4,4,4,4,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,4,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,4,4,4,4,4,4,4,4,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1},
  {1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1}}

  posX = 12.0
  posY = 12.0
  dirX = 1.0
  dirY = 0.0
  planeX = 0.0
  planeY = 0.66

  wallCoords = [][]int{}
  wallSide = [screenWidth]int{}

)


func (g *Game) Update() error {
	timeDelta := float64(time.Since(prevUpdateTime))/1e9
	prevUpdateTime = time.Now()
  
  for x:=0; x < screenWidth; x++{
    cameraX := 2.0 * float64(x) / float64(screenWidth) - 1 //x-coordinate in camera space
    rayDirX := dirX + planeX * cameraX
    rayDirY := dirY + planeY * cameraX


      //which box of the map we're in
    mapX := int(posX)
    mapY := int(posY)

      //length of ray from current position to next x or y-side
    sideDistX := 0.0
    sideDistY := 0.0

     //length of ray from one x or y-side to next x or y-side
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

      //what direction to step in x or y-direction (either +1 or -1)
    stepX := 0.0
    stepY:= 0.0

    hit := 0 //was there a wall hit?
    side := 0 //was a NS or a EW wall hit?



      //calculate step and initial sideDist
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

      //perform DDA
      for hit == 0 {
        //jump to next map square, either in x-direction, or in y-direction
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
        //Check if ray has hit a wall
        if worldMap[mapX][mapY] > 0{ 
          hit = 1
        }
      } 
      //Calculate distance projected on camera direction (Euclidean distance would give fisheye effect!)
      if side == 0 { 
        perpWallDist = (sideDistX - deltaDistX)
      }else {
        perpWallDist = (sideDistY - deltaDistY)
      }

      //Calculate height of line to draw on screen
    lineHeight := int((float64(screenHeight) / perpWallDist))

      //calculate lowest and highest pixel to fill in current stripe
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

  moveSpeed := (timeDelta) * 5.0; //the constant value is in squares/second
  rotSpeed := timeDelta * 3.0; //the constant value is in radians/second


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
      // if (x_index >= 0 || x_index <= 24) && (y_index >=0 || y_index <= 24){
      if worldMap[x_index][y_index] == 0 {posX -= dirX * moveSpeed}
      if worldMap[int(posX)][int(posY - dirY * moveSpeed)] == 0 {posY -= dirY * moveSpeed}
      // }
		case "ArrowUp":
      if worldMap[int(posX + dirX * moveSpeed)][int(posY)] == 0 { posX += dirX * moveSpeed}
      if worldMap[int(posX)][int(posY + dirY * moveSpeed)] == 0 {posY += dirY * moveSpeed}
		case "ArrowLeft":
      //both camera direction and camera plane must be rotated
      oldDirX := dirX
      dirX = dirX * math.Cos(-rotSpeed) - dirY * math.Sin(-rotSpeed)
      dirY = oldDirX * math.Sin(-rotSpeed) + dirY * math.Cos(-rotSpeed)
      oldPlaneX := planeX
      planeX = planeX * math.Cos(-rotSpeed) - planeY * math.Sin(-rotSpeed)
      planeY = oldPlaneX * math.Sin(-rotSpeed) + planeY * math.Cos(-rotSpeed)
		case "ArrowRight":
      //both camera direction and camera plane must be rotated
      oldDirX := dirX;
      dirX = dirX * math.Cos(rotSpeed) - dirY * math.Sin(rotSpeed);
      dirY = oldDirX * math.Sin(rotSpeed) + dirY * math.Cos(rotSpeed);
      oldPlaneX := planeX;
      planeX = planeX * math.Cos(rotSpeed) - planeY * math.Sin(rotSpeed);
      planeY = oldPlaneX * math.Sin(rotSpeed) + planeY * math.Cos(rotSpeed);
		}
	}

	return nil
}


func (g *Game) Draw(screen *ebiten.Image) {
  screen.Fill(color.RGBA{137, 196, 244, 255})


	whiteClr := color.RGBA{255, 255, 255, 255}
	grayClr := color.RGBA{200, 200, 200, 255}
	goldClr := color.RGBA{212,175,55, 255}

  for x:=0; x < screenWidth; x++ {
    x0 := float32(x)
    y0 := float32(wallCoords[x][0])
    x1 := x0+1
    y1 := float32(wallCoords[x][1])
    clr := whiteClr
    if wallSide[x] == 1 {
      clr = grayClr
    }
    vector.StrokeLine(screen, x0, y0, x1, y1, 1, clr, false)
    vector.StrokeLine(screen, x0, y1, x1, screenHeight, 1, goldClr, false)

  }
}
func main() {

  for i:=0; i < screenWidth; i++ {
    wallCoords = append(wallCoords, []int{})
  }

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("The Library of Alexandria")
	fmt.Print("The Library of Alexandria")

	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
