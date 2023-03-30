package shellspy

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Session struct {
	input  io.Reader
	output io.Writer

	transcript io.Writer
}

func NewSession(input io.Reader, output io.Writer) *Session {
	return &Session{input: input, output: output}
}

func (s *Session) Run() error {
	f, err := os.OpenFile("shellspy.txt",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	s.transcript = f

	mw := io.MultiWriter(s.output, s.transcript)

	fmt.Fprintln(mw, "Recording session to 'shellspy.txt'")
	err = Scan(s.input, s.output, s.transcript)
	fmt.Fprintln(mw, "Session saved as 'shellspy.txt'")

	if err := f.Sync(); err != nil {
		log.Print(err)
	}

	return err
}

func CommandFromString(input string) (*exec.Cmd, error) {
	tokens := strings.Split(input, " ")
	return exec.Command(tokens[0], tokens[1:]...), nil
}

func Scan(input io.Reader, output io.Writer, transcript io.Writer) error {
	sc := bufio.NewScanner(input)

	mw := io.MultiWriter(output, transcript)

	fmt.Fprint(mw, "\n> ")
	for sc.Scan() {
		line := sc.Text()
		fmt.Fprintln(transcript, line)

		if line == "exit" {
			break
		}

		cmd, err := CommandFromString(line)
		if err != nil {
			return err
		}

		stdoutStderr, err := ExecCmd(cmd)
		if err != nil {
			return err
		}

		fmt.Fprint(mw, stdoutStderr)
		fmt.Fprint(mw, "\n> ")
	}

	return sc.Err()
}

func ExecCmd(cmd *exec.Cmd) (string, error) {
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return out.String(), nil
}
