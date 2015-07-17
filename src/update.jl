using Packages
# update svn and git packages
for pkg in get_packages()
    if contains(url(pkg),"svn") && tobuild(pkg) 
        cd(path(pkg))
        rev = version(pkg)
        println(string("Updating ",name(pkg)," to svn revision $rev."))
        if rev!="latest"
            run(`svn update -r$rev`)
        else
            run(`svn update`)
        end
    end
    if contains(url(pkg),"git") && tobuild(pkg) 
        cd(path(pkg))
        rev = version(pkg)
        println(string("Updating ",name(pkg)," to git revision $rev."))
        if rev!="latest"
            run(`git checkout $rev`)
        else
            run(`git pull`)
        end
    end
end
