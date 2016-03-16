using Environs,Packages
# build sim-recon subdirectory
printenv() # expose env variables and save them to env-setup file
if length(ARGS) == 0 error("'hdpm build-dir' requires a sim-recon subdirectory as the argument.") end
if length(ARGS) > 2 error("Too many arguments") end
dir = (ARGS[1] == "-c") ? ARGS[2] : ARGS[1]
cd(dir)
if !contains(dir,"sim-recon") error("$dir does not appear to be a subdirectory of sim-recon.") end
cmd = (ARGS[1] == "-c") ? "scons -u -c install": "scons -u install"
run(`sh -c $cmd`)
