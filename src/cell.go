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
        return c.Value
    } else {
        return c.Display
    }
}

func (c *Cell) Mark() {
    // fmt.Println("c was: ", c.IsMarked)
    c.IsMarked = !c.IsMarked
    // fmt.Println("c is:  ", c.IsMarked)
    if c.IsMarked {
        c.Display = "?"
    } else {
        c.Display = "-"
    }
}
