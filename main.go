package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var log *logrus.Logger
var oldOrg string
var newOrg string
var dryRun bool

func walk(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d.IsDir() {
		return nil
	}

	configFileName := filepath.Base(path)
	if configFileName != "config" {
		return nil
	}

	dir := filepath.Dir(path)

	gitDirName := filepath.Base(dir)
	if gitDirName != ".git" {
		return nil
	}

	log.Infof("Analyzing '%s'\n", path)

	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("os.ReadFile error: %s\n", fmt.Errorf("%w", err))
	}

	content := string(bytes)

	// Parse and replace git config 'url' tags
	var regexp1 = regexp.MustCompile(
		fmt.Sprintf(`(\s*url\s*=\s*(https:\/\/github\.com\/)|(git@github\.com:))%s(/\S+\.git)`, oldOrg))
	fixedContent := regexp1.ReplaceAllString(content,
		fmt.Sprintf(`${1}%s${4}`, newOrg))

	// Parse and replace git config 'lfs' tags
	var regexp2 = regexp.MustCompile(
		fmt.Sprintf(`(\s*\[\s*lfs\s*"(https:\/\/github\.com\/)|(git@github\.com:))%s(/\S+\.git\S+\s*"\s*])`, oldOrg))
	fixedContent = regexp2.ReplaceAllString(fixedContent,
		fmt.Sprintf(`${1}%s${4}`, newOrg))

	results := compare(content, fixedContent)

	if !dryRun && len(results) != 0 {
		err := os.WriteFile(path, []byte(fixedContent), fs.ModePerm)
		if err != nil {
			log.Fatalf("os.WriteFile error: %s\n", fmt.Errorf("%w", err))
		}
	}

	for i := 0; i < len(results); i++ {
		log.Infof("Changed '%s'\n", results[i].after)
	}

	return nil
}

type comparisonResult struct {
	before string
	after  string
}

func compare(contentBefore string, contentAfter string) []comparisonResult {
	linesBefore := strings.Split(contentBefore, "\n")
	linesAfter := strings.Split(contentAfter, "\n")

	linesBeforeLength := len(linesBefore)
	linesAfterLength := len(linesAfter)

	if linesBeforeLength != linesAfterLength {
		log.Fatalf("file content can't be parsed properly\n")
	}

	var result []comparisonResult
	for i := 0; i < linesAfterLength; i++ {
		if strings.Compare(linesBefore[i], linesAfter[i]) != 0 {
			result = append(result, comparisonResult{before: linesBefore[i], after: linesAfter[i]})
		}
	}

	return result
}

func main() {
	log = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%] %time% - %msg%",
		},
	}

	dryRunPtr := flag.Bool("dryRun", true, "dry run (don't change actual files)")
	directoryPtr := flag.String("directory", ".", "parent directory")
	oldOrgPtr := flag.String("oldOrg", "PerkinElmer", "old organization name")
	newOrgPtr := flag.String("newOrg", "PerkinElmerAES", "new organization name")
	flag.Parse()

	log.Infof("Flag dryRun assigned to '%t'\n", *dryRunPtr)
	log.Infof("Flag directory assigned to '%s'\n", *directoryPtr)
	log.Infof("Flag oldOrg assigned to '%s'\n", *oldOrgPtr)
	log.Infof("Flag newOrg assigned to '%s'\n", *newOrgPtr)
	oldOrg = *oldOrgPtr
	newOrg = *newOrgPtr
	dryRun = *dryRunPtr

	fileInfo, err := os.Stat(*directoryPtr)
	if err != nil {
		log.Fatalf("os.Stat error: %s\n", fmt.Errorf("%w", err))
	}

	if !fileInfo.IsDir() {
		log.Fatalf("directory command line flag should be a folder\n")
	}

	log.Infof("===============================================\n")
	if dryRun {
		log.Infof("This is DRY RUN only! No data will be modified!\n")
	} else {
		log.Infof("Your local GitHub repositories will be updated!\n")
	}
	log.Infof("===============================================\n")

	err = filepath.WalkDir(*directoryPtr, walk)
	if err != nil {
		log.Fatalf("filepath.WalkDir error: %s\n", fmt.Errorf("%w", err))
	}
}
