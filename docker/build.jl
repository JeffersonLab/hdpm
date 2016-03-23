# build sim-recon in Docker containers
dtags = Dict("c6"=>"centos6","c7"=>"centos7","u14"=>"ubuntu14","f22"=>"fedora22")
bstages = ["base","deps","sim-recon"]
info("Available OS tags: ",join(keys(dtags),", "))
info("Build stages: ",join(bstages,", "))
for arg in ARGS
    if !(arg in keys(dtags)) && !(arg in bstages)
        println("Usage error: Improper argument ",arg,":\n\tMust be in OS tags or build stages."); exit()
    end end
function get_list(c,args)
    list = ASCIIString[]
    if length(args) > 0
        for item in c
            if item in args push!(list,item) end end
        if length(list) == 0 list = c end
    else list = c end
    list
end
tags = get_list(keys(dtags),ARGS)
stages = get_list(bstages,ARGS)
for tag in tags
    dir="buildfiles/$(dtags[tag])"
    for stage in stages
        if stage in ["base","deps"] name = string("hd",stage); dfile = string("Dockerfile-",stage)
        else name = stage; dfile = "Dockerfile" end
        try run(`docker rmi $name:$tag`)
        catch info("Image not available to remove (ignore previous error).") end
        f = open(joinpath(pwd(),".log-$name-$tag"),"w")
        write(f,readall(`docker build --no-cache -t $name:$tag -f $dir/$dfile $dir`)); close(f)
    end
end
