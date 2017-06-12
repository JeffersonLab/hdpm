# hdpm - Hall D Package Manager
`hdpm` is a tool for installing Hall-D offline simulation and reconstruction software (aka [sim-recon](https://github.com/JeffersonLab/sim-recon)) and its dependencies. It provides package management commands, such as **select**, **fetch**, **install**, and **clean**. Package builds can be customized by editing simple JSON files, for example, by changing package versions or options used by build commands. `hdpm` sets environment variables by using the **path** settings and builds each package using its standard SCons or Make-based build system. Packages can be built individually or in groups, with dependency builds being triggered when needed. Binary-package distributions are provided for CentOS 6, CentOS 7, and Ubuntu 16.04 LTS.

## Default packages
`xerces-c`, `cernlib`, `root`, `amptools`, `geant4`, `evio`, `rcdb`, `ccdb`, `jana`, `hdds`, `sim-recon`, `hdgeant4`, `gluex_root_analysis`, `hd_utilities`

## Documentation
For documentation on how to install and use `hdpm`, visit its wiki at https://github.com/JeffersonLab/hdpm/wiki.

## Tested platforms (x86-64)
- CentOS/RHEL 6
- CentOS/RHEL 7
- Ubuntu 16.04 LTS
