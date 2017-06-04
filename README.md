# hdpm - Hall D Package Manager
`hdpm` is a tool for building Hall-D offline simulation and reconstruction software (aka [sim-recon](https://github.com/JeffersonLab/sim-recon)) and its dependencies. Specify a list of paths to sim-recon and its dependencies, and it will set the Hall-D environment variables before building each package, using the standard SCons or Make-based build system of each package. Additional commands are provided for downloads, updates, and clean builds. Packages can be built individually or all in turn, with dependency builds being triggered when needed. `hdpm` can be used to install precompiled packages on CentOS 6/7 and Ubuntu 16.04 LTS.

## Default packages
`xerces-c`, `cernlib`, `root`, `amptools`, `geant4`, `evio`, `rcdb`, `ccdb`, `jana`, `hdds`, `sim-recon`, `hdgeant4`, `gluex_root_analysis`, `hd_utilities`

## Documentation
For documentation on how to install and use `hdpm`, visit its wiki at https://github.com/JeffersonLab/hdpm/wiki.

## Tested platforms (x86-64)
- CentOS/RHEL 6
- CentOS/RHEL 7
- Ubuntu 16.04 LTS
