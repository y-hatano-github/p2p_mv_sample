package main

import (
	"math"
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetOutputMode(termbox.Output256)

	key := make(chan string)
	go keyEvent(key)

	fps := NewFPS(60) // FPS 60

	r := *rand.New(rand.NewSource(time.Now().UnixNano()))
	w, h := termbox.Size()
	sX := float64(r.Intn(w))
	sY := float64(r.Intn(h))
	eX := float64(r.Intn(w))
	eY := float64(r.Intn(h))

	t := math.Atan2(eY-sY, eX-sX) // angle Between Two Points

	vX := math.Cos(t) // amount of movement of X
	vY := math.Sin(t) // amount of movement of Y

	d := int(math.Sqrt(math.Pow(eX-sX, 2) + math.Pow(eY-sY, 2))) // distance between two points

	m := 0 // amount of movement

loop:
	for {
		termbox.SetCell(int(sX), int(sY), ' ', termbox.ColorBlack, termbox.ColorRed)
		termbox.SetCell(int(eX), int(eY), ' ', termbox.ColorBlack, termbox.ColorBlue)
		fps.Update()
		select {
		case k := <-key:
			if k == "esc" {
				break loop
			}
		default:
		}

		sX += vX
		sY += vY
		m++

		termbox.Flush()
		fps.Wait()

		if m >= d {
			sX = eX
			sY = eY
			eX = float64(r.Intn(w))
			eY = float64(r.Intn(h))
			t := math.Atan2(eY-sY, eX-sX)

			vX = math.Cos(t)
			vY = math.Sin(t)

			d = int(math.Sqrt(math.Pow(eX-sX, 2) + math.Pow(eY-sY, 2)))
			m = 0
			termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}

func keyEvent(key chan string) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				key <- "esc"
			case termbox.KeyCtrlC:
				key <- "esc"
			default:
				key <- string(ev.Ch)
			}
		}
	}
}

type FPS struct {
	startTime int64
	endTime   int64
	fps       int64
	frame     int64
}

func NewFPS(fps int64) *FPS {
	return &FPS{
		fps:       fps,
		startTime: 0,
		endTime:   0,
		frame:     0,
	}
}

func (f *FPS) Update() {
	if f.frame == 0 {
		f.startTime = time.Now().UnixNano() / int64(time.Millisecond)
	}
	if f.frame == f.fps {
		f.frame = 0
		f.startTime = time.Now().UnixNano() / int64(time.Millisecond)
	}
	f.frame++
}

func (f *FPS) Wait() {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	wait := f.frame*1000/f.fps - (now - f.startTime)

	for time.Now().UnixNano()/int64(time.Millisecond)-now < wait {
	}
}
