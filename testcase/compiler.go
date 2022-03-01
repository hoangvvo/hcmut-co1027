package testcase

import (
	"bytes"
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	cp "github.com/otiai10/copy"
)

//go:embed main.cpp.txt
var mainCpp []byte

type Result struct {
	Name       string
	Error      string
	TestOutput string
	MyOutput   string
}

func CompileAndRun(caseDirs []string, answer string) ([]Result, error) {
	var results []Result

	wd, err := ioutil.TempDir(os.TempDir(), "testrun")
	defer os.Remove(wd)
	if err != nil {
		return results, err
	}

	// write main.cpp and studyInPink2
	if err = os.WriteFile(filepath.Join(wd, "main.cpp"), mainCpp, 0644); err != nil {
		return results, err
	}

	if err = os.WriteFile(filepath.Join(wd, "studyInPink2.h"), []byte(answer), 0644); err != nil {
		return results, err
	}
	// compile file
	var stderr bytes.Buffer
	cmd := exec.Command("g++", "main.cpp", "-o", "main")
	cmd.Dir = wd
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		return results, fmt.Errorf("%w: %s", err, stderr.String())
	}

	for _, caseDir := range caseDirs {
		caseName := filepath.Base(caseDir)
		// copy test case
		if err = cp.Copy(caseDir, wd, cp.Options{
			OnDirExists: func(src, dest string) cp.DirExistsAction {
				return cp.Replace
			},
		}); err != nil {
			return results, err
		}

		// run file
		var stderrRun bytes.Buffer
		cmdRun := exec.Command("./main")
		cmdRun.Dir = wd
		cmdRun.Stderr = &stderrRun

		result := Result{Name: caseName}

		testOutput, err := os.ReadFile(filepath.Join(caseDir, "result.txt"))
		if err != nil {
			return results, err
		}

		result.TestOutput = string(testOutput)

		output, err := cmdRun.Output()
		if err != nil {
			result.Error = fmt.Errorf("%w: %s", err, stderr.String()).Error()
		} else {
			result.MyOutput = string(output)
			if string(result.MyOutput) != string(result.TestOutput) {
				result.Error = ErrResultMismatch.Error()
			}
		}
		results = append(results, result)
	}

	return results, nil
}
