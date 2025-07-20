package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
	paddleWidth  = 100
	paddleHeight = 20
	ballSize     = 8
	blockWidth   = 50
	blockHeight  = 20
	blockRows    = 6
	blockCols    = 12
	paddleSpeed  = 5
	ballSpeed    = 3
)

// ブロック構造体
type Block struct {
	x, y    float64
	width   float64
	height  float64
	visible bool
	color   color.Color
}

// パドル構造体
type Paddle struct {
	x, y   float64
	width  float64
	height float64
}

// ボール構造体
type Ball struct {
	x, y   float64
	vx, vy float64
	size   float64
}

// ゲーム構造体
type Game struct {
	paddle   Paddle
	ball     Ball
	blocks   []Block
	score    int
	gameOver bool
	gameWon  bool
}

// ゲーム初期化
func NewGame() *Game {
	g := &Game{
		paddle: Paddle{
			x:      screenWidth/2 - paddleWidth/2,
			y:      screenHeight - 50,
			width:  paddleWidth,
			height: paddleHeight,
		},
		ball: Ball{
			x:    screenWidth / 2,
			y:    screenHeight / 2,
			vx:   ballSpeed,
			vy:   -ballSpeed,
			size: ballSize,
		},
		score:    0,
		gameOver: false,
		gameWon:  false,
	}

	// ブロックの初期化
	g.blocks = make([]Block, blockRows*blockCols)
	for row := 0; row < blockRows; row++ {
		for col := 0; col < blockCols; col++ {
			idx := row*blockCols + col
			g.blocks[idx] = Block{
				x:       float64(col*blockWidth + 20),
				y:       float64(row*blockHeight + 50),
				width:   blockWidth - 2,
				height:  blockHeight - 2,
				visible: true,
				color:   getBlockColor(row),
			}
		}
	}

	return g
}

// ブロックの色を行に応じて設定
func getBlockColor(row int) color.Color {
	colors := []color.Color{
		color.RGBA{255, 0, 0, 255},   // 赤
		color.RGBA{255, 165, 0, 255}, // オレンジ
		color.RGBA{255, 255, 0, 255}, // 黄
		color.RGBA{0, 255, 0, 255},   // 緑
		color.RGBA{0, 0, 255, 255},   // 青
		color.RGBA{128, 0, 128, 255}, // 紫
	}
	return colors[row%len(colors)]
}

// 衝突検出
func (g *Game) checkCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

// ゲーム更新
func (g *Game) Update() error {
	if g.gameOver || g.gameWon {
		// ゲーム終了時、スペースキーでリスタート
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			*g = *NewGame()
		}
		return nil
	}

	// パドルの移動
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.paddle.x > 0 {
		g.paddle.x -= paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && g.paddle.x < screenWidth-g.paddle.width {
		g.paddle.x += paddleSpeed
	}

	// ボールの移動
	g.ball.x += g.ball.vx
	g.ball.y += g.ball.vy

	// 壁との衝突（左右）
	if g.ball.x <= 0 || g.ball.x >= screenWidth-g.ball.size {
		g.ball.vx = -g.ball.vx
	}

	// 天井との衝突
	if g.ball.y <= 0 {
		g.ball.vy = -g.ball.vy
	}

	// 底面に到達（ゲームオーバー）
	if g.ball.y >= screenHeight {
		g.gameOver = true
		return nil
	}

	// パドルとの衝突
	if g.checkCollision(g.ball.x, g.ball.y, g.ball.size, g.ball.size,
		g.paddle.x, g.paddle.y, g.paddle.width, g.paddle.height) {
		// パドルの位置に応じて角度を変更
		relativePos := (g.ball.x + g.ball.size/2 - g.paddle.x) / g.paddle.width
		angle := (relativePos - 0.5) * math.Pi / 3 // -60度から+60度
		speed := math.Sqrt(g.ball.vx*g.ball.vx + g.ball.vy*g.ball.vy)
		g.ball.vx = speed * math.Sin(angle)
		g.ball.vy = -math.Abs(speed * math.Cos(angle)) // 常に上向き
	}

	// ブロックとの衝突
	for i := range g.blocks {
		if !g.blocks[i].visible {
			continue
		}

		if g.checkCollision(g.ball.x, g.ball.y, g.ball.size, g.ball.size,
			g.blocks[i].x, g.blocks[i].y, g.blocks[i].width, g.blocks[i].height) {
			g.blocks[i].visible = false
			g.score += 10

			// 横方向の衝突か縦方向の衝突かを判定
			overlapX := math.Min(g.ball.x+g.ball.size-g.blocks[i].x, g.blocks[i].x+g.blocks[i].width-g.ball.x)
			overlapY := math.Min(g.ball.y+g.ball.size-g.blocks[i].y, g.blocks[i].y+g.blocks[i].height-g.ball.y)

			if overlapX < overlapY {
				g.ball.vx = -g.ball.vx
			} else {
				g.ball.vy = -g.ball.vy
			}
			break
		}
	}

	// 勝利条件の確認
	allDestroyed := true
	for _, block := range g.blocks {
		if block.visible {
			allDestroyed = false
			break
		}
	}
	if allDestroyed {
		g.gameWon = true
	}

	return nil
}

// 描画
func (g *Game) Draw(screen *ebiten.Image) {
	// 背景を黒で塗りつぶし
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// パドルの描画
	vector.DrawFilledRect(screen, float32(g.paddle.x), float32(g.paddle.y),
		float32(g.paddle.width), float32(g.paddle.height),
		color.RGBA{255, 255, 255, 255}, false)

	// ボールの描画
	vector.DrawFilledCircle(screen, float32(g.ball.x+g.ball.size/2), float32(g.ball.y+g.ball.size/2),
		float32(g.ball.size/2), color.RGBA{255, 255, 255, 255}, false)

	// ブロックの描画
	for _, block := range g.blocks {
		if block.visible {
			vector.DrawFilledRect(screen, float32(block.x), float32(block.y),
				float32(block.width), float32(block.height), block.color, false)
		}
	}

	// スコアの表示
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.score))

	// ゲーム終了メッセージ
	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, "GAME OVER! Press SPACE to restart", screenWidth/2-150, screenHeight/2)
	} else if g.gameWon {
		ebitenutil.DebugPrintAt(screen, "YOU WIN! Press SPACE to restart", screenWidth/2-150, screenHeight/2)
	} else {
		// 操作説明
		ebitenutil.DebugPrintAt(screen, "Use LEFT/RIGHT arrow keys to move paddle", 10, screenHeight-40)
	}
}

// レイアウト
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Breakout Game - ブロック崩し")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}