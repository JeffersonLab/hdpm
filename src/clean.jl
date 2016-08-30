using Environs,Packages
# clean packages
BMS_OSNAME = install_dirname()
for pkg in get_packages(); if length(ARGS) > 0 if !(name(pkg) in ARGS) continue end end
    if !is_external(pkg)
        cd(path(pkg))
        if name(pkg) == "ccdb"
            run(`rm -f .sconsign.dblite`)
            run(`rm -f success.hdpm`)
            if length(cmds(pkg)) != 1 warn("ccdb cannot be cleaned. Check 'commands.txt'"); continue end
            cmd = replace(cmds(pkg)[1],"scons","scons -c")
            run(setenv(`sh -c $cmd`,getenv()))
        end
        if name(pkg) in ["jana","hdds","sim-recon","gluex_root_analysis"]
            run(`rm -f .sconsign.dblite`)
            run(`rm -f src/.sconsign.dblite`)
            run(`rm -f success.hdpm`)
            run(`rm -rf $BMS_OSNAME`)
            if name(pkg) == "gluex_root_analysis"
                for dir in ["libraries/DSelector","programs/MakeDSelector","programs/tree_to_amptools"]
                    run(`rm -rf $dir/$BMS_OSNAME`)
                end
            end
            if ispath("src") cd("src") end; run(`rm -rf .$BMS_OSNAME`)
        end
    end
end
