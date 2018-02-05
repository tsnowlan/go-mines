package main

import (
    "bufio"
    "fmt"
    // "math/rand"
    "os"
    "strings"
    // "time"

    ui "github.com/gizak/termui"
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

    ih = 3
    th = 4
)

// Structs

type UserAction struct {
    Action    string
    Subject   string
}

// Globals

var brd *GameBoard
var rdr *bufio.Reader

// general funcs

func GetAction() (act UserAction, err error) {
    line, rd_err := rdr.ReadString('\n')
    if rd_err != nil {
        return act, fmt.Errorf("Error parsing input: %v", err)
    }
    line = strings.Trim(line, "\n" )
    if line == "" {
        return act, fmt.Errorf("Cannot parse empty string")
    }
    resp := strings.Split(line, " ")
    resp[0] = strings.ToLower(resp[0])
    if len(resp) >= 2 {
        act.Subject = resp[1]
    }

    if resp == nil {
        err = fmt.Errorf("Unable to parse '%v'", line)
    } else if resp[0] == "reveal" || resp[0] == "r" {
        act.Action = "reveal"
        act.Subject = resp[1]
    } else if resp[0] == "mark" || resp[0] == "m" {
        act.Action = "mark"
        act.Subject = resp[1]
    } else if resp[0] == "quit" || resp[0] == "q" {
        act.Action = "quit"
    } else if resp[0] == "help" || resp[0] == "h" {
        act.Action = "help"
    } else {
        err = fmt.Errorf("Invalid action: '%v'\n", resp[0])
    }
    return act, err
}

func GetYesNo(msg string, default_value string) (str string) {
    fmt.Printf(msg)
    resp := ""
    for resp == "" {
        line, rd_err := rdr.ReadString('\n')
        if rd_err != nil {
            fmt.Printf("Error parsing input: %v\n", rd_err)
            fmt.Println(msg)
            continue
        }

        line = strings.Trim(line, "\n" )
        if line == "" && default_value != "" {
            line = default_value
        }

        if strings.ToLower(line) == "y" || strings.ToLower(line) == "yes" {
            resp = "y"
        } else if strings.ToLower(line) == "n" || strings.ToLower(line) == "no" {
            resp = "n"
        } else {
            fmt.Println("You must specify y/yes or n/no")
            fmt.Println(msg)
            continue
        }
    }
    return resp
}

func coord2idx(str string) (i int, j int, err error) {
    if len(str) != 3 || strings.Index(str, ",") == -1 {
        i, j = -1, -1
        err = fmt.Errorf("Invalid coord string: '%v'", str)
    } else {
        i_j := strings.Split(str, ",")
        // convert individual runes to int representation
        i = int(i_j[0][0]) - ROW_OFFSET
        j = int(i_j[1][0]) - COL_OFFSET
    }
    return i, j, err
}

// Main

func main() {
    rdr = bufio.NewReader(os.Stdin)
    // checkReplay := true

    err := ui.Init()
    if err != nil {
        panic(err)
    }
    defer ui.Close()

    brd = NewGameBoard(8, 0.15)
    brd.Parse()
    fmt.Println("Generated GameBoard with dims %v", brd.Dims)

    ui_width := 3 * brd.Dims + 10
    // max_width = ui.TermWidth - 5

    header := ui.NewPar("Go Mines!")
    header.Height = 3
    header.Width = ui_width
    // header.Border = false

    board_ui := ui.NewPar(brd.String())
    board_ui.BorderLabel = "Gameboard"
    board_ui.Width = ui_width
    board_ui.Height = 2 * brd.Dims
    board_ui.Y = header.Height

    tb_msg := fmt.Sprintf("%v/%v revealed, %v/%v marked/bombs\n[r]eveal or [m]ark a cell", brd.Revealed, brd.Success, brd.Marked, brd.NumMines)
    // fmt.Println("[r]eveal or [m]ark a cell: ")
    tb := ui.NewPar(tb_msg)
    tb.Height = th
    tb.Y = board_ui.Y + board_ui.Height
    tb.Width = ui_width

    ib := ui.NewPar("")
	ib.Height = ih
	ib.BorderLabel = "Input"
	ib.BorderLabelFg = ui.ColorYellow
	ib.BorderFg = ui.ColorYellow
	ib.TextFgColor = ui.ColorWhite
    ib.Y = tb.Y + tb.Height
    ib.Width = ui_width

    ui.Body.AddRows(
        ui.NewRow(ui.NewCol(12, 0, header)),
        ui.NewRow(ui.NewCol(12, 0, board_ui)),
        ui.NewRow(ui.NewCol(12, 0, tb)),
        ui.NewRow(ui.NewCol(12, 0, ib)))

    ui.Render(ui.Body)

    ui.Handle("/sys/kbd/q", func(ui.Event) {
        ui.StopLoop()
    })
    ui.Handle("/sys/kbd/r", func(ui.Event) {
        ui.StopLoop()
    })
    ui.Loop()

}
