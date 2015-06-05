using Packages
top = gettop()
for pkg in get_packages()
    if tobuild(pkg) && ispath(path(pkg)) println("Warning: ",path(pkg)," already exists. Skipping it.") end
    if tobuild(pkg) && !ispath(path(pkg))
        mk_cd(top)
        if name(pkg) in ["scripts","online-monitoring","online-sbms"] 
            mk_cd(getbase(path(pkg))) 
        end
        URL = url(pkg)
        # checkout svn and git packages
        if contains(URL,"svn")
            rev = version(pkg)
            if rev!="latest"
                run(`svn checkout -r $rev $URL`)
            else
                run(`svn checkout $URL`)
            end
        end
        if contains(URL,"git")
            run(`git clone $URL $(path(pkg))`)
            rev = version(pkg)
            if rev!="latest"
                cd(path(pkg))
                run(`git checkout -b $rev $rev`)
            end
        end
        # download/unpack other packages
        if name(pkg) in ["xerces-c","evio","amptools"]
            get_unpack_file(URL)
            if name(pkg) == "amptools" run(`mv AmpTools $(string("AmpTools_",version(pkg)))`) end
        end
        if name(pkg) == "cernlib" && version(pkg) == "2005"
            mkdir(path(pkg))
            cd(path(pkg))
            get_unpack_file(replace(URL,".2005.corr.2014.04.17","-2005-all-new")) # get the "all" file
            get_unpack_file(replace(URL,"corr","install")) # get the "install" file
            run(`wget $URL`) # get the "corr" file
            file = split(URL,"/")[end]; run(`mv -f $file cernlib.2005.corr.tgz`)
        end
    end
end
