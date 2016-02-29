package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "strings"
    "time"
)

const (
    MAX_BOARD_SIZE = 26
    MIN_BOARD_SIZE = 3
    DEF_BOARD_SIZE = 5
    MAX_MINE_PCT = .99
    MIN_MINE_PCT = .05
    DEF_MINE_PCT = .2
    ROW_NAMES = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    COL_NAMES = "abcdefghijklmnopqrstuvwxyz"
    ROW_OFFSET = 65
    COL_OFFSET = 97
)

// Structs

type Cell struct {
    IsBomb      bool
    IsMarked    bool
    IsRevealed  bool
    Display     string
    Value       string
}

type GameBoard struct {
    Rows      [][]Cell
    Dims      int
    State     string
    NumMoves  int
}

type UserAction struct {
    Action    string
    Subject   string
}

// Globals

var brd *GameBoard
var rdr *bufio.Reader

// Class funcs

func (c Cell) String() string {
    if brd.State == "dead" {
        return c.Value
    } else {
        return c.Display
    }
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

func (c Cell) Mark() {
    c.IsMarked = !c.IsMarked
    if c.IsMarked {
        c.Display = "?"
    } else {
        c.Display = "-"
    }
}

func (b GameBoard) Parse() {
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

func (b GameBoard) Reveal(i int, j int) {
    if b.State != "active" {
        b.State = "active"
    }
    c := &b.Rows[i][j]
    c.IsRevealed = true
    b.NumMoves++
    if c.IsBomb {
        brd.State = "dead"
        c.Value = "!"
    } else {
        c.Display = c.Value
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

// general funcs

func NewGameBoard(dim int, mine_pct float64) *GameBoard {
    if dim == 0 {
        dim = DEF_BOARD_SIZE
    }
    if mine_pct == 0 {
        mine_pct = DEF_MINE_PCT
    }
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    num_mines := int(mine_pct * float64(dim * dim))

    rows := make([][]Cell, dim)
    for i := range rows {
        rows[i] = make([]Cell, dim)
        for j := range rows[i] {
            is_bomb := false
            if r.Float64() < mine_pct && num_mines > 0 {
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
    new_brd := GameBoard{
        Rows:      rows,
        Dims:      dim,
        State:     "new",
        NumMoves:  0,
    }

    return &new_brd
}

func GetAction() *UserAction {
    line, err := rdr.ReadString('\n')
    if err != nil {
        fmt.Printf("Error parsing input: %v", err)
        os.Exit(1)
    }
    line = strings.ToLower(strings.Trim(line, "\n" ))
    if line == "" {
        return nil
    }
    resp := strings.Split(line, " ")
    var action, subject string
    if len(resp) >= 2 {
        subject = resp[1]
    }

    if resp == nil {

    } else if resp[0] == "reveal" || resp[0] == "r" {
        action = "reveal"
    } else if resp[0] == "mark" || resp[0] == "m" {
        action = "mark"
    } else if resp[0] == "quit" || resp[0] == "q" {
        action = "quit"
    } else if resp[0] == "help" || resp[0] == "h" {
        action = "help"
    } else {
        fmt.Printf("Invalid action: '%v'\n", resp[0])
        action = ""
    }
    return &UserAction{
        Action: action,
        Subject: subject,
    }
}

func coord2idx(str string) int, int {

}

// Main

func main() {
    rdr = bufio.NewReader(os.Stdin)
    brd = NewGameBoard(8, 0.2)
    brd.Parse()

    guess := [2]int{0,0}
    for brd.State != "dead" {
        fmt.Println("[r]eveal or [m]ark a cell: ")
        uact := GetAction()
        if !brd.Rows[guess[0]][guess[1]].IsRevealed {
            fmt.Printf("Guessing %v,%v\n", string(ROW_NAMES[guess[0]]), string(COL_NAMES[guess[1]]))
            brd.Reveal(guess[0], guess[1])
            fmt.Println(brd)
        }
        if guess[1] < len(brd.Rows[guess[0]])-1 {
            guess[1]++
        } else {
            guess[0]++
            guess[1] = 0
        }
    }

    fmt.Println("Game over")
}
