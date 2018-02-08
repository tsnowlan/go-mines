package main

/*import (
    "fmt"
)*/

type Cell struct {
    IsBomb      bool
    IsMarked    bool
    IsRevealed  bool
    Display     string
    Value       string
}

func (c Cell) String() string {
    if brd.State == "dead" || brd.State == "success" {
        return string(c.Value)
    } else {
        return string(c.Display)
    }
}

// func (c Cell) Rune() rune {
//     if brd.State == "dead" || brd.State == "success" {
//         return c.Value
//     } else {
//         return c.Display
//     }
// }

func (c *Cell) Mark() (net int) {
    c.IsMarked = !c.IsMarked
    if c.IsMarked {
        c.Display = "?"
        net = 1
    } else {
        c.Display = "-"
        net = -1
    }
    return net
}
