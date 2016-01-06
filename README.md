# hdpm - Hall D Package Manager
`hdpm` is a tool for building Hall-D offline simulation and reconstruction software (aka [sim-recon](https://github.com/JeffersonLab/sim-recon)) and its dependencies. Specify a list of paths to sim-recon and its dependencies, and it will set the Hall-D environment variables before building each package, using the standard SCons or Make-based build system of each package. Additional commands are provided for downloads, updates, and clean builds. Packages can be built individually or all in turn, with dependency builds being triggered when needed. `hdpm` can also be used to install precompiled packages on Linux, speeding up the setup of a Hall-D offline software development environment.

## Listed packages
`xerces-c`, `cernlib`, `root`, `amptools`, `geant4`, `evio`, `ccdb`, `jana`, `hdds`, `sim-recon`

## Documentation
For documentation on how to setup and use `hdpm`, visit the wiki at https://github.com/JeffersonLab/hdpm/wiki.

## Tested platforms
- CentOS/RHEL 6
- CentOS/RHEL 7
- Mac OS X 10.10+
- Ubuntu 14
- Fedora 22
