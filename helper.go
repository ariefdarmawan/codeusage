package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func readLibraries(input, output string) error {
	fis, err := ioutil.ReadDir(input)
	if err != nil {
		return fmt.Errorf("unable to read libraries. %s", err.Error())
	}

	for _, fi := range fis {
		if fi.IsDir() {
			name := fi.Name()
			if name[0] != '.' && name[len(name)-1] != '_' {
				pathName := filepath.Join(input, name)
				if lineCount, err := processLibrary(name, pathName, output); err == nil {
					log.Infof("Library %s (%s) = %d", name, pathName, lineCount)
					w.WriteString(fmt.Sprintf("%s,%s,%d\n", name, pathName, lineCount))
				} else {
					return fmt.Errorf("unable to process library %s. %s", pathName, err.Error())
				}
			}
		}
	}

	return nil
}

func processLibrary(name, folderPath, output string) (int, error) {
	fis, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return 0, err
	}

	totalLines := 0
	for _, fi := range fis {
		fiName := fi.Name()
		if fiName[0] != '.' {
			lineCount := 0
			fiPath := filepath.Join(folderPath, fiName)
			if fi.IsDir() {
				if lineCount, err = processLibrary(name, fiPath, output); err != nil {
					return 0, err
				}
			} else {
				bs, err := ioutil.ReadFile(fiPath)
				if err == nil {
					lineCount = len(strings.Split(string(bs), "\n"))
				} else {
					return 0, err
				}
			}
			totalLines += lineCount
		}
	}

	return totalLines, nil
}

func readProjects(name, folderPath, output string, pw *bufio.Writer) error {
	fis, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}

	if pw == nil {
		f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("unable to create project usage output. %s", err.Error())
		}
		defer f.Close()

		pw = bufio.NewWriter(f)
		defer pw.Flush()

		pw.WriteString("project,file,line,uselib,libraries\n")
	}

	for _, fi := range fis {
		fiName := fi.Name()
		if fiName[0] != '.' {
			fiPath := filepath.Join(folderPath, fiName)
			if fi.IsDir() {
				if err = readProjects(name, fiPath, "", pw); err != nil {
					return err
				}
			} else {
				if fiPath[len(fiPath)-3:] == ".go" && !strings.Contains(fiPath, "vendor") {
					bs, err := ioutil.ReadFile(fiPath)
					useLib := false
					libs := []string{}
					if err == nil {
						lines := strings.Split(string(bs), "\n")
						lineCount := len(lines)
						for _, line := range lines {
							cleanedLine := strings.Trim(strings.Trim(line, " "), "\t")
							if len(cleanedLine) > 2 {
								if cleanedLine[0:2] != "//" && cleanedLine[0:2] != "/*" {
									for _, l := range config.Libraries {
										if strings.Contains(cleanedLine, l) {
											libNames := strings.Split(cleanedLine, "\"")
											if len(libNames) > 1 {
												useLib = true
												libs = append(libs, libNames[1])
											}
										}
									}
								}
							}
						}
						log.Infof("%s,%s,%d,%v,\"%s\"", name, fiPath, lineCount, useLib, strings.Join(libs, ", "))
						pw.WriteString(fmt.Sprintf("%s,%s,%d,%v,\"%s\"\n", name, fiPath, lineCount, useLib, strings.Join(libs, ", ")))
					} else {
						return err
					}
				}
			}
		}
	}

	return nil
}
