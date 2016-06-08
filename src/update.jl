using Packages
# update svn and git packages
for pkg in get_packages(); if length(ARGS) > 0 if !(name(pkg) in ARGS) continue end end
    if contains(url(pkg),"svn") && !contains(url(pkg),"tags") && !is_external(pkg)
        if !ispath(path(pkg)) continue end
        cd(path(pkg))
        rev = version(pkg)
        println(string("Updating ",name(pkg)," to svn revision $rev."))
        if rev!="master"
            run(`svn update -r$rev`)
        else
            run(`svn update`)
        end
    end
    if contains(url(pkg),"git") && !contains(url(pkg),"archive") && !is_external(pkg)
        if !ispath(path(pkg)) continue end
        cd(path(pkg))
        branch = version(pkg)
        println("Updating $branch branch of ",name(pkg),".")
        run(`git checkout $branch`)
        run(`git pull`)
    end
end
