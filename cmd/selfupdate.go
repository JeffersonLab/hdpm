package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// Create the selfupdate command
var cmdSelfupdate = &cobra.Command{
	Use:   "selfupdate",
	Short: "Update hdpm",
	Long: `Update hdpm to the latest release.

hdpm is installed to the $GLUEX_TOP/.hdpm directory.
Any previous installation of hdpm at that location is removed.`,
	Example: `1. hdpm selfupdate
2. hdpm selfupdate -v 0.6.0`,
	Run: runSelfupdate,
}

var version string

func init() {
	cmdHDPM.AddCommand(cmdSelfupdate)

	cmdSelfupdate.Flags().StringVarP(&version, "version", "v", "", "Request a specific version")
	cmdSelfupdate.Flags().BoolVarP(&noCheckURL, "no-check-url", "", false, "Do not check URL")
}

func runSelfupdate(cmd *cobra.Command, args []string) {
	pathInit()

	// Set proxy env. variables if on JLab CUE
	setenvJLabProxy()

	// Update hdpm
	mkcd(PD)
	selfupdate()
}

func selfupdate() {
	ver := "latest"
	if version != "" {
		ver = version
	}
	url := "https://halldweb.jlab.org/dist/hdpm/hdpm-" + ver + ".linux.tar.gz"
	if ver == "latest" {
		ver = latestRelease("hdpm")
		url = strings.Replace(url, "latest", ver, 1)
	}
	if ver == VERSION {
		fmt.Printf("Already up-to-date: hdpm version %s\n", ver)
		return
	}
	if runtime.GOOS == "darwin" {
		url = strings.Replace(url, "linux", "macOS", 1)
	}
	checkURL(url)
	if isPath(HD + "/bin/hdpm") {
		run("cp", "-p", HD+"/bin/hdpm", HD+"/bin/hdpm.old")
		os.Remove(HD + "/bin/hdpm")
	}
	if err := fetchTarfileError(url, HD); err != nil {
		log.SetPrefix("fetch failed: ")
		log.SetFlags(0)
		log.Println(err)
		if isPath(HD + "/bin/hdpm.old") {
			os.Rename(HD+"/bin/hdpm.old", HD+"/bin/hdpm")
			fmt.Printf("Update failed: restored hdpm version %s\n", VERSION)
		}
		return
	}
	os.Remove(HD + "/bin/hdpm.old")
	fmt.Printf("Installed: %s\n", HD+"/bin/hdpm")
	fmt.Printf("Update succeeded: hdpm version %s\n", ver)
}

func latestRelease(name string) string {
	latest_release := "0.0.0"
	page := output("curl", "-s", "https://halldweb.jlab.org/dist/hdpm/")
	lines := strings.Split(page, "\n")
	for _, line := range lines {
		re := regexp.MustCompile("href=\".{20,30}\"")
		r := re.FindString(line)
		if r == "" {
			continue
		}
		file := r[6 : len(r)-1]
		prefix, suffix := name+"-", ".linux.tar.gz"
		if strings.HasPrefix(file, prefix) && strings.HasSuffix(file, suffix) && !strings.HasPrefix(file, name+"-dev.") {
			file = strings.TrimPrefix(file, prefix)
			file = strings.TrimSuffix(file, suffix)
			if strings.Contains(file, ".") {
				if isLater(file, latest_release) {
					latest_release = file
				}
			}
		}
	}
	if latest_release == "0.0.0" {
		fmt.Fprintf(os.Stderr, "No releases found at https://halldweb.jlab.org/dist/hdpm/ for %s.\n", name)
		os.Exit(2)
	}
	fmt.Printf("Latest release: %s version %s\n\n", name, latest_release)
	return latest_release
}

func isLater(v1, v2 string) bool {
	V1 := vnSlice(v1)
	V2 := vnSlice(v2)
	if len(V1) != len(V2) {
		return false
	}
	if len(V1) != 3 {
		return false
	}
	for i, _ := range V1 {
		if V1[i] == V2[i] {
			continue
		}
		return V1[i] > V2[i]
	}
	return false
}

func vnSlice(v string) (V []int) {
	for _, s := range strings.Split(v, ".") {
		i, _ := strconv.Atoi(s)
		V = append(V, i)
	}
	return V
}
