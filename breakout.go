package main

import (
	"fmt"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	WINDOW_WIDTH  = 800
	WINDOW_HEIGHT = 600
	PADDLE_WIDTH  = 100
	PADDLE_HEIGHT = 20
	BALL_SIZE     = 10
	BLOCK_WIDTH   = 75
	BLOCK_HEIGHT  = 30
	BLOCK_ROWS    = 5
	BLOCK_COLS    = 10
	PADDLE_SPEED  = 400.0
	BALL_SPEED    = 300.0
	FPS           = 60
)

type Vec2 struct {
	X, Y float64
}

type Paddle struct {
	Position Vec2
	Width    float64
	Height   float64
}

type Ball struct {
	Position Vec2
	Velocity Vec2
	Size     float64
}

type Block struct {
	Position Vec2
	Width    float64
	Height   float64
	Active   bool
}

type Game struct {
	Paddle     Paddle
	Ball       Ball
	Blocks     [][]Block
	Score      int
	Lives      int
	GameState  int // 0: playing, 1: game over, 2: win
	LastUpdate time.Time
}

// Dummy text rendering functions (can be extended with SDL_ttf)
func renderText(renderer *sdl.Renderer, text string, x, y int32, color sdl.Color) {
	// ダミー実装 - SDL_ttfで拡張可能
	// Dummy implementation - can be extended with SDL_ttf
	fmt.Printf("Text at (%d,%d): %s\n", x, y, text)
}

func drawFilledRect(renderer *sdl.Renderer, x, y, w, h int32, color sdl.Color) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	rect := sdl.Rect{X: x, Y: y, W: w, H: h}
	renderer.FillRect(&rect)
}

func (g *Game) initGame() {
	// パドル初期化
	g.Paddle = Paddle{
		Position: Vec2{X: WINDOW_WIDTH/2 - PADDLE_WIDTH/2, Y: WINDOW_HEIGHT - 50},
		Width:    PADDLE_WIDTH,
		Height:   PADDLE_HEIGHT,
	}

	// ボール初期化
	g.Ball = Ball{
		Position: Vec2{X: WINDOW_WIDTH / 2, Y: WINDOW_HEIGHT / 2},
		Velocity: Vec2{X: BALL_SPEED * 0.6, Y: -BALL_SPEED * 0.8},
		Size:     BALL_SIZE,
	}

	// ブロック初期化
	g.Blocks = make([][]Block, BLOCK_ROWS)
	for row := 0; row < BLOCK_ROWS; row++ {
		g.Blocks[row] = make([]Block, BLOCK_COLS)
		for col := 0; col < BLOCK_COLS; col++ {
			g.Blocks[row][col] = Block{
				Position: Vec2{
					X: float64(col*BLOCK_WIDTH + 10),
					Y: float64(row*BLOCK_HEIGHT + 50),
				},
				Width:  BLOCK_WIDTH,
				Height: BLOCK_HEIGHT,
				Active: true,
			}
		}
	}

	g.Score = 0
	g.Lives = 3
	g.GameState = 0
	g.LastUpdate = time.Now()
}

func (g *Game) checkCollision(rect1X, rect1Y, rect1W, rect1H, rect2X, rect2Y, rect2W, rect2H float64) bool {
	return rect1X < rect2X+rect2W &&
		rect1X+rect1W > rect2X &&
		rect1Y < rect2Y+rect2H &&
		rect1Y+rect1H > rect2Y
}

func (g *Game) updatePaddle(deltaTime float64, keys []uint8) {
	if keys[sdl.SCANCODE_LEFT] != 0 && g.Paddle.Position.X > 0 {
		g.Paddle.Position.X -= PADDLE_SPEED * deltaTime
	}
	if keys[sdl.SCANCODE_RIGHT] != 0 && g.Paddle.Position.X < WINDOW_WIDTH-g.Paddle.Width {
		g.Paddle.Position.X += PADDLE_SPEED * deltaTime
	}
}

func (g *Game) updateBall(deltaTime float64) {
	// ボール移動
	g.Ball.Position.X += g.Ball.Velocity.X * deltaTime
	g.Ball.Position.Y += g.Ball.Velocity.Y * deltaTime

	// 壁との当たり判定
	if g.Ball.Position.X <= 0 || g.Ball.Position.X >= WINDOW_WIDTH-g.Ball.Size {
		g.Ball.Velocity.X = -g.Ball.Velocity.X
	}
	if g.Ball.Position.Y <= 0 {
		g.Ball.Velocity.Y = -g.Ball.Velocity.Y
	}

	// パドルとの当たり判定
	if g.checkCollision(
		g.Ball.Position.X, g.Ball.Position.Y, g.Ball.Size, g.Ball.Size,
		g.Paddle.Position.X, g.Paddle.Position.Y, g.Paddle.Width, g.Paddle.Height,
	) {
		// パドルの中央からの距離に基づいて反射角度を計算
		paddleCenter := g.Paddle.Position.X + g.Paddle.Width/2
		ballCenter := g.Ball.Position.X + g.Ball.Size/2
		normalizedPosition := (ballCenter - paddleCenter) / (g.Paddle.Width / 2)

		speed := math.Sqrt(g.Ball.Velocity.X*g.Ball.Velocity.X + g.Ball.Velocity.Y*g.Ball.Velocity.Y)
		angle := normalizedPosition * math.Pi / 3 // 最大60度

		g.Ball.Velocity.X = speed * math.Sin(angle)
		g.Ball.Velocity.Y = -math.Abs(speed * math.Cos(angle))

		g.Ball.Position.Y = g.Paddle.Position.Y - g.Ball.Size
	}

	// ブロックとの当たり判定
	for row := 0; row < BLOCK_ROWS; row++ {
		for col := 0; col < BLOCK_COLS; col++ {
			block := &g.Blocks[row][col]
			if !block.Active {
				continue
			}

			if g.checkCollision(
				g.Ball.Position.X, g.Ball.Position.Y, g.Ball.Size, g.Ball.Size,
				block.Position.X, block.Position.Y, block.Width, block.Height,
			) {
				block.Active = false
				g.Score += 10

				// 簡単な反射計算
				ballCenterX := g.Ball.Position.X + g.Ball.Size/2
				ballCenterY := g.Ball.Position.Y + g.Ball.Size/2
				blockCenterX := block.Position.X + block.Width/2
				blockCenterY := block.Position.Y + block.Height/2

				if math.Abs(ballCenterX-blockCenterX) > math.Abs(ballCenterY-blockCenterY) {
					g.Ball.Velocity.X = -g.Ball.Velocity.X
				} else {
					g.Ball.Velocity.Y = -g.Ball.Velocity.Y
				}
				return
			}
		}
	}

	// ボールが下に落ちた場合
	if g.Ball.Position.Y > WINDOW_HEIGHT {
		g.Lives--
		if g.Lives <= 0 {
			g.GameState = 1 // ゲームオーバー
		} else {
			// ボールをリセット
			g.Ball.Position = Vec2{X: WINDOW_WIDTH / 2, Y: WINDOW_HEIGHT / 2}
			g.Ball.Velocity = Vec2{X: BALL_SPEED * 0.6, Y: -BALL_SPEED * 0.8}
		}
	}
}

