package runner

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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

type CompileResult struct {
	RunDir    string
	SuiteName string
}

func preformat(output string) string {
	output = strings.TrimSpace(output)
	output = strings.Trim(output, "\n")
	return output
}

func DeleteCompiled(runDir string) error {
	return os.RemoveAll(runDir)
}

var lRun = log.New(os.Stdout, "run: ", log.LstdFlags)

func Run(runDir string, caseDirs []string) ([]Result, error) {
	_, err := os.Stat(runDir)
	if err != nil {
		lRun.Println(err)
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
			lRun.Println(err)
			return nil, err
		}

		result := Result{Name: caseName}

		testOutput, err := os.ReadFile(filepath.Join(caseDir, "output.txt"))
		if err != nil {
			lRun.Println(err)
			return nil, err
		}
		result.ResultExpected = preformat(string(testOutput))

		// run file
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.TimeoutSec)*time.Second)
		defer cancel()

		var stdErr bytes.Buffer
		cmdRun := exec.CommandContext(ctx, "./main")
		cmdRun.Dir = runDir
		cmdRun.Stderr = &stdErr

		output, err := cmdRun.Output()
		if ctx.Err() == context.DeadlineExceeded {
			// this means the user takes too long on one task
			// and may take a long time for later tasks
			// we consider just stop
			return nil, ErrDeadlineExceeded
		}

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

func Compile(answer string, answerFileName string, suiteName string) (*CompileResult, error) {
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

	studyInPinkHPath := filepath.Join(wd, answerFileName)
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
