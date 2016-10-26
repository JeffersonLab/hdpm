package cmd

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// Create the cull command
var cmdCull = &cobra.Command{
	Use:   "cull [TAG...]",
	Short: "Cull distribution tarfiles",
	Long: `Cull distribution tarfiles.
	
tags: c6, c7, u14, u16

All tags will be culled if no arguments are given.

Usage examples:
1. gxd cull -n 10 c6
`,
	Run: runCull,
}

var N int
var show bool

func init() {
	cmdGXD.AddCommand(cmdCull)

	cmdCull.Flags().IntVarP(&N, "", "n", 10, "Number of tarfiles to keep per tag.")
	cmdCull.Flags().BoolVarP(&show, "show", "s", false, "Show tarfiles to be culled.")
}

func runCull(cmd *cobra.Command, args []string) {
	var tags = []string{"c6", "c7", "u14", "u16"}

	for _, arg := range args {
		if !in(tags, arg) {
			fmt.Fprintf(os.Stderr, "%s: Unknown tag\n", arg)
			os.Exit(2)
		}
	}

	if len(args) == 0 {
		args = tags
	}

	for _, tag := range tags {
		if !in(args, tag) {
			continue
		}
		cull(tag)
	}
}

func cull(tag string) {
	target := "/group/halld/www/halldweb/html/dist"

	rmGlob(target + "/sim-recon--*-" + tag + ".tar.gz")
	rmGlob(target + "/sim-recon-deps--*-" + tag + ".tar.gz")

	files := make(map[time.Time]string)
	latest_file := ""
	var latest_dt time.Time
	page := output("curl", "-s", "https://halldweb.jlab.org/dist/")
	lines := strings.Split(page, "\n")
	for _, line := range lines {
		re := regexp.MustCompile("href=\".{25,50}\"")
		r := re.FindString(line)
		if r == "" {
			continue
		}
		file := r[6 : len(r)-1]
		if strings.HasPrefix(file, "sim-recon-") && strings.HasSuffix(file, "-"+tag+".tar.gz") &&
			!strings.HasPrefix(file, "sim-recon-deps-") {
			re = regexp.MustCompile("(\\d{4})-(\\d{2})-(\\d{2})")
			r = re.FindString(line)
			s1 := strings.Split(r, "-")
			re = regexp.MustCompile("(\\d{2}):(\\d{2})")
			r = re.FindString(line)
			s2 := strings.Split(r, ":")
			dt := mkDate(s1, s2)
			files[dt] = file
			if latest_dt.Before(dt) {
				latest_dt = dt
				latest_file = file
			}
		}
	}
	if latest_file == "" {
		fmt.Fprintf(os.Stderr, "File not found at https://halldweb.jlab.org/dist for %s OS tag.\n", tag)
		return
	}
	var times []time.Time
	for t, _ := range files {
		times = append(times, t)
	}
	sort.Sort(ts(times))
	i := 0
	for _, t := range times {
		i++
		if i > N {
			if show {
				fmt.Printf("%s        %v\n", files[t], t)
			} else {
				os.Remove(target + "/" + files[t])
			}
		}
	}
}

// Sort by time (satisfy sort interface)
type ts []time.Time

func (a ts) Len() int           { return len(a) }
func (a ts) Less(j, i int) bool { return a[i].Before(a[j]) }
func (a ts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func mkDate(s1 []string, s2 []string) time.Time {
	s1 = append(s1, s2...)
	var n []int
	for _, s := range s1 {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalln(err)
		}
		n = append(n, i)
	}
	return time.Date(n[0], time.Month(n[1]), n[2], n[3], n[4], 0, 0, time.Local)
}
