# hdpm - Hall D Package Manager
This is a set of tools for managing **top-level builds** of the Jefferson Lab Hall-D offline software and online monitoring plugins; Specify a list of paths to the software and its dependencies, and it will set the Hall-D environment variables before building each package, using the standard SCons or Make-based build of each package. Additional tools are provided for top-level checkouts, updates, and clean builds. Testing of these tools has been carried out on 64-bit RedHat/CentOS 6 Linux systems. The user is assumed to already have SVN, Git, python 2.7, cURL, SCons, GNU Make, and CMake (for geant4) installed.

## Listed packages
`xerces-c`, `cernlib`, `root`, `amptools`, `geant4`, `evio`, `ccdb`, `jana`, `hdds`, `sim-recon`

## Build-settings templates
Builds are configured with text files which can serve as templates for future builds. Example templates are provided in the **"example-templates"** directory. Each subdirectory in the **"example-templates"** directory, with name **"settings-[id]"**, is a build-settings template. **"example-templates"** is copied to the **"templates"** directory when selecting a template for the first time. This directory is for storing user-defined templates. Create a new template by editing the text files located in the **"settings"** directory and copying them to a directory named **"templates/settings-[id]"**. The `hdpm save` command is available to make this a bit easier, and is described in the user interface section.

1. **top.txt**
   - **path** of the top-level build directory and **build tag** separated by a tab or space; Supply an absolute path in order to set the build directory to an arbitrary location; A relative path will be joined with the **"pkgs"** directory in the current working directory. The build tag is used as the name of the install-directory (instead of the usual *BMS_OSNAME*). Unique tags can be used to house side-by-side builds of the package with different dependencies or different source code, such as after switching branches in a Git repository. All packages built by `hdpm` have a file written to the install-directory as a record when the build finishes, named **"success.hdpm"**. It stores build statistics and dependencies.
   - alternatively, if you want to use the *defaults*, use **"default"** for the path and/or tag name; the *defaults* are a directory named **"pkgs"** in the current working directory and *BMS_OSNAME* as the installation directory name.
   - **note:** it is often convenient to make the *build tag* the same as the *template id*; however, this is not mandatory.
2. **paths.txt**
   - path of each package (**format:** `<name> <path>`, 1 package per line); if a relative path is given, it will be joined with the path of the top directory before being used. If the path contains the version number and/or *BMS_OSNAME*, they can be replaced by **"[VER]"** and/or **"[OS]"**, respectively. To exclude a package from the build environment, set the path to **"NA"** (non-applicable). Core dependencies cannot be excluded, so this is only relevant to `amptools` and `geant4` at this time.
   - **important:** do not change the default names of the standard packages; these names cannot be customized.
3. **urls.txt**
   - URL of each SVN, Git, or other package. If the URL contains the version number of a package, it can be replaced with **"[VER]"**.
4. **vers.txt**
   - version, SVN revision, or Git tag/branch of each package. Use this setting to control which version of each package is checked-out/downloaded and built. For SVN and Git packages use **"latest"** to get the most recent revision. This will get the *master* branch for Git.
5. **commands.txt**
   - list of build/configure commands to run for each package, including configuration options and number of threads to use in each build. Specify one command per line, with multi-word commands enclosed in double quotes. Unlike the other configuration files, this one supports multiple lines per package. To disable a package from being acquired, updated, and/or built, comment out each line for that package. Note, however, that unless the path of a disabled-dependency is set to a valid external installation, the build of the dependent package will fail. All of the essential SCons and Make-based builds are supported at this time; However, for **cernlib** just the 2005 version of the Vogt-64-bit build is supported, and cannot be customized. Setting its command(s) will have no effect, but it should remain uncommented for it to be enabled.

## Julia language
Julia is needed to run the package management scripts. If working on a 64-bit Linux machine on the JLab CUE (such as iFarm65), one can simply use the group copy of the Julia binary. Otherwise, you will need to download the Julia binary or build it from source. It is recommended to download the binary instead of building it from source, since it takes a very long time to build its LLVM dependencies.

### Quick start options on the JLab CUE
A setup script is provided which will put the 64-bit Linux group installation of Julia into your PATH on the JLab CUE. This script also makes an alias to the `hdpm` script (see below).
 - `source setup.(c)sh`

