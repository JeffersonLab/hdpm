if length(ARGS) != 1 error("Please provide Docker username as single argument") end
duser = ARGS[1]
name="sim-recon-deps"
short_tags = ["c6","c7","u14","f22"]
i=0
for tag in ["centos6","centos7","ubuntu14","fedora22"]
    i+=1
    repo = (i < 3) ? string(joinpath("quay.io",duser,name),":",tag) : string(joinpath(duser,name),":",tag)
    f = open(joinpath(pwd(),".id-deps-$(short_tags[i])"),"w")
    try
        write(f,readchomp(`docker inspect --format='{{.Id}}' $repo`)[1:5])
    catch
        info(repo," does not exist")
        rm(joinpath(pwd(),".id-deps-$(short_tags[i])"))
    finally
        close(f)
    end
end
