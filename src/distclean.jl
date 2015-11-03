using Packages
# clean packages for distribution
safe_to_proceed =
"yes" in ARGS ? true : input("You are about to obliterate your src code!
Do you really want to proceed (yes or no)? ") == "yes"
if safe_to_proceed
    top = gettop()
    cwd = pwd(); cd(joinpath(top,"sim-recon"))
    info("Asserting that working directory of sim-recon is clean ...")
    assert(readall(`git status`) == "On branch master\nYour branch is up-to-date with 'origin/master'.\nnothing to commit, working directory clean\n"); cd(cwd)
    rm_regex(r"^julia-.+",top)
    BMS_OSNAME = install_dirname()
    for pkg in get_packages(); if length(ARGS) > 0 if !(name(pkg) in ARGS) continue end end
        if !is_external(pkg) && ispath(path(pkg))
            cd(path(pkg))
            run(`rm -rf src`); run(`rm -rf .$BMS_OSNAME`)
            rm_regex(r".+gz$"); rm_regex(r".+\.contents$")
            rm_regex(r"^\.g.+"); rm_regex(r"^\.s.+")
            rm_regex(r"^setenv\..+")
            rm_regex(r"^setenv\..+",joinpath(pwd(),BMS_OSNAME))
        end
    end
end
