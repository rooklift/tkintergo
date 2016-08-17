# tkintergo

This Go package spawns a Python subprocess that uses Tkinter to draw sprites according to commands sent to it via the Stdin pipe. Only GIF files are supported. The Python process also sends information about keypresses via Stderr, which the Go package uses to update a map of keys.
