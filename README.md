# hdpm - Hall D Package Manager
This is a set of tools for managing **top-level builds** of the Jefferson Lab Hall-D offline software and online monitoring plugins; Specify a list of paths to the software and its dependencies, and it will set the Hall-D environment variables before building each package, using the standard SCons or Make-based build of each package. Additional tools are provided for top-level checkouts, updates, and clean builds. Testing of these tools has been carried out on 64-bit RedHat/CentOS 6 Linux systems. The user is assumed to already have SVN, Git, python 2.7, Wget or cURL, SCons, GNU Make, and CMake (for clhep and geant4) installed. 

## Listed packages
`python`, `xerces-c`, `cernlib`, `root`, `clhep`, `amptools`, `geant4`, `evio`, `ccdb`, `jana`, `hdds`, `sim-recon`, `online-monitoring`, `online-sbms`, `scripts(analysis)`

## Build-settings templates
Builds are configured with text files which can serve as templates for future builds. Example templates are provided in the **"example-templates"** directory. Each subdirectory in the **"example-templates"** directory, with name **"settings-[id]"**, is a build-settings template. **"example-templates"** is copied to the **"templates"** directory when selecting a template for the first time. This directory is for storing user-defined templates. Create a new template by editing the text files located in the **"settings"** directory and copying them to a directory named **"templates/settings-[id]"**.

1. **top.txt**
   - **path** of the top-level build directory and **build tag** separated by a tab or space; Supply an absolute path in order to set the build directory to an arbitrary location; A relative path will be joined with the **"pkgs"** directory in the current working directory. The build tag is appended to the *BMS_OSNAME* environment variable in order to support side-by-side builds.
   - alternatively, if you want to use the *defaults*, use **"default"** for the path and/or tag name; the *defaults* are a directory named **"pkgs/[date]"** in the current working directory and **no tag**.
   - **note:** it is often convenient to make the *build tag* the same as the *template id*; however, this is not mandatory.
2. **paths.txt**
   - path of each package (**format:** `<name> <path>`, 1 package per line); if a relative path is given, it will be joined with the path of the top directory before being used. If the path contains the version number and/or *BMS_OSNAME*, they can be replaced by **"[VER]"** and/or **"[OS]"**, respectively. To exclude a package from the build environment, set the path to **"NA"** (non-applicable).
   - **important:** do not change the default names of the standard packages; these names cannot be customized without modifying the Julia modules (see below).
3. **urls.txt**
   - URL of each SVN, Git, or other package. If the URL contains the version number of a package, it can be replaced with **"[VER]"**.
4. **vers.txt**
   - version, SVN revision, or Git tag/branch of each package. Use this setting to control which version of each package is checked-out/downloaded and built. For SVN and Git packages use **"latest"** to get the most recent revision. 
5. **tobuild.txt**
   - which packages to checkout, update, and/or build; use **"true"** or **"false"**. All of the essential SCons and Make-based builds are supported at this time; However, for **cernlib** just the 2005 version of the Vogt-64-bit build is supported. SVN and Git packages which require no building should be set to **"true"** for checkouts and updates.  
6. **nthreads.txt**
   - number of threads to use in each build.

## Julia language
Julia is needed to run the package management scripts. If working on a 64-bit Linux machine on the JLab CUE (such as iFarm65), one can simply use the group copy of the Julia binary. Otherwise, you will need to download the Julia binary or build it from source. It is recommended to download the binary instead of building it from source, since it takes a very long time to build its LLVM dependencies.

### Quick start options on the JLab CUE
Put the latest 64-bit Linux Julia executable into your path by doing one of the following:

1. `ln -s /u/group/halld/Software/ExternalPackages/julia-latest/bin/julia ~/bin/julia`
2. `setenv PATH /u/group/halld/Software/ExternalPackages/julia-latest/bin:$PATH`

### Download options
1. [Julia binaries](http://julialang.org/downloads)
   - `source get_julia.csh` (downloads 64-bit Linux binary using **wget** and puts it in your path) 
   - **note:** this will not work on the JLab iFarm because external HTTP traffic is blocked (see below)
2. [Julia source at GitHub](https://github.com/JuliaLang/julia)

## Package management scripts
The following Julia scripts are used to prepare, manage, and execute top-level builds. For typical usage, the user should not need to modify these.

1. **select_template.jl**
   - Select the settings template specified by the identifier **id** for your next build; all scripts will use the settings which have been copied from the templates directory by running this script.
   - **usage:** `julia select_template.jl <id>`
2. **show_settings.jl**
   - Show the current build settings. 2 optional arguments specify which column of settings to show (**"name"**, **"version"**, **"url"**, **"path"**, **"nthreads"**, or **"tobuild"**) and/or the integer number of spaces between columns. Use no arguments to show all settings with the default spacing (8 spaces).
   - **usage:** `julia show_settings.jl [column name] [column spacing]`
3. **copkgs.jl**
   - Checkout SVN and Git packages; Download others using **wget** or **curl**; use curl if wget is not available.
   - **usage:** `julia copkgs.jl`
4. **mkpkgs.jl**
   - Build all selected packages.
   - **usage:** `julia mkpkgs.jl`
5. **update.jl**
   - Update SVN and Git packages. For Git packages not set to the *latest* revision, this will checkout the tag/branch denoted by the revision.
   - **usage:** `julia update.jl`
6. **clean_build.jl**
   - Do a clean build of the packages. This is normally used after running the *update.jl* script. It will first delete your old executables, includes, and shared libs before building. This is currently supported for the *ccdb*, *jana*, *hdds*, *sim-recon*, and *online-monitoring* packages.
   - **usage:** `julia clean_build.jl`

## Julia modules
The package management scripts depend on these Julia modules. For typical usage, the user should not need to modify these.

1. **Environs.jl**
   - Sets the environment required by the various Hall-D package builds. A C-shell script and bash script, for setting the environment variables, are saved to the **"[top]/scripts/env"** directory; *source* the appropriate one before using the packages.
2. **Packages.jl**
   - Provides a composite type and various functions for organizing the package information and build settings.

## Recommended usage
It is possible to build all the packages in one go, within a single **top** directory. However, for most users this is probably not the best usage. I recommmend two patterns of usage. One is to utilize group installations of packages which are relatively stable when they are available. This helps to conserve computing resources and means less work for you. Make a template with the paths to the group installations and just add your own for the ones you want to build yourself, like **hdds**, **sim-recon**, and **online-monitoring**. Second is a two stage approach; first build the relatively stable packages; choose any convenient **top** directory for these. Maybe you will want to share these installations. For the second stage, make a separate build template for development builds of hdds, sim-recon, and online-monitoring; set the **top** to your development directory. Example templates are provided in the **"example-templates"** directory.

## Limitations on JLab iFarm
External HTTP traffic is blocked on the JLab iFarm, preventing one from directly downloading packages which are hosted externally. These packages include **xerces-c**, **cernlib**, **root**, **clhep**, **amptools**, and **geant4**. Git protocol traffic is **not** blocked on the iFarm, allowing one to clone the *ccdb* and *hdpm* GitHub repos by running `git clone git://github.com/JeffersonLab/[name]`. However, the Git protocol does not support pushing changes back to the remote: it is read-only. In practice these limitations are not very constraining since group installations of almost all of these packages are available, and it usually makes more sense to use these instead. If absolutely necessary one can, for example, download the package(s) to their scratch folder using another machine at JLab not behind the HTTP firewall, such as jlabl3 or jlabl4, and then move it to the desired build directory.
