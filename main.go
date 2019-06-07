package main

import (
	"bufio"
	"flag"
	"path/filepath"

	"github.com/eaciit/toolkit"
)

var (
	config     = new(Config)
	err        error
	configPath = flag.String("config", "./config", "path to config folder")
	output     = flag.String("output", "./output.csv", "Flag of output file")
	projectout = flag.String("project", "./project.csv", "project usage output")
	log        = getLogger("codeusage")
	w          *bufio.Writer
)

func main() {
	flag.Parse()

	if config, err = readConfig(*configPath); err != nil {
		log.Errorf("unable to read config. %s", err.Error())
		return
	}
	log.Infof("reading from %s\n\t\t%s", *configPath, toolkit.JsonString(config))

	/*
		fo, err := os.OpenFile(*output, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Error(err.Error())
			return
		}
		defer fo.Close()

		w = bufio.NewWriter(fo)
		defer w.Flush()

		w.WriteString("name,path,line\n")

		for _, p := range config.Libraries {
			formatedPath := p
			if formatedPath[0] != '/' {
				formatedPath = filepath.Join(config.WorkingDir, p)
			}
			if err := readLibraries(formatedPath, *output); err != nil {
				log.Error(err.Error())
			}
		}
	*/

	for _, p := range config.Projects {
		formatedPath := p
		if formatedPath[0] != '/' {
			formatedPath = filepath.Join(config.WorkingDir, p)
		}

		if err := readProjects(p, formatedPath, *projectout, nil); err != nil {
			log.Error(err.Error())
		}
	}
}
