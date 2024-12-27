package main

import (
    "fmt"
    "time"
    "github.com/msoulier/mlib"
)

var Reset  = "\033[0m"
var Red    = "\033[31m"
var Green  = "\033[32m"
var Yellow = "\033[33m"
var Blue   = "\033[34m"
var Purple = "\033[35m"
var Cyan   = "\033[36m"
var Gray   = "\033[37m"
var White  = "\033[97m"
var width = 70

func telltime() {
    fmt.Printf("\r")
    now := time.Now().UTC()
    year := now.Local().Year()
    next_year := year+1
    newyears := time.Date(next_year, time.January, 1, 0, 0, 0, 0, time.Local)
    diff := newyears.Sub(now)
    printstring := fmt.Sprintf("Countdown: %s until %d", mlib.Duration2Human(diff), next_year)
    finalstring := printstring[:]
    if len(printstring) > width {
        finalstring = printstring[:width]
    }
    format := fmt.Sprintf("%%s%%%ds%%s", width)
    fmt.Printf(format, Purple, finalstring, Reset)
}

func main() {
    for {
        telltime()
        time.Sleep(1*time.Second)
    }
}
