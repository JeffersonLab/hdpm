dtags = Dict("c6"=>"centos6","c7"=>"centos7","u14"=>"ubuntu14","f22"=>"fedora22")
info("Available OS tags: ",keys(dtags))
if length(ARGS) == 0
    error("Please provide Docker username as first argument.
    Specify a subset of tags by listing them as additional arguments.") end
duser = ARGS[1]
name="sim-recon-deps"
for tag in keys(dtags); if length(ARGS) > 1 if !(tag in ARGS) continue end end
    repo = (tag == "c6" || tag == "c7") ? string(joinpath("quay.io",duser,name),":",dtags[tag]) : string(joinpath(duser,name),":",dtags[tag])
    run(`docker tag hddeps:$tag $repo`)
    run(`docker push $repo`); run(`docker rmi $repo`)
end
