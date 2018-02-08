package main

import (
    "bufio"
    "fmt"
    "strings"
    // "reflect"

    "github.com/nsf/termbox-go"
)

const (
    coldef = termbox.ColorDefault
    NUM_ROWS = 20
    lbl_offset = 4
    x_offset = 30 + lbl_offset
    y_offset = 5

    MAX_BOARD_SIZE = 26
    MIN_BOARD_SIZE = 3
    DEF_BOARD_SIZE = 5
    MAX_MINE_PCT = .99
    MIN_MINE_PCT = .05
    DEF_MINE_PCT = .15
    ROW_NAMES = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    COL_NAMES = "abcdefghijklmnopqrstuvwxyz"
    ROW_OFFSET = 65
    COL_OFFSET = 97
)

// Structs

type UserAction struct {
    Action    string
    Subject   string
}

// Globals

var output_mode = termbox.OutputNormal
var COLORS = map[string]termbox.Attribute{
    "default": coldef,
    "white": termbox.ColorWhite,
    "bomb": termbox.ColorRed,
    "marked": termbox.ColorYellow,
}
var brd *GameBoard
var rdr *bufio.Reader

// general fuckery

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

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
    // fmt.Printf("Setting %v,%v to %v\n", x, y, msg)
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func draw_brd(brd *GameBoard) {
    header := header_str(brd)
    tbprint(0, y_offset - 1, COLORS["default"], COLORS["default"], header)
    for x := 0; x < len(brd.Rows); x++ {
        x_pos := 2 * x + x_offset

        if len(brd.Rows[x]) > 0 {
            for y := 0; y < len(brd.Rows[x]); y++ {
                y_pos := y + y_offset
                c := &brd.Rows[x][y]
                // fmt.Printf("%v,%v: %v\n", x, y, c.Value)

                if x == 0 {
                    tbprint(x_pos - lbl_offset, y_pos, COLORS["default"], COLORS["default"],
                            fmt.Sprintf("%v [ ", string(ROW_NAMES[y])))
                } else if x == len(brd.Rows) - 1 {
                    tbprint(x_pos + 1, y_pos, COLORS["default"], COLORS["default"],
                            fmt.Sprintf(" ] %v", string(ROW_NAMES[y])))
                }

                if brd.State == "success" || brd.State == "dead" {
                    if c.IsBomb {
                        tbprint(x_pos, y_pos, COLORS["bomb"], COLORS["default"], c.Value)
                    } else {
                        // fmt.Printf("%v,%v Val: %v (%v)\n", x, y, c.Value, reflect.TypeOf(c.Value))
                        tbprint(x_pos, y_pos, COLORS["white"], COLORS["default"], c.Value)
                    }
                } else {
                    if c.IsMarked {
                        tbprint(x_pos, y_pos, COLORS["marked"], COLORS["default"], c.Display)
                    } else {
                        tbprint(x_pos, y_pos, COLORS["default"], COLORS["default"], c.Display)
                    }
                }
            }
        }
    }
    tbprint(0, y_offset + len(brd.Rows), COLORS["default"], COLORS["default"], header)

    termbox.Flush()
}

func header_str(brd *GameBoard) string {
    hdr := strings.Repeat(" ", x_offset)
    // hdr := ""
    for x := range brd.Rows[0] {
        hdr += fmt.Sprintf("%v ", string(COL_NAMES[x]))
    }
    return hdr
}

func main() {
    err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

    brd = NewGameBoard(NUM_ROWS, DEF_MINE_PCT)
    brd.Parse()

loop:
	for {
    	draw_brd(brd)
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
            default:
                switch ev.Ch {
                case 'q':
                    break loop
                case 'r':
                    // reveal cell
                case 'm':
                    // mark cell
                case 'n':
                    brd = NewGameBoard(NUM_ROWS, DEF_MINE_PCT)
                    brd.Parse()
    			}
            }
        }
	}
}
