using Packages
# clean packages for distribution
safe_to_proceed =
"yes" in ARGS ? true : input("You are about to obliterate your src code!
Do you really want to proceed (yes or no)? ") == "yes"
if safe_to_proceed
    top = gettop()
    if ispath(joinpath(top,"sim-recon","master"))
        cwd = pwd(); cd(joinpath(top,"sim-recon","master"))
        info("Asserting that working directory of sim-recon is clean ...")
        assert(contains(readall(`git status`),"working directory clean")); cd(cwd); rm_regex(r"^julia-.+",top)
    end
    for pkg in get_packages()
        if length(ARGS) > 0 && !(length(ARGS) == 1 && ARGS[1] == "yes")
            if !(name(pkg) in ARGS) continue end end
        if !is_external(pkg) && ispath(path(pkg))
            cd(path(pkg))
            if name(pkg) == "root" && contains(cmds(pkg)[1],"./configure")
                if !ispath("Makefile") continue end
                run(`cp -p success.hdpm ../`)
                run(`make dist`); cd("../")
                run(`rm -rf $(version(pkg))`)
                for item in filter(r".+gz$",readdir("."))
                    run(`tar xf $item`); rm(item)
                end
                run(`mv success.hdpm $(version(pkg))`)
            else
                BMS_OSNAME = install_dirname()
                run(`rm -rf src`); run(`rm -rf .$BMS_OSNAME`)
                rm_regex(r".+gz$"); rm_regex(r".+\.contents$")
                rm_regex(r"^\.g.+"); rm_regex(r"^\.s.+")
                rm_regex(r"^setenv\..+")
                rm_regex(r"^setenv\..+",joinpath(pwd(),BMS_OSNAME))
            end
        end
    end
end
