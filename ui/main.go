package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	engine "github.com/g0g05arui/chess-engine/game_state"
)

var pieceImages map[string]image.Image

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

func main() {
	loadPieceImages()

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

			drawChessBoard(gtx, board)

			move, _ := engine.BestMove(board, 6, color)
			board = engine.BoardAfterMove(move, board)
			if color == engine.WhiteColor {
				color = engine.BlackColor
			} else {
				color = engine.WhiteColor
			}
			fmt.Println(move)
			e.Frame(gtx.Ops)
			w.Invalidate()

		}
	}
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
