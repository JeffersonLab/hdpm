package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// cmdEnv represents the env command
var cmdEnv = &cobra.Command{
	Use:   "env [VAR]",
	Short: "Print GlueX environment variables",
	Long: `Print GlueX environment variables.

Pass an environment variable as an argument to select it.

-w,--write flag:
Write bash and tcsh environment-setup scripts to
$GLUEX_TOP/.hdpm/env/NAME.[c]sh.

This command can also be used to set/unset GlueX env variables in your shell.
See examples below.`,
	Example: `1. hdpm env
2. hdpm env -s bash` +
		"\n3. (a) eval `hdpm env`        (Set env for tcsh shell)\n   (b) eval `hdpm env -u`     (Unset env for tcsh shell)\n" +
		`4. (a) eval "$(hdpm env)"     (Set env for bash shell)
   (b) eval "$(hdpm env -u)"  (Unset env for bash shell)`,
	Run: runEnv,
}

var shell string
var unset, write bool

func init() {
	cmdHDPM.AddCommand(cmdEnv)

	cmdEnv.Flags().BoolVarP(&write, "write", "w", false, "Write env scripts to $GLUEX_TOP/.hdpm/env/")
	cmdEnv.Flags().BoolVarP(&unset, "unset", "u", false, "Print unset commands for GlueX env variables")
	cmdEnv.Flags().StringVarP(&shell, "shell", "s", "", "Print commands for bash or tcsh shell")
}

func runEnv(cmd *cobra.Command, args []string) {
	pkgInit()
	arg := ""
	if len(args) >= 1 {
		arg = args[0]
	}
	ENV := getEnv()
	if write {
		writeEnv("sh", ENV)
		writeEnv("csh", ENV)
		return
	}
	env(arg, ENV)
}

func setEnv() {
	ENV := getEnv()
	writeEnv("sh", ENV)
	writeEnv("csh", ENV)
	delete(ENV, "PATH0")
	for k, v := range ENV {
		if v != "" {
			setenv(k, v)
		}
	}
}

