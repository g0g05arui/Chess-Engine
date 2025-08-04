package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	engine "github.com/g0g05arui/chess-engine/game_state"
)

// game state
var pieceImages map[string]image.Image
var gameStarted = false
var gameEnded = false
var gameEndReason = ""
var board = engine.CreateBoard()
var colorTurn = engine.WhiteColor
var depthSlider widget.Float
var selectedDepth int = 4
var moveStartTime time.Time
var isCalculatingMove bool = false

// UI elements
var startButton widget.Clickable
var newGameButton widget.Clickable
var theme *material.Theme

func loadPieceImages() {
	pieceImages = make(map[string]image.Image)

	files := []string{
		"bishop_black", "king_black", "knight_black", "pawn_black", "queen_black", "rook_black",
		"bishop_white", "king_white", "knight_white", "pawn_white", "queen_white", "rook_white",
	}

	for _, name := range files {
		path := "assets/pieces/" + name + ".png"
		f, err := os.Open(path)
		if err != nil {
			log.Fatalf("error loading %s: %v", name, err)
		}
		img, _, err := image.Decode(f)
		if err != nil {
			log.Fatalf("error decoding %s: %v", name, err)
		}
		pieceImages[name] = img
		f.Close()
	}
}

func initTheme() {
	theme = material.NewTheme()
}

func main() {
	loadPieceImages()
	initTheme()
	go func() {
		w := new(app.Window)
		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func checkGameEnd(board engine.Board) (bool, string) {
	// Check for threefold repetition
	fen := engine.BoardToFEN(board)
	if count, exists := board.Played[fen]; exists && count >= 3 {
		return true, "Draw by threefold repetition"
	}

	// You can add other endgame conditions here:
	// - Checkmate detection
	// - Stalemate detection
	// - Fifty-move rule
	// - Insufficient material

	return false, ""
}

func resetGame() {
	gameStarted = false
	gameEnded = false
	gameEndReason = ""
	board = engine.CreateBoard()
	colorTurn = engine.WhiteColor
	isCalculatingMove = false
}

func run(w *app.Window) error {
	var ops op.Ops
	board := engine.CreateBoard()

	w.Option(app.Size(600, 600), app.MaxSize(600, 600), app.MinSize(600, 600))

	color := engine.WhiteColor

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if gameStarted && !gameEnded {
				drawChessBoard(gtx, board)

				// Check for game end conditions
				if ended, reason := checkGameEnd(board); ended {
					gameEnded = true
					gameEndReason = reason
					fmt.Println("Game ended:", reason)
				} else if !isCalculatingMove {
					// Start calculating the move
					isCalculatingMove = true
					moveStartTime = time.Now()

					go func() {
						move, _ := engine.BestMove(board, selectedDepth, color)

						// Ensure at least 1 second has passed
						elapsed := time.Since(moveStartTime)
						if elapsed < time.Second {
							time.Sleep(time.Second - elapsed)
						}

						// Apply the move
						board = engine.BoardAfterMove(move, board)
						if color == engine.WhiteColor {
							color = engine.BlackColor
						} else {
							color = engine.WhiteColor
						}
						fmt.Println(move)

						isCalculatingMove = false
						w.Invalidate()
					}()
				}

				e.Frame(gtx.Ops)
				w.Invalidate()
			} else if gameEnded {
				drawGameEndScreen(gtx, w)
				e.Frame(gtx.Ops)
			} else {
				draw_menu(gtx, w)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func drawGameEndScreen(gtx layout.Context, w *app.Window) {
	// Initialize theme if not already done
	if theme == nil {
		initTheme()
	}

	// Fill background with light gray
	paint.Fill(gtx.Ops, color.NRGBA{R: 240, G: 240, B: 240, A: 255})

	// Center the content
	layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(gtx,
			// Game Over title
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H1(theme, "Game Over")
				title.Alignment = text.Middle
				title.Color = color.NRGBA{R: 180, G: 50, B: 50, A: 255}
				return title.Layout(gtx)
			}),
			// Spacing
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
			}),
			// End reason
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				reason := material.H3(theme, gameEndReason)
				reason.Alignment = text.Middle
				reason.Color = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
				return reason.Layout(gtx)
			}),
			// Spacing
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(40)}.Layout(gtx)
			}),
			// New Game button
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(theme, &newGameButton, "New Game")
				btn.CornerRadius = unit.Dp(8)
				btn.Background = color.NRGBA{R: 70, G: 130, B: 180, A: 255}

				// Check if button was clicked
				if newGameButton.Clicked(gtx) {
					resetGame()
					w.Invalidate()
				}

				return btn.Layout(gtx)
			}),
		)
	})
}

