package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// return types
func TestNewConfig(t *testing.T) {
	cfg, err := newConfig()
	if err != nil {
		t.Fatal(err)
	}

	isDefault(t, cfg.TransactionsLogPath, filename_transactions_log, "TransactionsLogPath")
	isDefault(t, cfg.StatusPath, filename_status, "StatusPath")
	isDefault(t, cfg.CtlFilePath, filename_ctlfile, "CtlFilePath")
}

func TestParseLine(t *testing.T) {
	/* test case setup */
	var abspaths []string
	relpaths := []string{
		"../foo.txt",
		"./bar.json",
		"++z",
		"./",
	}

	for _, rel := range relpaths {
		p, err := filepath.Abs(rel)
		if err != nil {
			t.Fatal(err)
		}
		abspaths = append(abspaths, p)
	}

	testCases := [][]string{
		[]string{"foo=bar", "FAIL"},
		[]string{"\n", "DEFAULT"},
		[]string{"", "DEFAULT"},
		[]string{"transactionslog=../foo.txt", abspaths[0], transactionslog},
		[]string{"statusfile=./bar.json", abspaths[1], statusfile},
		[]string{"ctlfile=++z", abspaths[2], ctlfile},
		[]string{"ctlfile=", "FAIL"},
		[]string{"transactionslog=/tmp/abc.txt", "/tmp/abc.txt", transactionslog},
		[]string{"ctlfile=./", abspaths[3], ctlfile},
		[]string{"abc", "FAIL"},
	}
	/*******************/

	for _, tc := range testCases {
		cfg, err := newConfig()
		if err != nil {
			t.Fatal(err)
		}

		err = parseLine(tc[0], cfg)

		switch tc[1] {
		case "DEFAULT":
			if err != nil {
				t.Error("this test case shouldn't return an error but it did: [", err, "]")
			}
			isDefault(t, cfg.TransactionsLogPath, filename_transactions_log, "TransactionsLogPath")
			isDefault(t, cfg.StatusPath, filename_status, "StatusPath")
			isDefault(t, cfg.CtlFilePath, filename_ctlfile, "CtlFilePath")
			break

		case "FAIL":
			if err == nil {
				t.Error("this test case should return an error")
			}
			break

		default:
			if err != nil {
				t.Error("this test case shouldn't return an error but it did: [", err, "]")
			}

			var savedPath string

			switch tc[2] {
			case transactionslog:
				savedPath = cfg.TransactionsLogPath

			case statusfile:
				savedPath = cfg.StatusPath

			case ctlfile:
				savedPath = cfg.CtlFilePath
			}

			if savedPath != tc[1] {
				t.Error("the path is", savedPath, "but should be", tc[1])
			}
		}
	}
}

func isDefault(t *testing.T, configPath string, defaultFilename string, attrName string) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	completePath := wd + "/" + defaultFilename
	if configPath != completePath {
		t.Error(attrName + "is" + configPath + "and it should be" + completePath)
	}
}

func TestParseConfig(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	/* test cases setup */
	testCases := [][]string{
		[]string{"transactionslog=/tmp/abc.txt",
			"/tmp/abc.txt",
			wd + "/" + filename_status,
			wd + "/" + filename_ctlfile},
		[]string{"statusfile=/foo/bar\n\ntransactionslog=/tmp/foo/bar.txt\n  \nctlfile=/tmp/x",
			"/tmp/foo/bar.txt",
			"/foo/bar",
			"/tmp/x"},
		[]string{"foobar\n\n", "FAIL"},
		[]string{"\n\nfoobar", "FAIL"},
	}
	// fail line list
	failList := [4]int{}
	failList[2] = 1
	failList[3] = 3
	/*******************/

	for i, tc := range testCases {
		t.Log("Testing [", tc[0], "]")

		r := bufio.NewReader(strings.NewReader(tc[0]))
		cfg, err := parseConfig(r)

		switch tc[1] {
		case "FAIL":
			if err == nil {
				t.Error("This should have failed but it did not")
				break
			}

			e, ok := err.(ErrCfgFormat)
			if !ok {
				t.Error("This should have return a ErrCfgFormat but return another type: [", err, "]")
				break
			}
			if int(e) != failList[i] {
				t.Error("Count line was", int(e), "but it should have been", failList[i])
			}

		default:
			if err != nil {
				t.Error("This shouldn't have return an error but it did: [", err, "]")
				break
			}

			if cfg.TransactionsLogPath != tc[1] {
				t.Error("TransactionsLogPath is", cfg.TransactionsLogPath, "but should be", tc[1])
			} else if cfg.StatusPath != tc[2] {
				t.Error("StatusPath is", cfg.StatusPath, "but should be", tc[2])
			} else if cfg.CtlFilePath != tc[3] {
				t.Error("CtlFilePath is", cfg.CtlFilePath, "but should be", tc[3])
			}
		}
	}
}
