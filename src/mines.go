package main

import (
    "bufio"
    "fmt"
    // "math/rand"
    "os"
    "strings"
    // "time"
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
    checkReplay := true

    for true {
        brd = NewGameBoard(8, 0.15)
        brd.Parse()

        // guess := [2]int{0,0}
        for brd.State != "dead" {
            fmt.Println(brd)
            fmt.Printf("%v/%v revealed, board state: %v\n", brd.Revealed, brd.Success, brd.State)
            fmt.Println("[r]eveal or [m]ark a cell: ")
            uact, err := GetAction()
            if err != nil {
                fmt.Println(err)
                continue
            }
            // fmt.Println(uact)
            if uact.Action == "quit" || uact.Action == "help" {
                brd.State = "dead"
                checkReplay = false
                break
            } else if uact.Action == "reveal" {
                i,j,err := coord2idx(uact.Subject)
                if err != nil {
                    fmt.Println(err)
                    continue
                }
                // fmt.Printf("Revealing %v,%v\n", i, j)
                err = brd.Reveal(i, j)
                if err != nil {
                    fmt.Printf("\n  %v\n", err)
                }
            } else if uact.Action == "mark" {
                i,j,err := coord2idx(uact.Subject)
                if err != nil {
                    fmt.Println(err)
                    continue
                }
                // fmt.Printf("Marking %v,%v\n", i, j)
                brd.Mark(i, j)
            }

            if brd.State == "success" {
                fmt.Println("\n\t *** Congratulations, you win! ***")
                // brd.ShowAll()
                break
            }
        }

        fmt.Println(brd)
        fmt.Println("Game over")
        if checkReplay {
            playAgain := GetYesNo("Play again? [Y/n]  ", "y")
            if playAgain == "n" {
                break
            }
        } else {
            break
        }
    }
}
