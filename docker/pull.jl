if length(ARGS) != 1 error("Please provide Docker username as single argument") end
duser = ARGS[1]
name="sim-recon"
short_tags = ["c6","c7","u14","f22"]
i=0
for tag in ["centos6","centos7","ubuntu14","fedora22"]
    i+=1
    repo = string(joinpath(duser,name),":",tag)
    run(`docker pull $repo`)
    try run(`docker rmi $name:$(short_tags[i])`)
    catch info("Image not available to remove (ignore previous 2 errors)") end
    run(`docker tag $repo $name:$(short_tags[i])`); run(`docker rmi $repo`)
end
