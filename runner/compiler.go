package runner

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hoangvvo/hcmut-co1027/conf"
	cp "github.com/otiai10/copy"
)

type Result struct {
	Name           string
	Error          string
	ResultExpected string
	ResultGot      string
	Dir            string
}

func preformat(output string) string {
	output = strings.TrimSpace(output)
	output = strings.Trim(output, "\n")
	return output
}

func DeleteCompiled(runDir string) error {
	return os.RemoveAll(runDir)
}

func Run(runDir string, caseDirs []string) ([]Result, error) {
	l := log.New(os.Stdout, "run/runDir: ", log.Lshortfile)

	l.Println("running " + fmt.Sprint(len(caseDirs)) + " cases")

	_, err := os.Stat(runDir)
	if err != nil {
		l.Println(err)
		return nil, err
	}

	var results []Result
	for _, caseDir := range caseDirs {
		caseName := filepath.Base(caseDir)

		// copy content of caseDir to runDir
		if err = cp.Copy(caseDir, runDir, cp.Options{
			OnDirExists: func(src, dest string) cp.DirExistsAction {
				return cp.Replace
			},
		}); err != nil {
			l.Println(err)
			return nil, err
		}

		result := Result{Name: caseName}

		testOutput, err := os.ReadFile(filepath.Join(caseDir, "output.txt"))
		if err != nil {
			l.Println(err)
			return nil, err
		}
		result.ResultExpected = preformat(string(testOutput))

		// run file
		var stdErr bytes.Buffer
		cmdRun := exec.Command("./main")
		cmdRun.Dir = runDir
		cmdRun.Stderr = &stdErr

		output, err := cmdRun.Output()
		if err != nil {
			result.Error = fmt.Errorf("%w: %s", err, stdErr.String()).Error()
		} else {
			result.ResultGot = preformat(string(output))
			// compare output
			if strings.Compare(string(result.ResultGot), string(result.ResultExpected)) != 0 {
				result.Error = ErrResultMismatch.Error()
			}
		}

		results = append(results, result)
	}

	return results, nil
}

type CompileResult struct {
	RunDir    string
	SuiteName string
}

func Compile(answer string, suiteName string) (*CompileResult, error) {
	suiteName = fileNameWithoutExtSliceNotation(filepath.Base(suiteName)) // prevent malicious paths
	suiteDir := filepath.Join(conf.SuitesDir, suiteName)

	stat, err := os.Stat(suiteDir)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, errors.New("not a directory")
	}

	wd, err := os.MkdirTemp(os.TempDir(), "run")
	if err != nil {
		return nil, err
	}

	mainCppPath := filepath.Join(wd, "main.cpp")
	err = cp.Copy(filepath.Join(suiteDir, "main.cpp"), mainCppPath)
	if err != nil {
		return nil, err
	}

	studyInPinkHPath := filepath.Join(wd, "studyInPink2.h")
	if err = os.WriteFile(studyInPinkHPath, []byte(answer), 0644); err != nil {
		return nil, err
	}
	// compile file
	var stderr bytes.Buffer
	cmd := exec.Command("g++", "-std=c++11", "main.cpp", "-o", "main")
	cmd.Dir = wd
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		return nil, fmt.Errorf("%w: %s", err, stderr.String())
	}

	// remove source files
	if err = os.Remove(mainCppPath); err != nil {
		return nil, err
	}
	if err = os.Remove(studyInPinkHPath); err != nil {
		return nil, err
	}

	return &CompileResult{
		RunDir:    wd,
		SuiteName: suiteName,
	}, nil
}