### Download options
1. [Julia binaries](http://julialang.org/downloads)
   - `source setup.(c)sh` (64-bit Linux download/setup)
   - `source setup-osx.(c)sh` (Mac OS X download/setup)
2. [Julia source at GitHub](https://github.com/JuliaLang/julia)

## User interface
Julia scripts (located in the **"src"** directory) are used to prepare, manage, and execute top-level builds. For typical usage, the user should not need to modify these. The scripts are controlled through the `hdpm` command interface.

* **hdpm**
   - Unified interface for managing packages. For convenience, setup.(c)sh creates an alias named `hdpm`  (hdpm='julia src/hdpm.jl').
   - **usage:** `hdpm <command> |<args>|`
   - commands: `help`, `select`, `save`, `show`, `fetch`, `build`, `update`, `clean`, `clean-build`
* **hdpm help**
    - Display available commands or arguments for a particular command.
    - **usage:** `hdpm help [command]`
* **hdpm select**
   - Select the settings template specified by the identifier **id** for your next build; all scripts will use the settings which have been copied from the templates directory by running this command.
   - **usage:** `hdpm select <id>`
* **hdpm save**
    - Save the current build settings as a new template for future builds.
    - **usage:** `hdpm save <new id>`
* **hdpm show**
   - Show the current build settings. 2 optional arguments specify which column of settings to show (**"name"**, **"version"**, **"url"**, **"path"**,**"deps"**, or **"cmds"**) and/or the integer number of spaces between columns. Use no arguments to show the first four columns of settings with the default spacing (2 spaces).
   - **usage:** `hdpm show [column name] [column spacing]`
* **hdpm fetch**
   - Checkout/clone SVN and Git packages; Download others using **curl**.
   - **usage:** `hdpm fetch [pkgs...]`
* **hdpm build**
   - Build all selected packages (fetch them if needed). All packages will be built if no arguments are given. There are also options for selecting the template or package(s) to build. Dependencies will be built if not already available. If a package has already been built, this command will display information about the build including how long it took to build, disk use, and the timestamp of when it finished.
   - **usage:** `hdpm build [id]`
   - **usage:** `hdpm build [pkgs...]`
* **hdpm update**
   - Update SVN and Git packages. For Git packages not set to *latest* (*master* branch), this will checkout and switch to a local branch denoted by the version.
   - **usage:** `hdpm update [pkgs...]`
* **hdpm clean**
   - Completely remove build products of selected packages. Only the  *ccdb*, *jana*, *hdds*, and *sim-recon* packages are currently supported.
   - **usage:** `hdpm clean [pkgs...]`
* **hdpm clean-build**
   - Do a clean build of the packages. This is normally used after running `hdpm update`. It will first delete your old executables, includes, and shared libs before building. This is currently supported for the *ccdb*, *jana*, *hdds*, and *sim-recon* packages.
   - **usage:** `hdpm clean-build [pkgs...]`

## Julia modules
The package management scripts depend on these Julia modules. For typical usage, the user should not need to modify these.

1. **Environs.jl**
   - Sets the environment required by the various Hall-D package builds. A C-shell script and bash script, for setting the environment variables, are saved to the **"[top]/env-setup"** directory; *source* the appropriate one before using the packages.
2. **Packages.jl**
   - Provides a composite type and various functions for organizing package information and build settings.

## Example templates
Example templates are provided in the **"example-templates"** directory. The **"F14"** and **"Sp15"** templates, for Fall-2014 and Spring-2015 data, respectively, are designed for use on the iFarm at Jefferson Lab, and use the official group installations of all packages except for **hdds** and **sim-recon**. The **"all"** template should be a suitable starting point for other use cases. It is used to install/manage all listed Hall-D packages, but can also be modified in order to utilize pre-existing/external installations if desired.

## Limitations on JLab iFarm
Almost all external HTTP traffic is blocked on the JLab iFarm, preventing one from directly downloading most packages which are hosted externally. These packages include **xerces-c**, **cernlib**, **root**, **amptools**, and **geant4**. Git HTTP traffic is **not** blocked on the iFarm; Clone JLab GitHub repos by running `git clone https://github.com/JeffersonLab/[name]`. In practice these limitations are not very constraining since group installations of almost all of these packages are available, and it usually makes more sense to use these instead. If absolutely necessary one can, for example, download the package(s) to their scratch folder using another machine at JLab not behind the HTTP firewall, such as jlabl3 or jlabl4, and then move it to the desired build directory.
