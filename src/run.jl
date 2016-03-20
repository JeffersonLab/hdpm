using Environs,Packages
# run a command in the Hall-D offline environment
@osx_only error("The 'hdpm run' command is not supported on Mac OS X.\n")
printenv() # expose env variables and save them to env-setup file
ts = readchomp(`date "+%Y-%m-%d"`)
mk_cd("work/ByDate/$ts")
ENV["CCDB_USER"] = ENV["USER"]
if length(ARGS) > 1 error("'hdpm run' does not support multiple arguments. \n\t\t\tPut multi-word commands in double quotes.\n") end
cmd = (length(ARGS) == 1) ? ARGS[1] : "bash"
run(`sh -c $cmd`)
