# build sim-recon in Docker containers
short_tags = ["c6","c7","u14","f22"]
i=0
info("Available OS tags: ",short_tags)
info("Build stages: ",["base","deps","sim-recon"])
for arg in ARGS
    if !(arg in vcat(short_tags,["base","deps","sim-recon"]))
        error("Improper argument ",arg,": Must be in OS tags or build stages")
    end end
for tag in ["centos6","centos7","ubuntu14","fedora22"]; if length(ARGS) > 0 if !(short_tags[i+1] in ARGS) continue end end
    i+=1
    dir="buildfiles/$tag"
    for stage in ["base","deps","sim-recon"]; if length(ARGS) > 0 if !(stage in ARGS) continue end end
        if stage in ["base","deps"] name = string("hd",stage); dfile = string("Dockerfile-",stage)
        else name = stage; dfile = "Dockerfile" end
        try run(`docker rmi $name:$(short_tags[i])`)
        catch info("Image not available to remove (ignore previous 2 errors)") end
        f = open(joinpath(pwd(),".log-$name-$(short_tags[i])"),"w")
        write(f,readall(`docker build --no-cache -t $name:$(short_tags[i]) -f $dir/$dfile $dir`)); close(f)
    end
end
