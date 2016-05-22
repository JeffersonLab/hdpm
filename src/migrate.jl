using Packages
# change pkgs directory structure to <pkg>/<ver> from <pkg>-<ver>
const top = gettop()
for pkg in get_packages()
    n = name(pkg)
    if n == "cernlib" || is_external(pkg) continue end
    p = path(pkg)
    v = version(pkg)
    cd(top)
    if v == "master"
        if !ispath(n) warn("$n does not exist (skipping it)."); continue end
        if ispath(joinpath(n,"master")) info("Nothing to migrate for $n."); continue end
        info("Migrating $n to $n/master.")
        run(`mv $n master`)
        run(`mkdir -p $n`)
        run(`mv master $n`)
    else
        dir = string(n,"-",v)
        if !ispath(dir) warn("$dir does not exist (skipping it)."); continue end
        new_dir = joinpath(n,v)
        if ispath(new_dir) info("Nothing to migrate for $n."); continue end
        info("Migrating $dir to $new_dir.")
        run(`mkdir -p $n`)
        run(`mv $dir $new_dir`)
    end
end
