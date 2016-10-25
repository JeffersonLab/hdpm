package cmd

import (
	"fmt"
	"io/ioutil"
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

func ver_i(ver string, i int) string {
	if i < 0 || i > 2 {
		return ver
	}
	for _, v := range strings.Split(ver, "-") {
		if strings.Contains(v, ".") {
			return strings.Split(v, ".")[i]
		}
	}
	return ver
}

func write_text(fname, text string) {
	f, err := os.Create(fname)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintln(f, text)
	f.Close()
}

func extractNames(args []string) []string {
	var names []string
	for _, arg := range args {
		if strings.Contains(arg, "@") {
			names = append(names, strings.Split(arg, "@")[0])
		} else {
			names = append(names, arg)
		}
	}
	return names
}

func extractVersions(args []string) map[string]string {
	versions := make(map[string]string)
	unchanged := true
	for _, arg := range args {
		if strings.Contains(arg, "@") {
			versions[strings.Split(arg, "@")[0]] = strings.Split(arg, "@")[1]
			unchanged = false
		} else {
			versions[arg] = ""
		}
	}
	if unchanged {
		versions = nil
	}
	return versions
}

func extractVersion(arg string) string {
	ver := ""
	if strings.Contains(arg, "@") {
		ver = strings.Split(arg, "@")[1]
	}
	return ver
}

func readFile(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	return strings.TrimRight(string(b), "\n")
}

func in(args []string, item string) bool {
	for _, arg := range args {
		if item == arg {
			return true
		}
	}
	return false
}

func mk(path string) {
	if err := os.MkdirAll(path, 0777); err != nil {
		log.Fatalln(err)
	}
}

func cd(path string) {
	if err := os.Chdir(path); err != nil {
		log.Fatalln(err)
	}
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

func mkcd(path string) {
	mk(path)
	cd(path)
}

func isPath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func isSymLink(path string) bool {
	if !isPath(path) {
		return false
	}
	stat, err := os.Lstat(path)
	if err != nil {
		log.Fatalln(err)
	}
	return stat.Mode()&os.ModeSymlink == os.ModeSymlink
}

func readLink(path string) string {
	link, err := os.Readlink(path)
	if err != nil {
		log.Fatalln(err)
	}
	return link
}

func readDir(path string) []string {
	if !isPath(path) {
		return nil
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	var names []string
	for _, file := range files {
		names = append(names, file.Name())
	}
	return names
}

func contains(list []string, sublist []string) bool {
	found := 0
	for _, sl := range sublist {
		for _, l := range list {
			if sl == l {
				found++
			}
		}
	}
	return found == len(sublist)
}

func setenv(key, value string) {
	if err := os.Setenv(key, value); err != nil {
		log.Fatalln(err)
	}
}

func unsetenv(key string) {
	if err := os.Unsetenv(key); err != nil {
		log.Fatalln(err)
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

func commande(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
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
