using Environs,Packages
# build sim-recon directory
printenv() # expose env variables and save them to env-setup file
ts = readchomp(`date "+%Y-%m-%d"`)
mk_cd("work/ByDate/$ts")
ENV["CCDB_USER"] = ENV["USER"]
if length(ARGS) != 1 error("'hdpm run' requires a command to run as the argument.") end
cmd = ARGS[1]
run(`sh -c $cmd`)
