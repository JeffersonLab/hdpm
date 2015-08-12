using Packages
top = gettop()
for pkg in get_packages()
    @osx_only if name(pkg) == "cernlib" info("Mac OS X detected: skipping cernlib");continue end
    if tobuild(pkg) && ispath(path(pkg)) info(path(pkg)," exists") end
    if tobuild(pkg) && !ispath(path(pkg))
        println()
        mk_cd(top)
        URL = url(pkg)
        # checkout svn and git packages
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
            if name(pkg) == "amptools"
                get_unpack_file(URL,dirname(path(pkg)))
            else
                get_unpack_file(URL,path(pkg))
            end
        end
        if name(pkg) == "cernlib" && version(pkg) == "2005"
            mkdir(path(pkg))
            cd(path(pkg))
            get_unpack_file(replace(URL,".2005.corr.2014.04.17","-2005-all-new")) # get the "all" file
            get_unpack_file(replace(URL,"corr","install")) # get the "install" file
            run(`curl -O $URL`) # get the "corr" file
            run(`mv -f $(basename(URL)) cernlib.2005.corr.tgz`)
        end
    end
end
