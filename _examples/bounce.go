package main

import (
    "math/rand"
    "time"

    "github.com/fohristiwhirl/gopycanvas"
)

const (
    WIDTH = 800
    HEIGHT = 600
)

type ball struct {
    x, y, speedx, speedy int32
}

func (b *ball) move() {
    b.x += b.speedx
    b.y += b.speedy
    if (b.x > WIDTH && b.speedx > 0) || (b.x < 0 && b.speedx < 0) {
        b.speedx *= -1
    }
    if (b.y > HEIGHT && b.speedy > 0) || (b.y < 0 && b.speedy < 0) {
        b.speedy *= -1
    }
}

func main() {
    gopycanvas.Start(WIDTH, HEIGHT, "gfx", "black")

    var balls []*ball

    for n := 0 ; n < 50 ; n++ {
        b := ball{x: rand.Int31n(WIDTH), y: rand.Int31n(HEIGHT), speedx: rand.Int31n(4), speedy: rand.Int31n(4)}
        balls = append(balls, &b)
    }

    for {
        for _, b := range balls {
            b.move()
            gopycanvas.Sprite("white.gif", b.x, b.y)
        }
        gopycanvas.EndFrame()
        time.Sleep(10 * time.Millisecond)
    }
}
