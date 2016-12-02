package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// Create the show command
var cmdShow = &cobra.Command{
	Use:   "show [FIELD]",
	Short: "Show the current package settings",
	Long: `Show the current package settings.

The package names and versions are printed by default.

fields: version, url, path, cmds, deps, isPrebuilt`,
	Example: `1. hdpm show
2. hdpm show url
3. hdpm show cmds
4. hdpm show -p`,
	Run: runShow,
}

var showPrereqs bool

func init() {
	cmdHDPM.AddCommand(cmdShow)

	cmdShow.Flags().BoolVarP(&showPrereqs, "prereqs", "p", false, "Show GlueX offline software prerequisites")
}

func runShow(cmd *cobra.Command, args []string) {
	pkgInit()
	arg := "version"
	if len(args) == 1 {
		arg = args[0]
	}
	if showPrereqs {
		prereqs(arg)
		return
	}
	id := "master"
	if isPath(SD + "/.id") {
		id = readFile(SD + "/.id")
	}
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("Settings id: %s\n", id)
	fmt.Printf("GLUEX_TOP:   %s\n", PD)
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-22s%-22s\n", "package", arg)
	fmt.Println(strings.Repeat("-", 80))
	for _, pkg := range packages {
		pkg.config()
		pkg.template()
		pkg.show(arg)
	}
}

func (p *Package) show(arg string) {
	switch arg {
	case "url":
		fmt.Printf("%-22s%-22s\n", p.Name, p.URL)
	case "path":
		fmt.Printf("%-22s%-22s\n", p.Name, p.Path)
	case "cmds":
		for _, cmd := range p.Cmds {
			fmt.Printf("%-22s%-22s\n", p.Name, cmd)
		}
	case "deps":
		fmt.Printf("%-22s%-22s\n", p.Name, strings.Join(p.Deps, ", "))
	case "isPrebuilt":
		fmt.Printf("%-22s%-22t\n", p.Name, p.IsPrebuilt)
	default:
		fmt.Printf("%-22s%-22s\n", p.Name, p.Version)
	}
}

func prereqs(arg string) {
	tag := ""
	if strings.Contains(OS, "CentOS6") || strings.Contains(OS, "RHEL6") {
		tag = "c6"
	}
	if strings.Contains(OS, "CentOS7") || strings.Contains(OS, "RHEL7") {
		tag = "c7"
	}
	if strings.Contains(OS, "Ubuntu14") || strings.Contains(OS, "LinuxMint17") {
		tag = "u14"
	}
	if strings.Contains(OS, "Ubuntu16") || strings.Contains(OS, "LinuxMint18") {
		tag = "u16"
	}
	if runtime.GOOS == "darwin" {
		tag = "macOS"
	}
	if arg != "version" {
		tag = arg
	}
	if tag == "" {
		fmt.Fprintf(os.Stderr, "%s: Unsupported operating system\n", OS)
		os.Exit(2)
	}
	var msg string
	switch {
	case tag == "c6":
		msg = `# CentOS/RHEL 6 prerequisites

yum update -y && yum install -y centos-release-SCL epel-release \
	centos-release-scl-rh \
	&& yum install -y python27 git make gcc-c++ gcc binutils \
	libX11-devel libXpm-devel libXft-devel libXext-devel \
	subversion scons gcc-gfortran imake patch expat-devel \
	blas-devel lapack-devel openmotif-devel mysql-devel sqlite-devel \
	fftw-devel bzip2 bzip2-devel tcsh devtoolset-3-toolchain \
	&& ln -s /usr/lib64/liblapack.a /usr/lib64/liblapack3.a
`
	case tag == "c7":
		msg = `# CentOS/RHEL 7 prerequisites

yum update -y && yum install -y epel-release && yum install -y \
	git make gcc-c++ gcc binutils python-devel \
	libX11-devel libXpm-devel libXft-devel libXext-devel \
	subversion scons gcc-gfortran imake patch expat-devel \
	mysql-devel sqlite-devel fftw-devel bzip2 bzip2-devel tcsh  \
	blas-devel blas-static lapack-devel lapack-static openmotif-devel \
	&& ln -s /usr/lib64/liblapack.a /usr/lib64/liblapack3.a
`
	case tag == "u14" || tag == "u16":
		msg = `# Ubuntu 14.04/16.04 LTS prerequisites

apt-get update && apt-get install -y curl git dpkg-dev make g++ gcc \
	binutils libx11-dev libxpm-dev libxft-dev libxext-dev libfftw3-dev \
	python-dev scons subversion gfortran xutils-dev libxt-dev \
	liblapack-dev libblas-dev libmotif-dev expect libgl1-mesa-dev \
	libglew-dev libmysqlclient-dev sqlite3 libsqlite3-dev tcsh libbz2-dev \
	&& ln -s /usr/bin/make /usr/bin/gmake \
	&& ln -s /usr/lib/liblapack.a /usr/lib/liblapack3.a
`
	case tag == "macOS":
		msg = `# macOS prerequisites

1. xcode-select --install
2. Install XQuartz (https://www.xquartz.org)
3. Install Homebrew (http://brew.sh)
4. brew install scons cmake gcc mariadb
`
	}
	fmt.Print(msg)
}
