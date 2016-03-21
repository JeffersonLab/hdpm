dtags = Dict("c6"=>"centos6","c7"=>"centos7","u14"=>"ubuntu14","f22"=>"fedora22")
info("Available OS tags: ",join(keys(dtags),", "))
if length(ARGS) == 0 println("Usage error: Please provide Docker username as first argument.\n\tSpecify a subset of tags by listing them as additional arguments."); exit() end
duser = ARGS[1]
name="sim-recon-deps"
for tag in keys(dtags); if length(ARGS) > 1 if !(tag in ARGS) continue end end
    repo = (tag == "c6" || tag == "c7") ? string(joinpath("quay.io",duser,name),":",dtags[tag]) : string(joinpath(duser,name),":",dtags[tag])
    f = open(joinpath(pwd(),".id-deps-$tag"),"w")
    try
        write(f,split(readchomp(`docker inspect --format='{{.Id}}' $repo`),":")[2][1:5])
    catch
        info(repo," does not exist.")
        rm(joinpath(pwd(),".id-deps-$tag"))
    finally
        close(f)
    end
end
