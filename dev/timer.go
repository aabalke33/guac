package timer

import (
	"fmt"
	"time"
)

type Timer struct {
    Timestamps []Time
}

type Time struct {
    Timestamp time.Time 
    Text string
}

func StartTimer() (*Timer) {
    start := Time{ Timestamp: time.Now(), Text: "Start" }
    println("⏰ Timer Started.")
    return &Timer{ Timestamps: []Time{start} }
}

func (t *Timer) Stop(title string) {

    start := t.Timestamps[0].Timestamp
    stop := time.Now()

    switch len(t.Timestamps) {
    case 0: println("⏰ Timer was never initiated.")
    case 1:
        time := fmt.Sprintf("⏰ Timer Ended. Total %s", stop.Sub(start))
        println(time)
    default:
        prev := t.Timestamps[len(t.Timestamps)-1].Timestamp
        lap := fmt.Sprintf("⏰ Timer Ended. %s %s", title, stop.Sub(prev))
        time := fmt.Sprintf("⏰ Timer Ended. Total %s", stop.Sub(start))
        println(lap)
        println(time)
    }

    //Dereference?
}

func (t *Timer) Lap(title string) {
    if len(t.Timestamps) < 1 {
        println("⏰ Timer was never initiated.")
        return
    }

    prev := t.Timestamps[len(t.Timestamps)-1].Timestamp
    lap := time.Now()

    t.Timestamps = append(t.Timestamps, Time{Timestamp: lap, Text: title})

    text := fmt.Sprintf("⏰ Timer Lap. %s %s", title, lap.Sub(prev))
    println(text)
}
