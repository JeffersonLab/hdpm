# delete dist. tarfiles corresponding to commits older than the latest Nkeep
target = joinpath(pwd(),".pkgs")
if length(ARGS) == 0 || length(ARGS) > 2 error("Wrong number of arguments. First argument is number of commits to keep; second, optional one, is OS tag.") end
Nkeep = ARGS[1]
if parse(Int,Nkeep) == 0 || parse(Int,Nkeep) > 10 info("Number of commits to keep ($Nkeep) is 0 or larger than 10.") end
cwd = pwd(); cd(".sim-recon")
commits = split(readchomp(`git log --grep="Merge pull request #" -$Nkeep --format="%h"`))
println("$Nkeep latest commits: ",commits)
flist = filter(r"^sim-recon-.{7}-.{5}-.{2,3}.tar.gz$",readdir(target))
if length(ARGS) == 2 filter!(x -> contains(x,ARGS[2]),flist) end
for file in flist
    c = split(file,"-")[3]
    if !(c in commits) rm(joinpath(target,file)) end
end
fbad = filter(r"^sim-recon--.{5}-.{2,3}.tar.gz$",readdir(target))
for file in fbad rm(joinpath(target,file)) end
cd(cwd)
