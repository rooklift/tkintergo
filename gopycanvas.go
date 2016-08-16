package gopycanvas

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"
)

const RENDERER = "tkinter_renderer.py"

var stdout_chan = make(chan string)

var stdin_pipe io.WriteCloser
var stdout_pipe io.ReadCloser
var stderr_pipe io.ReadCloser

func pipe_relay(p io.ReadCloser, ch chan string, f io.ReadWriteCloser) {

    // Relay a pipe to both a channel and a file (though either can be nil)

    scanner := bufio.NewScanner(p)

    for scanner.Scan() {

        s := scanner.Text()

        if ch != nil {
            ch <- s
        }

        if f != nil {
            fmt.Fprintf(f, "%s\n", s)
        }
    }
}

func Command(msg string, need_ack bool) {

    // It's important that need_ack is only true for commands that the renderer will actually acknowledge

    msg = strings.TrimSpace(msg)
    fmt.Fprintf(stdin_pipe, msg + "\n")

    if need_ack {
        <- stdout_chan
    }
}

func Sprite(filename string, x int32, y int32) {
    s := fmt.Sprintf("%s %d %d", filename, x, y)
    Command(s, false)
}

func EndFrame() {

    // This tells tkinter to update its idle tasks so that "sprites" actually get drawn.
    // It also waits for a response via stdout before continuing, which is helpful to
    // prevent the Golang program from getting ahead of tkinter.

    Command("ENDFRAME", true)
}

func Start(width int, height int, directory string, bg string) error {

    w := fmt.Sprintf("%d", width)
    h := fmt.Sprintf("%d", height)

    _, filename, _, _ := runtime.Caller(0)
    dir := filepath.Dir(filename)
    renderer := filepath.Join(dir, RENDERER)

    exec_command := exec.Command(
        "python", renderer, "--width", w, "--height", h, "--directory", directory, "--bg", bg)

    stdin_pipe, _ = exec_command.StdinPipe()
    stdout_pipe, _ = exec_command.StdoutPipe()
    stderr_pipe, _ = exec_command.StderrPipe()

    go pipe_relay(stdout_pipe, stdout_chan, nil)
    go pipe_relay(stderr_pipe, nil, os.Stderr)

    err := exec_command.Start()
    if err != nil {
        return fmt.Errorf("Start(): %v", err)
    }

    return nil
}
