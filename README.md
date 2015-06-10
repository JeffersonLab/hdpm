# HDBuildManager
This is a set of tools for managing **top-level builds** of the Jefferson Lab Hall-D offline software and online monitoring plugins; Specify a list of paths to the software and its dependencies, and it will set the Hall-D environment variables before building each package, using the standard SCons or Make-based build of each package. Additional tools are provided for top-level checkouts, updates, and clean builds. Testing of these tools has been carried out on 64-bit RedHat/CentOS 6 Linux systems. The user is assumed to already have SVN, Git, python 2.7, Wget or cURL, and CMake (for clhep and geant4) installed.

## Listed packages
python, xerces-c, cernlib, root, clhep, amptools, geant4, evio, ccdb, jana, hdds, sim-recon, online-monitoring, online-sbms, scripts(analysis)

## Build-settings templates
Each subdirectory in the **"templates"** directory, with name **"settings_[id]"**, is a build-settings template. Create new templates by copying and editing the text files to specify the desired build settings.

1. **top.txt**
   - **absolute path** of the top-level build directory and **build tag** name separated by a tab or space; this tag is appended to the *BMS_OSNAME* environment variable in order to support side-by-side builds.
   - alternatively, if you want to use the *defaults*, use **"default"** for the path and/or tag name; the *defaults* are a directory named **"build_[date]"** in the current working directory and **no tag**.
   - **note:** it is often convenient to make the *build tag* the same as the *template id*; however, this is not mandatory.
2. **paths.txt**
   - path of each package (**format:** [name] [path], 1 package per line); if a relative path is given, it will be joined with the path of the top directory before being used.
   - **important:** do not change the default names of the standard packages; these names cannot be customized without modifying the Julia modules (see below).
3. **urls.txt**
   - URL of each SVN, Git, or other package.
4. **vers.txt**
   - version, SVN revision, or Git tag/branch of each package. For SVN and Git packages use **"latest"** to get the most recent revision. This feature is only used for SVN and Git packages at this time, except for **cernlib** for setting the *CERN_LEVEL* environment variable.  
5. **tobuild.txt**
   - which packages to checkout, update, and/or build; use **"true"** or **"false"**. All of the essential SCons and Make-based builds are supported at this time; However, for **cernlib** just the 2005 version of the Vogt-64-bit build is supported. SVN and Git packages which require no building should be set to **"true"** for checkouts and updates.  
6. **nthreads.txt**
   - number of threads to use in each build.

## Julia language
Julia is needed to run the build management scripts. If working on the JLab ifarm, one can simply use my copy of the Julia binary. Otherwise, you will need to download the Julia binary or build it from source. It is recommended to download the binary instead of building it from source, since it takes a very long time to build its LLVM dependencies.

###Quick start options on the ifarm at JLab
1. Make an alias to my copy of the Julia binary.
   - alias julia /w/halld-scifs1a/nsparks/packages/julia-0.3.9/bin/julia
2. Make a symbolic link to it.
   - ln -s /w/halld-scifs1a/nsparks/packages/julia-0.3.9/bin/julia ~/bin/julia

###Download options
1. [Julia binaries](http://julialang.org/downloads)
   - **source get_julia.csh** (downloads 64-bit Linux binary using **wget** and puts it in your path) 
2. [Julia source at GitHub](https://github.com/JuliaLang/julia)

## Top-level build management scripts
The following Julia scripts are used to prepare, manage, and execute top-level builds. For typical usage, the user should not need to modify these.

1. **select_template.jl**
   - Select the settings template specified by the identifier **id** for your next build; all scripts will use the settings which have been copied from the templates directory by running this script.
   - **usage:** julia select_template.jl [id]
2. **show_settings.jl**
   - Show the current build settings.
   - **usage:** julia show_settings.jl
3. **copkgs.jl**
   - Checkout SVN and Git packages; Download others using **wget** or **curl**; use curl if wget is not available.
   - **usage:** julia copkgs.jl
4. **mkpkgs.jl**
   - Build all selected packages.
   - **usage:** julia mkpkgs.jl
5. **update.jl**
   - Update SVN and Git packages. For Git packages not set to the *latest* revision, this will checkout the tag/branch denoted by the revision.
   - **usage:** julia update.jl
6. **clean_build.jl**
   - Do a clean build of the packages. This is normally used after running the *update.jl* script. It will first delete your old executables, includes, and shared libs before building. This is currently supported for the *ccdb*, *jana*, *hdds*, *sim-recon*, and *online-monitoring* packages.
   - **usage:** julia clean_build.jl

## Julia modules
The build management scripts depend on these Julia modules. For typical usage, the user should not need to modify these.

1. **Environs.jl**
   - Sets the environment required by the various Hall-D package builds. A C-shell script and bash script, for setting the environment variables, are saved to the **"[top]/scripts/env"** directory; *source* the appropriate one before using the packages.
2. **Packages.jl**
   - Provides a composite type and various functions for organizing the package information and build settings.

## Recommended usage
It is possible to build all the packages in one go, within a single **top** directory. However, for most users this is probably not the best usage. I recommmend two patterns of usage. One is to utilize group installations of packages which are relatively stable when they are available. This helps to conserve computing resources and means less work for you. Make a template with the paths to the group installations and just add your own for the ones you want to build yourself, like **hdds**, **sim-recon**, and **online-monitoring**. Second is a two stage approach; first build the relatively stable packages; choose any convenient **top** directory for these. Maybe you will want to share these installations. For the second stage, make a separate build template for development builds of hdds, sim-recon, and online-monitoring; set the **top** to your development directory. Example templates are provided in the **templates** directory.
