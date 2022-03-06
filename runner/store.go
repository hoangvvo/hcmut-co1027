package runner

import (
	"errors"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hoangvvo/hcmut-co1027/conf"
	"github.com/hoangvvo/hcmut-co1027/unzip"
)

type TestSuite struct {
	Name      string
	Total     int
	TestCases []string
}

func GetSuite(suiteName string) (*TestSuite, error) {
	suiteDir := filepath.Join(conf.SuitesDir, suiteName)
	entries, err := os.ReadDir(suiteDir)
	if err != nil {
		return nil, err
	}
	ts := TestSuite{
		Name:      suiteName,
		TestCases: []string{},
	}
	for _, entry := range entries {
		if entry.IsDir() {
			ts.TestCases = append(ts.TestCases, entry.Name())
		}
	}
	ts.Total = len(ts.TestCases)
	return &ts, nil
}

func GetSuites() ([]TestSuite, error) {
	entries, err := os.ReadDir(conf.SuitesDir)
	if err != nil {
		return nil, err
	}
	var results []TestSuite
	for _, entry := range entries {
		dirName := entry.Name()

		dirsForCnt, _ := os.ReadDir(filepath.Join(conf.SuitesDir, dirName))
		total := len(dirsForCnt) - 1

		results = append(results, TestSuite{
			Name:  dirName,
			Total: total,
		})
	}
	return results, nil
}

func DeleteSuite(suiteName string) error {
	suiteName = filepath.Base(suiteName) // prevent malicious paths
	if suiteName == "case-1" {
		return errors.New("không được xóa cái này :)")
	}

	suiteDir := filepath.Join(conf.SuitesDir, suiteName)
	suiteZip := filepath.Join(conf.ArchiveDir, suiteName+".zip")

	stat, err := os.Stat(suiteDir)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return errors.New("not a directory")
	}

	os.Remove(suiteZip)
	err = os.RemoveAll(suiteDir)

	if err != nil {
		log.Println("removed test case: " + suiteName)
	}

	return err
}

func AddSuite(fileName string, file multipart.File) error {
	suiteName := fileNameWithoutExtSliceNotation(filepath.Base(fileName))
	suiteDir := filepath.Join(conf.SuitesDir, suiteName)

	existed, err := exists(suiteDir)
	if err != nil {
		return err
	}
	if existed {
		return errors.New("test suite with with the name " + suiteName + " has already existed")
	}

	// add zip file
	zipPath := filepath.Join(conf.ArchiveDir, fileName)
	if err = writeFile(zipPath, file); err != nil {
		os.Remove(zipPath)
		os.RemoveAll(suiteDir)
		return err
	}

	// unzip
	if _, err = unzip.New().Extract(zipPath, suiteDir); err != nil {
		// fail, remove all
		os.Remove(zipPath)
		os.RemoveAll(suiteDir)
		return err
	}

	existed, err = exists(filepath.Join(suiteDir, "main.cpp"))
	if err != nil {
		// fail, remove all
		os.Remove(zipPath)
		os.RemoveAll(suiteDir)
		return err
	}
	if !existed {
		// fail, remove all
		os.Remove(zipPath)
		os.RemoveAll(suiteDir)
		return errors.New("main.cpp is missing")
	}

	// attempt to convert crlf to lf
	// loop each dirs because arg can be very long for dos2unix
	dirEntries, err := os.ReadDir(suiteDir)
	if err != nil {
		// fail, remove all
		os.Remove(zipPath)
		os.RemoveAll(suiteDir)
		return err
	}

	// including main.cpp file
	if len(dirEntries) < 2 {
		// fail, remove all
		os.Remove(zipPath)
		os.RemoveAll(suiteDir)
		return errors.New("require at least 1 test case")
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			exec.Command("bash", "-c", "dos2unix -o "+filepath.Join(suiteDir, entry.Name(), "*.txt")).Run()
		}
	}

	log.Println("added test case: " + suiteName)

	return nil
}
