package main

import (
    "fmt"
    "math/rand"
    "os"
    "path/filepath"
    "time"

    py "github.com/fohristiwhirl/gopycanvas"
)

const (
    WIDTH = 800
    HEIGHT = 600

    GFX_DIR_NAME = "gfx"
)

type ball struct {
    x, y, speedx, speedy int32
}

func (b *ball) move() {
    if (b.x > WIDTH && b.speedx > 0) || (b.x < 0 && b.speedx < 0) {
        b.speedx *= -1
    }
    if (b.y > HEIGHT && b.speedy > 0) || (b.y < 0 && b.speedy < 0) {
        b.speedy *= -1
    }
    b.x += b.speedx
    b.y += b.speedy
}

func (b *ball) playermove() {
    b.speedx, b.speedy = 0, 0
    if py.KeyDown("a") {
        b.speedx -= 2
    }
    if py.KeyDown("d") {
        b.speedx += 2
    }
    if py.KeyDown("w") {
        b.speedy -= 2
    }
    if py.KeyDown("s") {
        b.speedy += 2
    }
    b.move()
}

func main() {

    our_dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        fmt.Fprintf(os.Stderr, "Couldn't find our own directory!\n")
        os.Exit(1)
    }
    gfx_path := filepath.Join(our_dir, GFX_DIR_NAME)

    py.Start(WIDTH, HEIGHT, gfx_path, "black")

    var balls []*ball
    var player *ball

    for n := 0 ; n < 50 ; n++ {
        b := ball{x: rand.Int31n(WIDTH), y: rand.Int31n(HEIGHT), speedx: rand.Int31n(4), speedy: rand.Int31n(4)}
        balls = append(balls, &b)
    }

    player = new(ball)
    player.x, player.y = 400, 300

    for {
        for _, b := range balls {
            b.move()
            py.Sprite("white.gif", b.x, b.y)
        }
        player.playermove()
        py.Sprite("green.gif", player.x, player.y)
        py.EndFrame()
        time.Sleep(10 * time.Millisecond)
    }
}
