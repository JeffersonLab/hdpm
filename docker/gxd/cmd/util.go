package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func workDir() string {
	pdir := os.Getenv("DOCKER")
	if pdir == "" {
		pdir, _ = os.Getwd()
	}
	return pdir
}

func write_text(fname, text string) {
	f, err := os.Create(fname)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintln(f, text)
	f.Close()
}

func in(args []string, item string) bool {
	for _, arg := range args {
		if item == arg {
			return true
		}
	}
	return false
}

func glob(pattern string) []string {
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalln(err)
	}
	return files
}

func rmGlob(pattern string) {
	files := glob(pattern)
	for _, file := range files {
		os.RemoveAll(file)
	}
}

func run(name string, args ...string) {
	if err := command(name, args...).Run(); err != nil {
		log.Fatalln(err)
	}
}

func command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func output(name string, args ...string) string {
	c := exec.Command(name, args...)
	b, err := c.CombinedOutput()
	if err != nil {
		log.Fatalln(err)
	}
	return strings.TrimRight(string(b), "\n")
}
