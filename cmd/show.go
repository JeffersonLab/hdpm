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

fields: version, url, path, cmds, deps, isPrebuilt, dependents`,
	Example: `1. hdpm show
2. hdpm show url
3. hdpm show cmds
4. hdpm show dependents
5. hdpm show -p`,
	Run: runShow,
}

var showPrereqs, showExtras, showMaster bool

func init() {
	cmdHDPM.AddCommand(cmdShow)

	cmdShow.Flags().BoolVarP(&showExtras, "extras", "e", false, "Show extra package settings")
	cmdShow.Flags().BoolVarP(&showMaster, "master", "m", false, "Show master package settings")
	cmdShow.Flags().BoolVarP(&showPrereqs, "prereqs", "p", false, "Show GlueX offline software prerequisites")
}

func runShow(cmd *cobra.Command, args []string) {
	arg := "version"
	if len(args) == 1 {
		arg = args[0]
	}
	if showPrereqs {
		prereqs(arg)
		return
	}
	pkgInit()
	if showMaster {
		fmt.Println("Master (default) packages")
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("%-22s%-22s\n", "package", arg)
		fmt.Println(strings.Repeat("-", 80))
		for _, pkg := range masterPackages {
			pkg.config()
			pkg.show(arg)
		}
		return
	}
	if showExtras {
		fmt.Println("Extra packages")
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("%-22s%-22s\n", "package", arg)
		fmt.Println(strings.Repeat("-", 80))
		for _, pkg := range extraPackages {
			pkg.config()
			pkg.show(arg)
		}
		return
	}
	s := &Settings{}
	if isPath(SD + "/.info.json") {
		s.read(SD)
	}
	if s.Name == "" {
		s.Name = "master"
	}
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("Settings:  %s\n", s.Name)
	if s.Comment != "" {
		fmt.Printf("Comment:   %s\n", s.Comment)
	}
	if s.Timestamp != "" {
		fmt.Printf("Timestamp: %s\n", s.Timestamp)
	}
	fmt.Printf("GLUEX_TOP: %s\n", PD)
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-22s%-22s\n", "package", arg)
	fmt.Println(strings.Repeat("-", 80))
	for _, pkg := range packages {
		pkg.config()
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
		fmt.Printf("%-22s%-22s\n", p.Name, strings.Join(p.Deps, " "))
	case "dependents":
		fmt.Printf("%-22s%-22s\n", p.Name, strings.Join(dependents(p.Name), " "))
	case "isPrebuilt":
		fmt.Printf("%-22s%-22t\n", p.Name, p.IsPrebuilt)
	default:
		fmt.Printf("%-22s%-22s\n", p.Name, p.gitVersion())
	}
}

func prereqs(arg string) {
	OS = osrelease()
	if runtime.GOOS == "darwin" {
		OS = "macOS"
	}
	if arg != "version" {
		OS = arg
	}
	var msg string
	switch OS {
	case "CentOS6", "RHEL6":
		msg = `# CentOS/RHEL 6 prerequisites
yum update -y && yum install -y centos-release-SCL epel-release \
    centos-release-scl-rh \
    && yum install -y python27 git make gcc-c++ gcc binutils cmake3 scons \
    libX11-devel libXpm-devel libXft-devel libXext-devel mesa-libGLU-devel \
    libXi-devel libXmu-devel gcc-gfortran imake patch expat-devel boost-devel \
    blas-devel lapack-devel openmotif-devel mysql-devel sqlite-devel \
    fftw-devel bzip2 bzip2-devel tcsh devtoolset-3-toolchain \
    && ln -s liblapack.a /usr/lib64/liblapack3.a
`
	case "CentOS7", "RHEL7":
		msg = `# CentOS/RHEL 7 prerequisites
yum update -y && yum install -y epel-release && yum install -y \
    git make gcc-c++ gcc binutils python-devel cmake3 scons boost-devel \
    libX11-devel libXpm-devel libXft-devel libXext-devel mesa-libGLU-devel \
    gcc-gfortran imake patch expat-devel libXi-devel libXmu-devel \
    mysql-devel sqlite-devel fftw-devel bzip2 bzip2-devel tcsh \
    blas-devel blas-static lapack-devel lapack-static openmotif-devel \
    && ln -s liblapack.a /usr/lib64/liblapack3.a
`
	case "Ubuntu14", "LinuxMint17", "Ubuntu16", "LinuxMint18":
		msg = `# Ubuntu 14.04/16.04 LTS prerequisites
apt-get update && apt-get install -y curl git dpkg-dev make g++ gcc \
   binutils libx11-dev libxpm-dev libxft-dev libxext-dev libfftw3-dev tcsh \
   python-dev cmake scons gfortran xutils-dev libxt-dev libboost-python-dev \
   liblapack-dev libblas-dev libmotif-dev expect libgl1-mesa-dev libxmu-dev \
   libxi-dev libglew-dev libmysqlclient-dev sqlite3 libsqlite3-dev libbz2-dev \
   && ln -s make /usr/bin/gmake \
   && ln -s liblapack.a /usr/lib/liblapack3.a
`
	case "macOS", "OSX":
		msg = `# macOS prerequisites

1. xcode-select --install
2. Install XQuartz (https://www.xquartz.org)
3. Install Homebrew (http://brew.sh)
4. brew install scons cmake gcc mariadb boost-python
`
	default:
		fmt.Fprintf(os.Stderr, "%s: Unknown operating system\n", OS)
	}
	fmt.Print(msg)
}
