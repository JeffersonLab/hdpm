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
	Long: `
Print GlueX environment variables in "key=value" format.

Pass an environment variable name as an argument to print it.

bash and tcsh environment-setup scripts are written to
$GLUEX_TOP/env-setup/<template-id>.[c]sh.
`,
	Run: runEnv,
}

func init() {
	cmdHDPM.AddCommand(cmdEnv)
}

func runEnv(cmd *cobra.Command, args []string) {
	arg := "ALL"
	if len(args) >= 1 {
		arg = args[0]
	}
	env(arg)
}

func env(arg string) {
	ENV := getEnv()
	if arg != "" {
		printEnv(arg, ENV)
	}
	printEnv("sh", ENV)
	printEnv("csh", ENV)
	for k, v := range ENV {
		if v != "" {
			setenv(k, v)
		}
	}
}

func printEnv(arg string, ENV map[string]string) {
	id := "master"
	dir := filepath.Join(ENV["GLUEX_TOP"], "settings")
	if isPath(dir) {
		id = readFile(dir + "/.id")
	}
	mk(filepath.Join(ENV["GLUEX_TOP"], "env-setup"))
	n := "kv"
	set := ""
	eq := "="
	if arg == "sh" {
		n = "bash"
		set = "export"
	}
	if arg == "csh" {
		n = "tcsh"
		set = "setenv"
		eq = " "
	}
	var keys []string
	for k, _ := range ENV {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if n == "kv" {
		fmt.Println("GLUEX_TOP=" + ENV["GLUEX_TOP"])
		fmt.Println("BMS_OSNAME=" + ENV["BMS_OSNAME"])
		for _, k := range keys {
			if k == "GLUEX_TOP" || k == "BMS_OSNAME" {
				continue
			}
			v := ENV[k]
			v = strings.Replace(v, ENV["GLUEX_TOP"], "${GLUEX_TOP}", -1)
			v = strings.Replace(v, ENV["BMS_OSNAME"], "${BMS_OSNAME}", -1)
			if arg == "ALL" || arg == k {
				fmt.Printf("%s=%s\n", k, v)
			}
		}
		return
	}
	f, err := os.Create(filepath.Join(ENV["GLUEX_TOP"], "env-setup", id+"."+arg))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprintln(f, "# "+n)
	fmt.Fprintln(f, set+" GLUEX_TOP"+eq+ENV["GLUEX_TOP"])
	fmt.Fprintln(f, set+" BMS_OSNAME"+eq+ENV["BMS_OSNAME"])
	for _, k := range keys {
		v := ENV[k]
		if k == "GLUEX_TOP" || k == "BMS_OSNAME" || strings.Contains(k, "PATH") {
			continue
		}
		v = strings.Replace(v, ENV["GLUEX_TOP"], "${GLUEX_TOP}", -1)
		v = strings.Replace(v, ENV["BMS_OSNAME"], "${BMS_OSNAME}", -1)
		fmt.Fprintln(f, set+" "+k+eq+v)
	}
	fmt.Fprintln(f, "#"+set+" JANA_CALIB_CONTEXT"+eq+"\"variation=mc\"")
	if ENV["JANA_RESOURCE_DIR"] == "" {
		fmt.Fprintln(f, "#"+set+" JANA_RESOURCE_DIR"+eq+"/path/to/resources")
	}
	ldlp := "LD_LIBRARY_PATH"
	if runtime.GOOS == "darwin" {
		ldlp = "DYLD_LIBRARY_PATH"
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
			path = strings.Replace(path, v, "${"+k+"}", -1)
		}
		path = strings.Replace(path, ENV["BMS_OSNAME"], "${BMS_OSNAME}", -1)
		if os.Getenv(pn) != "" {
			path = strings.Replace(path, os.Getenv(pn), "${"+pn+"}", -1)
		}
		path = strings.Replace(path, ENV["GLUEX_TOP"], "${GLUEX_TOP}", -1)
		if os.Getenv(pn) != "" {
			fmt.Fprintln(f, "\n"+set+" "+pn+eq+strings.Replace(os.Getenv(pn), ENV["GLUEX_TOP"], "${GLUEX_TOP}", -1))
		}
		fmt.Fprintln(f, "\n"+set+" "+pn+eq+path)
	}
	f.Close()
}

