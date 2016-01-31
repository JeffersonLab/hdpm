top="../pkgs/osx"
os = readchomp(`perl ../src/osrelease.pl`)
dir = filter(r"jana-.+",readdir(top))[1]
s = split(split(readall("$top/$dir/$os/success.hdpm"))[2],"-")
id_deps = string(s[1][4],s[2],s[3][1:2])
f = open(joinpath(pwd(),".id-deps-osx"),"w")
write(f,id_deps)
close(f)