func (g *Game) checkWinCondition() {
	allDestroyed := true
	for row := 0; row < BLOCK_ROWS; row++ {
		for col := 0; col < BLOCK_COLS; col++ {
			if g.Blocks[row][col].Active {
				allDestroyed = false
				break
			}
		}
		if !allDestroyed {
			break
		}
	}
	if allDestroyed {
		g.GameState = 2 // 勝利
	}
}

func (g *Game) update() {
	now := time.Now()
	deltaTime := now.Sub(g.LastUpdate).Seconds()
	g.LastUpdate = now

	if g.GameState != 0 {
		return
	}

	keys := sdl.GetKeyboardState()
	g.updatePaddle(deltaTime, keys)
	g.updateBall(deltaTime)
	g.checkWinCondition()
}

func (g *Game) render(renderer *sdl.Renderer) {
	// 背景クリア
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	// パドル描画
	drawFilledRect(renderer,
		int32(g.Paddle.Position.X), int32(g.Paddle.Position.Y),
		int32(g.Paddle.Width), int32(g.Paddle.Height),
		sdl.Color{R: 255, G: 255, B: 255, A: 255})

	// ボール描画
	drawFilledRect(renderer,
		int32(g.Ball.Position.X), int32(g.Ball.Position.Y),
		int32(g.Ball.Size), int32(g.Ball.Size),
		sdl.Color{R: 255, G: 255, B: 255, A: 255})

	// ブロック描画
	for row := 0; row < BLOCK_ROWS; row++ {
		for col := 0; col < BLOCK_COLS; col++ {
			block := &g.Blocks[row][col]
			if !block.Active {
				continue
			}

			// 行によって色を変える
			var color sdl.Color
			switch row {
			case 0:
				color = sdl.Color{R: 255, G: 0, B: 0, A: 255} // 赤
			case 1:
				color = sdl.Color{R: 255, G: 165, B: 0, A: 255} // オレンジ
			case 2:
				color = sdl.Color{R: 255, G: 255, B: 0, A: 255} // 黄
			case 3:
				color = sdl.Color{R: 0, G: 255, B: 0, A: 255} // 緑
			case 4:
				color = sdl.Color{R: 0, G: 0, B: 255, A: 255} // 青
			}

			drawFilledRect(renderer,
				int32(block.Position.X), int32(block.Position.Y),
				int32(block.Width), int32(block.Height),
				color)
		}
	}

	// UI表示（ダミー関数使用）
	scoreText := fmt.Sprintf("Score: %d", g.Score)
	livesText := fmt.Sprintf("Lives: %d", g.Lives)
	renderText(renderer, scoreText, 10, 10, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	renderText(renderer, livesText, 10, 30, sdl.Color{R: 255, G: 255, B: 255, A: 255})

	if g.GameState == 1 {
		renderText(renderer, "GAME OVER - Press R to restart", WINDOW_WIDTH/2-100, WINDOW_HEIGHT/2, sdl.Color{R: 255, G: 0, B: 0, A: 255})
	} else if g.GameState == 2 {
		renderText(renderer, "YOU WIN! - Press R to restart", WINDOW_WIDTH/2-100, WINDOW_HEIGHT/2, sdl.Color{R: 0, G: 255, B: 0, A: 255})
	}

	renderer.Present()
}

func main() {
	// SDL初期化
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// ウィンドウ作成
	window, err := sdl.CreateWindow("Breakout Game",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	// レンダラー作成
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	// ゲーム初期化
	game := &Game{}
	game.initGame()

	// メインループ
	running := true
	for running {
		// イベント処理
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case sdl.QuitEvent:
				running = false
			case sdl.KeyboardEvent:
				if t.Type == sdl.KEYDOWN {
					switch t.Keysym.Sym {
					case sdl.K_ESCAPE:
						running = false
					case sdl.K_r:
						if game.GameState != 0 {
							game.initGame()
						}
					}
				}
			}
		}

		// ゲーム更新
		game.update()

		// 描画
		game.render(renderer)

		// FPS制御
		sdl.Delay(1000 / FPS)
	}
}
