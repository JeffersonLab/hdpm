using Packages
const top = gettop()
const dist_dir = joinpath(top,".dist")
hz("="); info("Linking binaries into $top ...")
if !ispath(dist_dir) usage_error("No binary packages available (run 'hdpm install').") end
os = replace(osrelease(),"RHEL","CentOS")
os = replace(os,"LinuxMint17","Ubuntu14")
os = replace(os,"LinuxMint18","Ubuntu16")
commit(a,i) = split(split(readall("$dist_dir/$a/$os/success.hdpm"))[1],"-")[i]
for pkg in get_packages()
    n = name(pkg); v = version(pkg)
    if n=="hdds" v = commit(n,2) end
    if n=="sim-recon" v = commit(n,3) end
    p = joinpath(top,n,v)
    pd = n=="sim-recon" || n=="hdds" ? joinpath(dist_dir,n):joinpath(dist_dir,string(n,"-",v))
    if !ispath(pd) println("\t$n-$v is not included in distribution."); continue end
    if ispath(p) println("\t$n-$v is already installed."); continue end
    mkpath(dirname(p))
    if n == "hdds" || n == "sim-recon"
        d = dirname(p)
        for dir in readdir(d); if !islink(joinpath(d,dir)) continue end
            rm(joinpath(d,dir))
        end
    end
    run(`ln -s $pd $p`)
end
run(`rm -f $top/env-setup/dist.sh`); run(`rm -f $top/env-setup/dist.csh`)
if ispath("$dist_dir/env-setup/hdenv.sh") run(`ln -s $dist_dir/env-setup/hdenv.sh $top/env-setup/dist.sh`) end
if ispath("$dist_dir/env-setup/hdenv.csh") run(`ln -s $dist_dir/env-setup/hdenv.csh $top/env-setup/dist.csh`) end
