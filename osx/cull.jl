# delete dist. tarfiles corresponding to commits older than the latest Nkeep
target = joinpath(pwd(),".pkgs")
if length(ARGS) != 2 println("Usage error: Wrong number of arguments.\n\tFirst argument is number of files to keep, second is OS tag."); exit() end
Nkeep = ARGS[1]
cwd = pwd(); cd(".sim-recon")
commits = split(readchomp(`git log --grep="Merge pull request #" -$Nkeep --format="%h"`))
flist = filter(r"^sim-recon-.{7}-.{5}-.{2,3}.tar.gz$",readdir(target))
if length(ARGS) == 2 filter!(x -> contains(x,string("-",ARGS[2],".tar.gz")),flist) end
for file in flist
    c = split(file,"-")[3]
    if !(c in commits) rm(joinpath(target,file)) end
end
fbad = filter(r"^sim-recon--.{5}-.{2,3}.tar.gz$",readdir(target))
for file in fbad rm(joinpath(target,file)) end
cd(cwd)
