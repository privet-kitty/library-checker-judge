package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

func TestExecutorHello(t *testing.T) {
	cmd := exec.Command("echo", "Hello")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	result, err := SafeRun(cmd, 1.0, false)
	if err != nil {
		t.Fatalf("Fail Execute: %v", err)
	}
	if result.ReturnCode != 0 {
		t.Errorf("Error return code: %v", result.ReturnCode)
	}
	if 0.5 < result.Time {
		t.Error("Comsume too long time for Hello")
	}
}

func TestExecutorTimeOut(t *testing.T) {
	cmd := exec.Command("sleep", "5")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	result, err := SafeRun(cmd, 1.0, false)
	if err != nil {
		t.Fatal("Error ", err)
	}
	if !result.Tle {
		t.Fatal("Error not tle")
	}
	if result.Time < 0.9 || 1.1 < result.Time {
		t.Fatal("Error result time = ", result.Time)
	}
}

func TestOutputStripper(t *testing.T) {
	shortStr := []byte("short string")

	os := &outputStripper{N: 100}
	_, err := os.Write(shortStr)
	if err != nil {
		t.Fatal("outputStripper Error ", err)
	}
	res := os.Bytes()
	if !bytes.Equal(shortStr, res) {
		t.Fatal("outputStripper Differ")
	}
}

func TestOutputStripperLong(t *testing.T) {
	longStrBase := []byte("long string")
	longStr := []byte{}
	for i := 0; i < 100; i++ {
		longStr = append(longStr, longStrBase...)
	}

	os := &outputStripper{N: 100}
	_, err := os.Write(longStr)
	if err != nil {
		t.Fatal("outputStripper Error ", err)
	}
	res := os.Bytes()
	if len(res) > 100 {
		t.Fatal("outputStripper Differ")
	}
	t.Log(string(res))
}

func TestExecutorInfinityCE(t *testing.T) {
	checker, err := os.Open("./test_src/aplusb/checker.cpp")
	if err != nil {
		t.Fatal("Failed: Checker", err)
	}
	src, err := os.Open("./test_src/many_ce.d")
	if err != nil {
		t.Fatal("Failed: Source", err)
	}
	tempdir, err := ioutil.TempDir("", "judge")
	if err != nil {
		t.Fatal("Failed: tempdir", err)
	}
	defer os.RemoveAll(tempdir)
	judge, err := NewJudge(tempdir, "d", checker, src, 2.0)
	if err != nil {
		t.Fatal("Failed: NewJudge", err)
	}

	result, err := judge.CompileSource()

	if err != nil {
		t.Fatal("Failed: Failed Compile", err)
	}
	if result.ReturnCode == 0 {
		t.Fatal("Failed: Must be CE")
	}
	t.Log(string(result.Stderr))
}

func generateAplusB(t *testing.T, lang, srcName string) *Judge {
	checker, err := os.Open("./test_src/aplusb/checker.cpp")
	if err != nil {
		t.Fatal("Failed: Checker", err)
	}
	src, err := os.Open(path.Join("test_src/aplusb", srcName))
	if err != nil {
		t.Fatal("Failed: Source", err)
	}
	tempdir, err := ioutil.TempDir("", "judge")
	if err != nil {
		t.Fatal("Failed: tempdir", err)
	}
	judge, err := NewJudge(tempdir, lang, checker, src, 2.0)
	if err != nil {
		t.Fatal("Failed: NewJudge", err)
	}

	result, err := judge.CompileChecker()
	if err != nil || result.ReturnCode != 0 {
		t.Fatal("error CompileChecker", err, string(result.Stderr))
	}
	result, err = judge.CompileSource()
	if err != nil || result.ReturnCode != 0 {
		t.Fatal("error CompileSource", err, string(result.Stderr))
	}

	return judge
}

func TestAplusbAC(t *testing.T) {
	judge := generateAplusB(t, "cpp", "ac.cpp")
	in := strings.NewReader("1 1")
	expect := strings.NewReader("2")
	result, err := judge.TestCase(in, expect)
	log.Println(judge.dir)
	if err != nil {
		t.Fatal("error Run Test", err)
	}
	if result.Status != "AC" {
		t.Fatal("error Status", result)
	}
}

func TestAplusbRustAC(t *testing.T) {
	judge := generateAplusB(t, "rust", "ac.rs")
	in := strings.NewReader("1 1")
	expect := strings.NewReader("2")
	result, err := judge.TestCase(in, expect)
	log.Println(judge.dir)
	if err != nil {
		t.Fatal("error Run Test", err)
	}
	if result.Status != "AC" {
		t.Fatal("error Status", result)
	}
}

func TestAplusbHaskellAC(t *testing.T) {
	judge := generateAplusB(t, "haskell", "ac.hs")
	in := strings.NewReader("1 1")
	expect := strings.NewReader("2")
	result, err := judge.TestCase(in, expect)
	log.Println(judge.dir)
	if err != nil {
		t.Fatal("error Run Test", err)
	}
	if result.Status != "AC" {
		t.Fatal("error Status", result)
	}
}

func TestAplusbCSharpAC(t *testing.T) {
	judge := generateAplusB(t, "csharp", "ac.cs")
	in := strings.NewReader("1 1")
	expect := strings.NewReader("2")
	result, err := judge.TestCase(in, expect)
	log.Println(judge.dir)
	if err != nil {
		t.Fatal("error Run Test", err)
	}
	if result.Status != "AC" {
		t.Fatal("error Status", result)
	}
}

func TestAplusbWA(t *testing.T) {
	judge := generateAplusB(t, "cpp", "wa.cpp")
	in := strings.NewReader("1 1")
	expect := strings.NewReader("2")
	result, err := judge.TestCase(in, expect)
	if err != nil {
		t.Fatal("error Run Test", err)
	}
	if result.Status != "WA" {
		t.Fatal("error Status", result)
	}
}

func TestAplusbPE(t *testing.T) {
	judge := generateAplusB(t, "cpp", "pe.cpp")
	in := strings.NewReader("1 1")
	expect := strings.NewReader("2")
	result, err := judge.TestCase(in, expect)
	if err != nil {
		t.Fatal("error Run Test", err)
	}
	if result.Status != "PE" {
		t.Fatal("error Status", result)
	}
}

func TestAplusbFail(t *testing.T) {
	judge := generateAplusB(t, "cpp", "ac.cpp")
	in := strings.NewReader("1 1")
	expect := strings.NewReader("3") // !?
	result, err := judge.TestCase(in, expect)
	if err != nil {
		t.Fatal("error Run Test", err)
	}
	if result.Status != "Fail" {
		t.Fatal("error Status", result)
	}
}
