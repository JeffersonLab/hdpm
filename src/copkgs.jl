using Packages
top = gettop()
deps = get_deps(ARGS) # add deps
if length(deps) > 0 info("dependency list: ",string(deps)) end
for pkg in get_packages(); if length(ARGS) > 0 if !(name(pkg) in ARGS) && !(name(pkg) in deps) continue end end
    @osx_only if name(pkg) == "cernlib" info("Mac OS X detected: skipping cernlib"); continue end
    if is_external(pkg) && name(pkg) in deps warn(name(pkg)," is dependency under user control, assumed to be set to valid external installation.") end
    if !is_external(pkg) && ispath(path(pkg)) info(path(pkg)," exists") end
    if !is_external(pkg) && !ispath(path(pkg))
        println()
        mk_cd(top)
        URL = url(pkg)
        # checkout/clone svn and git packages
        if contains(URL,"svn")
            rev = version(pkg)
            if rev!="latest" && !contains(URL,"tags")
                run(`svn checkout --non-interactive --trust-server-cert -r $rev $URL $(path(pkg))`)
            else
                run(`svn checkout --non-interactive --trust-server-cert $URL $(path(pkg))`)
            end
        end
        if contains(URL,"git") && !contains(URL,"archive")
            run(`git clone $URL $(path(pkg))`)
            rev = version(pkg)
            if rev!="latest"
                cd(path(pkg))
                run(`git checkout -b $rev $rev`)
            end
        end
        # download/unpack other packages
        if name(pkg) != "cernlib" && (contains(URL,".tar.gz") || contains(URL,".tgz"))
            get_unpack_file(URL,path(pkg))
        end
        if name(pkg) == "cernlib" && version(pkg) == "2005"
            mk_cd(path(pkg))
            get_unpack_file(replace(URL,".2005.corr.2014.04.17","-2005-all-new")) # get the "all" file
            get_unpack_file(replace(URL,"corr","install")) # get the "install" file
            run(`curl -OL $URL`) # get the "corr" file
            run(`mv -f $(basename(URL)) cernlib.2005.corr.tgz`)
        end
    end
end
