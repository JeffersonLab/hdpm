# hdpm - Hall D Package Manager
`hdpm` is a tool for installing Hall-D offline simulation and reconstruction software (aka [sim-recon](https://github.com/JeffersonLab/sim-recon)) and its dependencies. The default package settings are suitable for most users but can be customized by editing simple JSON files. `hdpm` sets the required environment variables for builds based on the package settings, and installs each package using its standard SCons or Make-based build system. Various commands are provided for managing packages. Packages can be built individually or in groups, with dependency builds being triggered when needed. `hdpm` can be used to install precompiled packages on CentOS 6/7 and Ubuntu 16.04 LTS.

## Default packages
`xerces-c`, `cernlib`, `root`, `amptools`, `geant4`, `evio`, `rcdb`, `ccdb`, `jana`, `hdds`, `sim-recon`, `hdgeant4`, `gluex_root_analysis`, `hd_utilities`

## Documentation
For documentation on how to install and use `hdpm`, visit its wiki at https://github.com/JeffersonLab/hdpm/wiki.

## Tested platforms (x86-64)
- CentOS/RHEL 6
- CentOS/RHEL 7
- Ubuntu 16.04 LTS
