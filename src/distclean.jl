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
    BMS_OSNAME = install_dirname()
    for pkg in get_packages()
        if length(ARGS) > 0 && !(length(ARGS) == 1 && ARGS[1] == "yes")
            if !(name(pkg) in ARGS) continue end end
        if !is_external(pkg) && ispath(path(pkg))
            cd(path(pkg))
            cleanup(pkg)
        end
    end
end
