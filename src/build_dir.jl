using Environs,Packages
# build sim-recon subdirectory
printenv() # expose env variables and save them to env-setup file
if length(ARGS) == 0 usage_error("'hdpm build-dir' requires a sim-recon subdirectory as the argument.") end
if length(ARGS) > 2 usage_error("Too many arguments: Run 'hdpm help build' to see available arguments.") end
dir = (ARGS[1] == "-c") ? ARGS[2] : ARGS[1]
cd(dir)
if !contains(dir,"sim-recon") usage_error("$dir does not appear to be a subdirectory of sim-recon.") end
cmd = (ARGS[1] == "-c") ? "scons -u -c install": "scons -u install"
run(`sh -c $cmd`)
