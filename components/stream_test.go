package components

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/clakeboy/golib/utils"
	"github.com/creack/pty"
	"golang.org/x/term"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"testing"
)

func TestBinToBytes(t *testing.T) {
	bye, err := hex.DecodeString("FF2C2E1F")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(bye)

	fmt.Println(0xff2e)

	protocol := []byte{0x03, 0xEE, 0xFE, 0x02}
	fmt.Println(bytes.Compare(MainProtocol, protocol))
	fmt.Println(utils.IntToBytes(CMDClose, 16))
}

func TestMainStream_Build(t *testing.T) {
	cmd := NewMainStream()
	cmd.Command = CMDClose
	cmdStream := cmd.Build()
	fmt.Println(hex.EncodeToString(cmdStream))

	deCmd := NewMainStream()
	err := deCmd.Parse(cmdStream)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestRedCode(t *testing.T) {
	test()
}

func test() error {
	// Create arbitrary command.
	c := exec.Command("bash")

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)

	return nil
}