func getEnv() map[string]string {
	GLUEX_TOP := packageDir()
	BMS_OSNAME := osrelease()
	path, ver := make(map[string]string), make(map[string]string)
	for _, p := range packages {
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
		"AMPTOOLS":    filepath.Join(path["amptools"], "AmpTools"),
		"AMPPLOTTER":  filepath.Join(path["amptools"], "AmpPlotter"),
		"XERCESCROOT": path["xerces-c"],
		"EVIOROOT": filepath.Join(path["evio"], output("uname", "-s")+"-"+
			output("uname", "-m")),
		"RCDB_HOME":          path["rcdb"],
		"RCDB_CONNECTION":    "mysql://rcdb@hallddb.jlab.org/rcdb",
		"CCDB_HOME":          path["ccdb"],
		"CCDB_CONNECTION":    CCDB_CONNECTION,
		"CCDB_USER":          "${USER}",
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
	enames := []string{"HALLD_MY", "PATH", "LD_LIBRARY_PATH", "PYTHONPATH", "JANA_PLUGIN_PATH"}
	if runtime.GOOS == "darwin" {
		ENV["CERN"] = ""
		ENV["CERN_LEVEL"] = ""
		enames[2] = "DYLD_LIBRARY_PATH"
	}
	for _, n := range enames {
		ENV[n] = os.Getenv(n)
	}
	paths := []string{filepath.Join(ENV["CERN"], ENV["CERN_LEVEL"]), ENV["ROOTSYS"], ENV["XERCESCROOT"], ENV["EVIOROOT"], filepath.Join(ENV["RCDB_HOME"], "cpp"), ENV["CCDB_HOME"], ENV["JANA_HOME"], filepath.Join(ENV["HALLD_HOME"], ENV["BMS_OSNAME"]), filepath.Join(ENV["ROOT_ANALYSIS_HOME"], ENV["BMS_OSNAME"])}
	// PATH and LD_LIBRARY_PATH
	for _, p := range paths {
		if p != "" && p != "cpp" {
			ENV["PATH"] = addPath(ENV["PATH"], filepath.Join(p, "bin"))
			ENV[enames[2]] = addPath(ENV[enames[2]], filepath.Join(p, "lib"))
		}
	}
	if ENV["RCDB_HOME"] != "" {
		ENV["PATH"] = addPath(ENV["PATH"], filepath.Join(ENV["RCDB_HOME"], "bin"))
		ENV["PATH"] = ENV["RCDB_HOME"] + ":" + ENV["PATH"]
	}
	if path["cmake"] != "" {
		ENV["PATH"] = addPath(ENV["PATH"], filepath.Join(path["cmake"], "bin"))
	}
	// PYTHONPATH
	pypaths := []string{filepath.Join(ENV["ROOTSYS"], "lib"), filepath.Join(ENV["RCDB_HOME"], "python"), filepath.Join(ENV["CCDB_HOME"], "python") + ":" + filepath.Join(ENV["CCDB_HOME"], "python", "ccdb", "ccdb_pyllapi/"), filepath.Join(ENV["HALLD_HOME"], ENV["BMS_OSNAME"], "lib/python")}
	for _, pyp := range pypaths {
		ENV["PYTHONPATH"] = addPath(ENV["PYTHONPATH"], pyp)
	}
	// JANA_PLUGIN_PATH
	plugin_paths := []string{filepath.Join(ENV["JANA_HOME"], "plugins"), filepath.Join(ENV["HALLD_HOME"], ENV["BMS_OSNAME"], "plugins")}
	if ENV["HALLD_MY"] != "" {
		plugin_paths = append(plugin_paths, filepath.Join(ENV["HALLD_MY"], ENV["BMS_OSNAME"], "plugins"))
	}
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

func addPath(path, new_path string) string {
	if path == "" {
		return new_path
	}
	if !strings.Contains(path, new_path) && !strings.HasPrefix(new_path, "/usr/local/") {
		return new_path + ":" + path
	}
	return path
}

func envInit() {
	osr := osrelease()
	if (strings.Contains(osr, "CentOS6") || strings.Contains(osr, "RHEL6")) && (isPath(filepath.Join(packageDir(), ".dist")) || in(os.Args, "install")) {
		p := filepath.Join(packageDir(), ".dist")
		setenv("PATH", p+"/opt/rh/python27/root/usr/bin:"+p+"/opt/rh/devtoolset-3/root/usr/bin:"+os.Getenv("PATH"))
		a := p + "/opt/rh/python27/root/usr/lib64:" + p + "/opt/rh/devtoolset-3/root/usr/lib64:" + p + "/opt/rh/devtoolset-3/root/usr/lib"
		if os.Getenv("LD_LIBRARY_PATH") == "" {
			setenv("LD_LIBRARY_PATH", a)
		} else {
			setenv("LD_LIBRARY_PATH", a+":"+os.Getenv("LD_LIBRARY_PATH"))
		}
		setenv("LDFLAGS", "-L"+p+"/opt/rh/python27/root/usr/lib64")
	} else if (strings.Contains(osr, "CentOS6") || strings.Contains(osr, "RHEL6")) && isPath("/apps/gcc/4.9.2/bin") && isPath("/apps/python/PRO/bin") {
		p := "/apps"
		setenv("PATH", p+"/python/PRO/bin:"+p+"/gcc/4.9.2/bin:"+os.Getenv("PATH"))
		a := p + "/python/PRO/lib:" + p + "/gcc/4.9.2/lib64:" + p + "/gcc/4.9.2/lib"
		if os.Getenv("LD_LIBRARY_PATH") == "" {
			setenv("LD_LIBRARY_PATH", a)
		} else {
			setenv("LD_LIBRARY_PATH", a+":"+os.Getenv("LD_LIBRARY_PATH"))
		}
	}
}