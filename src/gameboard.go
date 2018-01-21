package main

import (
    "fmt"
    "math/rand"
    "time"
)

type GameBoard struct {
    Rows      [][]Cell
    Dims      int
    State     string
    NumMoves  int
    NumMines  int
    Revealed  int
    Marked    int
    Success   int
}

func (b GameBoard) String() string {
    header := "    "
    for x := range b.Rows {
        header += fmt.Sprintf("%v ", string(COL_NAMES[x]))
    }
    brd_str := fmt.Sprintf("\n%v\n", header)
    for i := 0; i < len(b.Rows); i++ {
        row_str := fmt.Sprintf("%v [ ", string(ROW_NAMES[i]))
        for j := 0; j < len(b.Rows[i]); j++ {
            row_str += fmt.Sprintf("%v ", b.Rows[i][j])
        }
        row_str += fmt.Sprintf("] %v\n", string(ROW_NAMES[i]))

        brd_str += row_str
    }
    brd_str += fmt.Sprintf("%v\n", header)
    return brd_str
}

func (b *GameBoard) Mark(i int, j int) {
    c := &b.Rows[i][j]
    if !c.IsRevealed {
        b.Marked += c.Mark()
    } else {
        fmt.Println("Cell has already been revealed")
    }
}

func (b *GameBoard) Parse() {
    for i := 0; i < len(b.Rows); i++ {
        for j := 0; j < len(b.Rows[i]); j++ {
            c := &b.Rows[i][j]
            if c.IsBomb {
                c.Value = "*"
            } else {
                val := 0
                for i0 := i-1; i0 < i+2; i0++ {
                    if i0 < 0 || i0 >= len(b.Rows[i]) {
                        continue
                    }
                    for j0 := j-1; j0 < j+2; j0++ {
                        if j0 < 0 || j0 >= len(b.Rows[i0]) {
                            continue
                        }
                        if b.Rows[i0][j0].IsBomb {
                            val++
                        }
                    }
                }
                if val == 0 {
                    c.Value = " "
                } else {
                    c.Value = fmt.Sprintf("%v", val)
                }
            }
        }
    }
}

func (b *GameBoard) Reveal(i int, j int) (err error) {
    cLabel := fmt.Sprintf("%v,%v", string(i+ROW_OFFSET), string(j+COL_OFFSET))
    if b.State != "active" {
        b.State = "active"
    }
    c := &b.Rows[i][j]
    if c.IsMarked {
        err = fmt.Errorf("Cannot reveal marked cell %v", cLabel)
    } else if c.IsRevealed {
        err = fmt.Errorf("Already revealed cell %v: pick another", cLabel)
    } else {
        c.IsRevealed = true
        b.NumMoves++
        if c.IsBomb {
            brd.State = "dead"
            c.Value = "!"
            fmt.Println("\n\n\t*** BOOOOOM ***\n")
        } else {
            c.Display = c.Value
            b.Revealed++
        }
        if c.Value == " " {
            for new_i := i-1; new_i < i+2; new_i++ {
                if new_i < 0 || new_i >= len(b.Rows) {
                    continue
                }
                for new_j := j-1; new_j < j+2; new_j++ {
                    if new_j < 0 || new_j >= len(b.Rows[new_i]) {
                        continue
                    } else if b.Rows[new_i][new_j].IsRevealed {
                        continue
                    } else if new_i == i && new_j == j {
                        continue
                    }
                    b.Reveal(new_i, new_j)
                }
            }
        }
    }

    if b.Revealed == b.Success {
        b.State = "success"
    }
    return err
}

func NewGameBoard(dim int, mine_pct float64) *GameBoard {
    if dim == 0 {
        dim = DEF_BOARD_SIZE
    }
    if mine_pct == 0 {
        mine_pct = DEF_MINE_PCT
    }
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    num_mines := int(mine_pct * float64(dim * dim))
    total_mines := num_mines
    fmt.Printf("Initializing board with %v tiles and %v mines", dim * dim, num_mines)

    rows := make([][]Cell, dim)
    for num_mines > 0 {
        for i := range rows {
            if rows[i] == nil {
                rows[i] = make([]Cell, dim)
            }
            for j := range rows[i] {
                if &rows[i][j] == nil || !rows[i][j].IsBomb {
                    is_bomb := false
                    if num_mines > 0 && r.Float64() < mine_pct {
                        is_bomb = true
                        num_mines--
                    }
                    rows[i][j] = Cell{
                        IsBomb: is_bomb,
                        IsMarked: false,
                        Display: "-",
                    }
                }
            }
        }
    }
    new_brd := GameBoard{
        Rows:      rows,
        Dims:      dim,
        State:     "new",
        NumMoves:  0,
        NumMines:  total_mines,
        Revealed:  0,
        Marked:    0,
        Success:   dim * dim - total_mines,
    }

    return &new_brd
}
