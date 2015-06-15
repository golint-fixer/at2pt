package convert

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestWritePlainText(t *testing.T) {

	inf, err := os.Open("knp_test/input.knp")
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}
	defer inf.Close()

	//make output
	//cf: http://stackoverflow.com/questions/10473800/
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = ConvertKNP(inf, os.Stdout, true)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	gotStdout := <-outC

	//get golds
	goldf, err := os.Open("knp_test/gold.txt")
	gold_reader := bufio.NewReader(goldf)
	if err != nil {
		t.Errorf("Error when open the gold file: %v", err)
		return
	}

	//Check
	for _, line := range strings.Split(gotStdout, "\n") {
		_gold_line, _, err := gold_reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Errorf("Error when getting gold line: %s", err)
			return
		}

		gold_line := string(_gold_line)
		if line != gold_line {
			t.Errorf("got\n[%s]\nbut want\n[%s]\n", line, gold_line)
		}
	}
}
