package main

import (
    "fmt"
    "time"
    "github.com/msoulier/mlib"
)

func main() {
    now := time.Now().UTC()
    fmt.Printf("The time is now %s\n", now.Local())
    newyears := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local)
    diff := newyears.Sub(now)
    fmt.Printf("It is %s until local New Years\n", mlib.Duration2Human(diff))
}
