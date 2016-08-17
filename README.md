# tkintergo

This Go package spawns a Python subprocess that uses Tkinter to draw sprites according to commands sent to it via the Stdin pipe. Only GIF files are supported. The Python process also sends information about keypresses via Stderr, which the Go package uses to update a map of keys.

Basic usage looks like this:

```
func main() {
  tkintergo.Start(WIDTH, HEIGHT, GFX_DIRECTORY, BG_COLOUR)    // e.g. (800, 600, "gfx", "black")
  for {
    /* Insert some logic here that moves stuff about */
    tkintergo.Sprite("somesprite.gif", x, y)
    tkintergo.Sprite("anothersprite.gif", x, y)
    tkintergo.EndFrame()
    time.Sleep(10 * time.Millisecond)
  }
}
```

Key states can be accessed with the KeyDown function, e.g. `tkintergo.KeyDown("space")` returns true if the spacebar is down.
