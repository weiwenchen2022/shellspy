package shellspy_test

import (
	"io"
	"strings"
	"testing"

	"shellspy"

	"github.com/google/go-cmp/cmp"
)

func TestCommandFromString(t *testing.T) {
	t.Parallel()

	input := "echo hello"
	cmd, err := shellspy.CommandFromString(input)
	if err != nil {
		t.Fatal(err)
	}

	wantArgs := []string{"echo", "hello"}
	gotArgs := cmd.Args
	if diff := cmp.Diff(wantArgs, gotArgs); diff != "" {
		t.Error(diff)
	}
}

func TestScan(t *testing.T) {
	t.Parallel()

	r := strings.NewReader("echo\nexit\n")
	err := shellspy.Scan(r, io.Discard, io.Discard)
	if err != nil {
		t.Error(err)
	}

	if r.Len() > 0 {
		t.Error()
	}
}

func TestExecCommand(t *testing.T) {
	t.Parallel()

	input := "echo -n hello world"
	cmd, err := shellspy.CommandFromString(input)
	if err != nil {
		t.Error(err)
	}

	gotOut, err := shellspy.ExecCmd(cmd)
	if err != nil {
		t.Error(err)
	}

	wantOut := "hello world"
	if diff := cmp.Diff(wantOut, gotOut); diff != "" {
		t.Error(diff)
	}
}

// func TestParse(t *testing.T) {
// 	input := `sh -c 'echo stdout; echo stderr 1>&2'`
// 	fmt.Println(input)

// 	input = strconv.Quote(input)
// 	fmt.Println(input)

// 	tokerns := strings.Split(input, " ")
// 	fmt.Println(len(tokerns))
// }
