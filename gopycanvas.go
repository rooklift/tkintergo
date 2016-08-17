package gopycanvas

// Comments:
//
// - We send the renderer a bunch of sprites to render, followed by "ENDFRAME".
// - The renderer sends "ENDFRAME" back as an acknowledgement.
// - It also sends key presses and key releases via Stderr.

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
    "strings"
    "sync"
)

const RENDERER = "tkinter_renderer.py"

var keymap = make(map[string]bool)
var keymap_MUTEX sync.Mutex

var stdin_pipe io.WriteCloser
var stdout_pipe io.ReadCloser
var stderr_pipe io.ReadCloser

func stderr_watcher() {

    scanner := bufio.NewScanner(stderr_pipe)

    for scanner.Scan() {

        s := scanner.Text()

        if strings.HasPrefix(s, "KEY") && len(s) > 4 {
            sym := s[4:]
            keymap_MUTEX.Lock()
            keymap[sym] = true
            keymap_MUTEX.Unlock()
        } else if strings.HasPrefix(s, "REL") && len(s) > 4 {
            sym := s[4:]
            keymap_MUTEX.Lock()
            keymap[sym] = false
            keymap_MUTEX.Unlock()
        } else if s == "QUIT" {
            fmt.Fprintf(os.Stderr, "%s has quit\n", RENDERER)
            os.Exit(0)
        } else {
            fmt.Fprintf(os.Stderr, "(renderer) " + s + "\n")
        }
    }
}

func KeyState(sym string) bool {
    keymap_MUTEX.Lock()
    ret := keymap[sym]      // Will be false if sym is not in the keymap at all
    keymap_MUTEX.Unlock()
    return ret
}

func Command(msg string, need_ack bool) {

    // It's important that need_ack is only true for commands that the renderer will actually acknowledge

    msg = strings.TrimSpace(msg)
    fmt.Fprintf(stdin_pipe, msg + "\n")

    if need_ack {
        scanner := bufio.NewScanner(stdout_pipe)
        scanner.Scan()
    }
}

func Sprite(filename string, x int32, y int32) {
    s := fmt.Sprintf("%s %d %d", filename, x, y)
    Command(s, false)
}

func EndFrame() {

    // This tells the renderer to draw the frame. (In the case of the Tkinter renderer,
    // this means calling update_idletasks()). It also waits for a response via stdout
    // before continuing, which is helpful to prevent the Golang program from getting
    // ahead of the renderer.

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

    go stderr_watcher()

    err := exec_command.Start()
    if err != nil {
        return fmt.Errorf("Start(): %v", err)
    }

    return nil
}
