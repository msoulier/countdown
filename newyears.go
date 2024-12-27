package main

import (
    "fmt"
    "time"
    "github.com/msoulier/mlib"
)

func main() {
    now := time.Now().UTC()
    fmt.Printf("%v\n", now)
    newyears := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
    fmt.Printf("%v\n", newyears)
    diff := newyears.Sub(now)
    fmt.Printf("%s until %s\n", mlib.Duration2Human(diff), newyears)
}