func drawChessBoard(gtx layout.Context, board engine.Board) {
	boardSize := gtx.Constraints.Max.X
	squareSize := boardSize / 8
	pieceDrawSize := int(float32(squareSize) * 0.9)
	lightColor := color.NRGBA{R: 240, G: 217, B: 181, A: 255}
	darkColor := color.NRGBA{R: 181, G: 136, B: 99, A: 255}

	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			x := file * squareSize
			y := rank * squareSize

			square := clip.Rect{
				Min: image.Point{x, y},
				Max: image.Point{x + squareSize, y + squareSize},
			}.Push(gtx.Ops)

			color := lightColor
			if (rank+file)%2 == 1 {
				color = darkColor
			}
			paint.Fill(gtx.Ops, color)
			square.Pop()

			// Draw piece if exists
			p := board.PiecesMatrix[8-rank][file+1]
			if p.Type != 0 {
				img := pieceImages[pieceKey(p)]
				if img == nil {
					continue
				}

				// Calculate centering offset
				margin := (squareSize - pieceDrawSize) / 2
				pieceX := x + margin
				pieceY := y + margin

				// Apply offset transformation
				off := op.Offset(image.Pt(pieceX, pieceY)).Push(gtx.Ops)

				// Create clipping rectangle for the piece
				pieceClip := clip.Rect{Max: image.Pt(pieceDrawSize, pieceDrawSize)}.Push(gtx.Ops)

				// Get original image dimensions
				imgBounds := img.Bounds()
				imgWidth := imgBounds.Dx()
				imgHeight := imgBounds.Dy()

				// Calculate scale to fit the piece in the desired size
				scaleX := float32(pieceDrawSize) / float32(imgWidth)
				scaleY := float32(pieceDrawSize) / float32(imgHeight)
				scale := scaleX
				if scaleY < scaleX {
					scale = scaleY
				}

				// Apply scaling transformation
				scaleOp := op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scale, scale))).Push(gtx.Ops)

				// Draw the image
				paint.NewImageOp(img).Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)

				scaleOp.Pop()
				pieceClip.Pop()
				off.Pop()
			}
		}
	}
}

func pieceKey(p engine.Piece) string {
	var name string
	switch p.Type {
	case engine.Pawn:
		name = "pawn"
	case engine.Knight:
		name = "knight"
	case engine.Bishop:
		name = "bishop"
	case engine.Rook:
		name = "rook"
	case engine.Queen:
		name = "queen"
	case engine.King:
		name = "king"
	default:
		return ""
	}
	if p.Color == engine.WhiteColor {
		return name + "_white"
	}
	return name + "_black"
}

func draw_menu(gtx layout.Context, w *app.Window) {
	// Initialize theme if not already done
	if theme == nil {
		initTheme()
	}

	// Fill background with white
	paint.Fill(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255})

	// Center the menu content
	layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(gtx,
			// Title
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H2(theme, "Chess Game")
				title.Alignment = text.Middle
				return title.Layout(gtx)
			}),
			// Spacing
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(40)}.Layout(gtx)
			}),
			// Depth selector label
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Body1(theme, fmt.Sprintf("AI Depth: %d", selectedDepth))
				label.Alignment = text.Middle
				return label.Layout(gtx)
			}),
			// Spacing
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
			}),
			// Depth slider
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				// Update selected depth based on slider value (1-8 range)
				sliderValue := depthSlider.Value
				selectedDepth = int(sliderValue*7) + 1 // Maps 0.0-1.0 to 1-8

				// Set initial slider position for depth 4
				if depthSlider.Value == 0 && selectedDepth == 1 {
					depthSlider.Value = 3.0 / 7.0 // Set to position for depth 4
					selectedDepth = 4
				}

				slider := material.Slider(theme, &depthSlider)
				slider.Color = color.NRGBA{R: 70, G: 130, B: 180, A: 255}

				// Constrain slider width
				gtx.Constraints.Max.X = gtx.Dp(unit.Dp(200))
				return slider.Layout(gtx)
			}),
			// Spacing
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
			}),
			// Depth description
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				var description string
				switch {
				case selectedDepth <= 2:
					description = "Easy - Fast moves"
				case selectedDepth <= 4:
					description = "Medium - Balanced"
				case selectedDepth <= 6:
					description = "Hard - Strong play"
				default:
					description = "Expert - Very slow"
				}

				desc := material.Caption(theme, description)
				desc.Alignment = text.Middle
				desc.Color = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
				return desc.Layout(gtx)
			}),
			// Spacing
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
			}),
			// Start button
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(theme, &startButton, "Start Game")
				btn.CornerRadius = unit.Dp(8)

				// Check if button was clicked
				if startButton.Clicked(gtx) {
					gameStarted = true
					w.Invalidate()
				}

				return btn.Layout(gtx)
			}),
		)
	})
}