func env(arg string, ENV map[string]string) {
	if shell == "" {
		shell = filepath.Base(os.Getenv("SHELL"))
	}
	type shSyntax struct {
		name, set, unset, eq, end, uend string
	}
	sh := shSyntax{"tcsh", "setenv", "unsetenv", " \"", "\";\n", ";\n"}
	if shell == "bash" || shell == "sh" {
		sh = shSyntax{"bash", "export", "unset", "=\"", "\"\n", "\n"}
	}
	var keys []string
	for k, _ := range ENV {
		if k != "PATH0" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	if unset {
		if arg == "" {
			fmt.Printf("%s %s%s%s%s", sh.set, "PATH", sh.eq, ENV["PATH"], sh.end)
		}
		for _, k := range keys {
			if k == "GLUEX_TOP" || k == "PATH" || k == "HALLD_MY" {
				continue
			}
			if arg == "" || arg == k {
				fmt.Printf("%s %s%s", sh.unset, k, sh.uend)
			}
		}
		ldlp := "LD_LIBRARY_PATH"
		if runtime.GOOS == "darwin" {
			ldlp = "DYLD_LIBRARY_PATH"
		}
		for _, k := range []string{ldlp, "PYTHONPATH", "JANA_PLUGIN_PATH"} {
			if ENV[k] == "" {
				if arg == "" || arg == k {
					fmt.Printf("%s %s%s", sh.unset, k, sh.uend)
				}
			}
		}
		return
	}
	for _, k := range []string{"GLUEX_TOP", "BMS_OSNAME"} {
		fmt.Printf("%s %s%s%s%s", sh.set, k, sh.eq, ENV[k], sh.end)
	}
	for _, k := range keys {
		if k == "GLUEX_TOP" || k == "BMS_OSNAME" {
			continue
		}
		if arg == "" || arg == k {
			fmt.Printf("%s %s%s%s%s", sh.set, k, sh.eq, ENV[k], sh.end)
		}
	}
}

func writeEnv(arg string, ENV map[string]string) {
	s := &Settings{}
	if isPath(SD + "/.info.json") {
		s.read(SD)
	}
	if s.Name == "" {
		s.Name = "master"
	}
	id := s.Name
	mk(filepath.Join(HD, "env"))
	type shSyntax struct {
		name, set, unset, eq, end, uend string
	}
	if arg != "sh" {
		arg = "csh"
	}
	sh := shSyntax{"tcsh", "setenv", "unsetenv", " \"", "\";\n", ";\n"}
	if arg == "sh" {
		sh = shSyntax{"bash", "export", "unset", "=\"", "\"\n", "\n"}
	}
	var keys []string
	for k, _ := range ENV {
		if k != "PATH0" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	f, err := os.Create(filepath.Join(HD, "env", id+"."+arg))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintln(f, "# "+sh.name)
	fmt.Fprintf(f, sh.set+" GLUEX_TOP"+sh.eq+ENV["GLUEX_TOP"]+sh.end)
	fmt.Fprintf(f, sh.set+" BMS_OSNAME"+sh.eq+ENV["BMS_OSNAME"]+sh.end)
	for _, k := range keys {
		v := ENV[k]
		if k == "GLUEX_TOP" || k == "BMS_OSNAME" || strings.Contains(k, "PATH") {
			continue
		}
		if k == "CCDB_USER" {
			v = "${USER}"
		}
		v = strings.Replace(v, ENV["GLUEX_TOP"], "${GLUEX_TOP}", 1)
		v = strings.Replace(v, ENV["BMS_OSNAME"], "${BMS_OSNAME}", 2)
		fmt.Fprintf(f, sh.set+" "+k+sh.eq+v+sh.end)
	}
	fmt.Fprintf(f, "#"+sh.set+" JANA_CALIB_CONTEXT"+sh.eq+"variation=mc"+sh.end)
	if ENV["JANA_RESOURCE_DIR"] == "" {
		fmt.Fprintf(f, "#"+sh.set+" JANA_RESOURCE_DIR"+sh.eq+"/path/to/resources"+sh.end)
	}
	ldlp := "LD_LIBRARY_PATH"
	if runtime.GOOS == "darwin" {
		ldlp = "DYLD_LIBRARY_PATH"
	}
	isCentOS6JLabFarm := false
	if strings.Contains(OS, "CentOS6") || strings.Contains(OS, "RHEL6") {
		if strings.Contains(OS, "gcc") && isPath("/apps/gcc/4.9.2/bin") {
			isCentOS6JLabFarm = true
		}
	}
	for _, pn := range []string{"PATH", ldlp, "PYTHONPATH", "JANA_PLUGIN_PATH"} {
		if ENV[pn] == "" {
			continue
		}
		path := ENV[pn]
		for k, v := range ENV {
			if k == "GLUEX_TOP" || k == "BMS_OSNAME" || k == "CCDB_USER" || strings.Contains(k, "PATH") {
				continue
			}
			if strings.HasPrefix(k, "G4") && k != "G4ROOT" && k != "G4WORKDIR" {
				continue
			}
			path = strings.Replace(path, v, "${"+k+"}", -1)
		}
		if pn == "PATH" && !isCentOS6JLabFarm {
			path = strings.Replace(path, ENV["PATH0"], "${PATH}", 1)
		}
		path = strings.Replace(path, ENV["BMS_OSNAME"], "${BMS_OSNAME}", -1)
		path = strings.Replace(path, ENV["GLUEX_TOP"], "${GLUEX_TOP}", -1)
		fmt.Fprintf(f, "\n"+sh.set+" "+pn+sh.eq+path+sh.end)
	}
	f.Close()
}

func getEnv() map[string]string {
	GLUEX_TOP := PD
	BMS_OSNAME := OS
	path, ver := make(map[string]string), make(map[string]string)
	for _, p := range packages {
		p.config()
		path[p.Name] = p.Path
		ver[p.Name] = p.Version
	}
	CCDB_CONNECTION := "mysql://ccdb_user@hallddb.jlab.org/ccdb"
	var ENV = map[string]string{
		"GLUEX_TOP":   GLUEX_TOP,
		"BMS_OSNAME":  BMS_OSNAME,
		"CERN":        path["cernlib"],
		"CERN_LEVEL":  ver["cernlib"],
		"ROOTSYS":     path["root"],
		"XERCESCROOT": path["xerces-c"],
		"EVIOROOT": filepath.Join(path["evio"], output("uname", "-s")+"-"+
			output("uname", "-m")),
		"RCDB_HOME":          path["rcdb"],
		"RCDB_CONNECTION":    "mysql://rcdb@hallddb.jlab.org/rcdb",
		"CCDB_HOME":          path["ccdb"],
		"CCDB_CONNECTION":    CCDB_CONNECTION,
		"HDDS_HOME":          path["hdds"],
		"JANA_HOME":          filepath.Join(path["jana"], BMS_OSNAME),
		"JANA_CALIB_URL":     CCDB_CONNECTION,
		"JANA_GEOMETRY_URL":  "xmlfile://" + path["hdds"] + "/main_HDDS.xml",
		"HALLD_HOME":         path["sim-recon"],
		"JANA_RESOURCE_DIR":  "/u/group/halld/www/halldweb/html/resources",
		"ROOT_ANALYSIS_HOME": path["gluex_root_analysis"],
	}
	if !isPath(ENV["JANA_RESOURCE_DIR"]) {
		ENV["JANA_RESOURCE_DIR"] = ""
	}
	ENV["CCDB_USER"] = os.Getenv("CCDB_USER")
	if ENV["CCDB_USER"] == "" {
		ENV["CCDB_USER"] = os.Getenv("USER")
	}
	if path["amptools"] != "" {
		ENV["AMPTOOLS"] = filepath.Join(path["amptools"], "AmpTools")
		ENV["AMPPLOTTER"] = filepath.Join(path["amptools"], "AmpPlotter")
	}
	if path["hd_utilities"] != "" {
		mcw := os.Getenv("MCWRAPPER_CENTRAL")
		if mcw == "" {
			ENV["MCWRAPPER_CENTRAL"] = filepath.Join(path["hd_utilities"], "MCwrapper")
		} else {
			ENV["MCWRAPPER_CENTRAL"] = mcw
		}
	}
	if path["hdgeant4"] != "" {
		ENV["G4WORKDIR"] = path["hdgeant4"]
		setenv("G4WORKDIR", ENV["G4WORKDIR"])
	}
	var isG4Installed bool
	if path["geant4"] != "" {
		ENV["G4ROOT"] = path["geant4"]
		ENV["G4VERSION"] = ver["geant4"]
		g4c := path["geant4"] + "/bin/geant4-config"
		if isPath(g4c) {
			v := output("sh", "-c", g4c+" --version")
			ENV["G4INSTALL"] = path["geant4"] + "/share/Geant4-" + v + "/geant4make"
			if isPath(ENV["G4INSTALL"]) {
				ENV = addG4(ENV)
				isG4Installed = true
			}
		}
	}
	if isJLabFarm() {
		ENV["http_proxy"] = "http://jprox.jlab.org:8081"
		ENV["https_proxy"] = "https://jprox.jlab.org:8081"
	}
	enames := []string{"HALLD_MY", "PATH", "LD_LIBRARY_PATH", "PYTHONPATH", "JANA_PLUGIN_PATH"}
	g4LibDir := "lib64"
	if runtime.GOOS == "darwin" {
		ENV["CERN"] = ""
		ENV["CERN_LEVEL"] = ""
		enames[2] = "DYLD_LIBRARY_PATH"
		g4LibDir = "lib"
	}
	for _, n := range enames {
		ENV[n] = os.Getenv(n)
	}
	// PATH and LD_LIBRARY_PATH
	// First remove old entries
	cpaths := []string{filepath.Join(os.Getenv("CERN"), os.Getenv("CERN_LEVEL")), os.Getenv("ROOTSYS"), os.Getenv("XERCESCROOT"), os.Getenv("EVIOROOT"), filepath.Join(os.Getenv("RCDB_HOME"), "cpp"), os.Getenv("CCDB_HOME"), os.Getenv("JANA_HOME"), filepath.Join(os.Getenv("HALLD_HOME"), os.Getenv("BMS_OSNAME")), filepath.Join(os.Getenv("ROOT_ANALYSIS_HOME"), os.Getenv("BMS_OSNAME"))}
	for _, p := range cpaths {
		ENV["PATH"] = cleanPath(ENV["PATH"], filepath.Join(p, "bin"))
		ENV[enames[2]] = cleanPath(ENV[enames[2]], filepath.Join(p, "lib"))
	}
	if ENV["RCDB_HOME"] != "" {
		ENV["PATH"] = cleanPath(ENV["PATH"], filepath.Join(os.Getenv("RCDB_HOME"), "bin"))
		ENV["PATH"] = cleanPath(ENV["PATH"], os.Getenv("RCDB_HOME"))
	}
	if ENV["G4ROOT"] != "" {
		g4root := os.Getenv("G4ROOT")
		ENV["PATH"] = cleanPath(ENV["PATH"], filepath.Join(g4root, "bin"))
		ENV[enames[2]] = cleanPath(ENV[enames[2]], filepath.Join(g4root, g4LibDir))
	}
	if ENV["G4WORKDIR"] != "" {
		ENV["PATH"] = cleanPath(ENV["PATH"], filepath.Join(os.Getenv("G4WORKDIR"), "bin", os.Getenv("G4SYSTEM")))
		ENV[enames[2]] = cleanPath(ENV[enames[2]], filepath.Join(os.Getenv("G4WORKDIR"), "tmp", os.Getenv("G4SYSTEM"), "hdgeant4"))
	}
	if ENV["MCWRAPPER_CENTRAL"] != "" {
		ENV["PATH"] = cleanPath(ENV["PATH"], os.Getenv("MCWRAPPER_CENTRAL"))
	}
	ENV["PATH0"] = ENV["PATH"]
	if isJLabFarm() && isPath("/apps/cmake/cmake-3.5.1") {
		ENV["PATH"] = addPath(ENV["PATH"], "/apps/cmake/cmake-3.5.1/bin")
	}
	paths := []string{filepath.Join(ENV["CERN"], ENV["CERN_LEVEL"]), ENV["ROOTSYS"], ENV["XERCESCROOT"], ENV["EVIOROOT"], filepath.Join(ENV["RCDB_HOME"], "cpp"), ENV["CCDB_HOME"], ENV["JANA_HOME"], filepath.Join(ENV["HALLD_HOME"], ENV["BMS_OSNAME"]), filepath.Join(ENV["ROOT_ANALYSIS_HOME"], ENV["BMS_OSNAME"])}
	for _, p := range paths {
		ENV["PATH"] = addPath(ENV["PATH"], filepath.Join(p, "bin"))
		ENV[enames[2]] = addPath(ENV[enames[2]], filepath.Join(p, "lib"))
	}
	if ENV["RCDB_HOME"] != "" {
		ENV["PATH"] = addPath(ENV["PATH"], filepath.Join(ENV["RCDB_HOME"], "bin"))
		if !unset {
			ENV["PATH"] = ENV["RCDB_HOME"] + ":" + ENV["PATH"]
		}
	}
	if ENV["G4ROOT"] != "" {
		if isG4Installed {
			ENV["PATH"] = addPath(ENV["PATH"], filepath.Join(ENV["G4ROOT"], "bin"))
			ENV[enames[2]] = addPath(ENV[enames[2]], filepath.Join(ENV["G4ROOT"], g4LibDir))
		}
	}
	if ENV["G4WORKDIR"] != "" {
		if isG4Installed {
			ENV["PATH"] = addPath(ENV["PATH"], filepath.Join(ENV["G4WORKDIR"], "bin", ENV["G4SYSTEM"]))
			ENV[enames[2]] = addPath(ENV[enames[2]], filepath.Join(ENV["G4WORKDIR"], "tmp", ENV["G4SYSTEM"], "hdgeant4"))
		}
	}
	if ENV["MCWRAPPER_CENTRAL"] != "" {
		ENV["PATH"] = addPath(ENV["PATH"], ENV["MCWRAPPER_CENTRAL"])
	}
	// PYTHONPATH
	cpypaths := []string{filepath.Join(os.Getenv("ROOTSYS"), "lib"), filepath.Join(os.Getenv("RCDB_HOME"), "python"), filepath.Join(os.Getenv("CCDB_HOME"), "python") + ":" + filepath.Join(os.Getenv("CCDB_HOME"), "python", "ccdb", "ccdb_pyllapi/"), filepath.Join(os.Getenv("HALLD_HOME"), os.Getenv("BMS_OSNAME"), "lib/python")}
	for _, pyp := range cpypaths {
		ENV["PYTHONPATH"] = cleanPath(ENV["PYTHONPATH"], pyp)
	}
	pypaths := []string{filepath.Join(ENV["ROOTSYS"], "lib"), filepath.Join(ENV["RCDB_HOME"], "python"), filepath.Join(ENV["CCDB_HOME"], "python") + ":" + filepath.Join(ENV["CCDB_HOME"], "python", "ccdb", "ccdb_pyllapi/"), filepath.Join(ENV["HALLD_HOME"], ENV["BMS_OSNAME"], "lib/python")}
	for _, pyp := range pypaths {
		ENV["PYTHONPATH"] = addPath(ENV["PYTHONPATH"], pyp)
	}
	// JANA_PLUGIN_PATH
	cplugin_paths := []string{filepath.Join(os.Getenv("JANA_HOME"), "plugins"), filepath.Join(os.Getenv("HALLD_HOME"), os.Getenv("BMS_OSNAME"), "plugins"), filepath.Join(os.Getenv("HALLD_MY"), os.Getenv("BMS_OSNAME"), "plugins")}
	for _, plugin_path := range cplugin_paths {
		ENV["JANA_PLUGIN_PATH"] = cleanPath(ENV["JANA_PLUGIN_PATH"], plugin_path)
	}
	plugin_paths := []string{filepath.Join(ENV["JANA_HOME"], "plugins"), filepath.Join(ENV["HALLD_HOME"], ENV["BMS_OSNAME"], "plugins"), filepath.Join(ENV["HALLD_MY"], ENV["BMS_OSNAME"], "plugins")}
	for _, plugin_path := range plugin_paths {
		ENV["JANA_PLUGIN_PATH"] = addPath(ENV["JANA_PLUGIN_PATH"], plugin_path)
	}
	for k, v := range ENV {
		if v == "" {
			delete(ENV, k)
		}
	}
	return ENV
}

func addG4(e map[string]string) map[string]string {
	d := output("bash", "-c", ". "+e["G4INSTALL"]+"/geant4make.sh; env")
	for _, line := range strings.Split(d, "\n") {
		if strings.HasPrefix(line, "G4") {
			s := strings.Split(line, "=")
			if len(s) == 2 {
				e[s[0]] = s[1]
			}
		}
	}
	return e
}

func addPath(path, new_path string) string {
	if !filepath.IsAbs(new_path) || unset {
		return path
	}
	if path == "" {
		return new_path
	}
	if !strings.Contains(path, new_path) && !strings.HasPrefix(new_path, "/usr/local/") {
		return new_path + ":" + path
	}
	return path
}

func cleanPath(path, old_path string) string {
	if path == "" || !filepath.IsAbs(old_path) {
		return path
	}
	if strings.Contains(path, old_path) && !strings.HasPrefix(old_path, "/usr/local/") {
		if strings.HasSuffix(path, old_path) {
			return strings.Replace(path, old_path, "", -1)
		} else {
			return strings.Replace(path, old_path+":", "", -1)
		}
	}
	return path
}

func setenvJLabProxy() {
	if isJLabFarm() {
		setenv("http_proxy", "http://jprox.jlab.org:8081")
		setenv("https_proxy", "https://jprox.jlab.org:8081")
	}
}

func setenvp(name, path string) {
	p0 := os.Getenv(name)
	if p0 == "" {
		setenv(name, path)
		return
	}
	if !strings.Contains(p0, path) {
		setenv(name, path+":"+p0)
	}
}

func setenvGCC() {
	if isPath("/apps/gcc/4.9.2/bin") {
		v := output("gcc", "-dumpversion")
		setenvp("PATH", "/apps/gcc/4.9.2/bin")
		setenvp("LD_LIBRARY_PATH", "/apps/gcc/4.9.2/lib64:/apps/gcc/4.9.2/lib")
		OS = strings.Replace(OS, v, "4.9.2", 1)
		pyp := JPD + "/" + OS + "/python/Python-2.7.13"
		if isPath(pyp) {
			setenvp("PATH", pyp+"/bin")
			setenvp("LD_LIBRARY_PATH", pyp+"/lib")
		}
	}
}
