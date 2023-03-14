package logger

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var LoggerT *Logger

func mockPrint(buf *bytes.Buffer) {
	printf = func(format string, a ...interface{}) (int, error) {
		buf.WriteString(fmt.Sprintf(format, a...))
		return 0, nil
	}
}

func checkFile() string {
	file, err := os.Open("log.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		return scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return ""
}

func init() {
	testLogAvalable := []int{ERROR, INFO, WARN}
	LoggerT, _ = New("log.log", "Test", testLogAvalable, false)
}

func TestErrorPrintEmty(t *testing.T) {
	var buf bytes.Buffer
	mockPrint(&buf)
	testLogAvalable := []int{}
	LoggerTest, _ := New("log.log", "Test", testLogAvalable, false)
	LoggerTest.Error(buf)
	assert.Empty(t, buf)
}

func TestErrorPrintNotEmty(t *testing.T) {
	var buf bytes.Buffer
	LoggerT.Error("TEST")
	mockPrint(&buf)
	testLogAvalable := []int{}
	LoggerTest, err := New("log.log", "Test", testLogAvalable, false)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	LoggerTest.SetAvalableLogging(ERROR)
	LoggerTest.Error(buf)
	assert.NotEmpty(t, buf)
}

func TestWarnPrintNotEmty(t *testing.T) {
	var buf bytes.Buffer
	testLogAvalable := []int{}
	LoggerTest, _ := New("log.log", "Test", testLogAvalable, true)
	LoggerTest.SetAvalableLogging(WARN)
	LoggerTest.Warning(buf)
	str := checkFile()
	assert.NotEmpty(t, str)
}
