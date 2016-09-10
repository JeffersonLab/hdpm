using Environs,Packages
# build packages
const home = dirname(dirname(@__FILE__))
printenv() # expose env variables and save them to env-setup file
BMS_OSNAME = install_dirname()
deps = get_deps(ARGS) # add deps
first_success = ""
for pkg in get_packages(); if length(ARGS) > 0 if !(name(pkg) in ARGS) && !(name(pkg) in deps) continue end end
    @osx_only if name(pkg) == "cernlib" continue end
    path_to_success = joinpath(path(pkg),"success.hdpm")
    if name(pkg) in ["jana","hdds","sim-recon"] path_to_success = joinpath(path(pkg),BMS_OSNAME,"success.hdpm") end
    if first_success == "" && ispath(path_to_success) && !is_external(pkg) first_success = name(pkg) end
    if !is_external(pkg) && !ispath(path_to_success)
        @assert(ispath(path(pkg)),"$(path(pkg)) does not exist")
        println();info("$(name(pkg)): checking dependencies")
        check_deps(pkg)
        info("building $(name(pkg))-$(git_version(pkg))")
        cd(path(pkg))
        tic() # start timer
        du = split(readchomp(`du -sh $(path(pkg))`))[1] # src code disk use
        if name(pkg) in ["xerces-c","root","amptools","geant4","evio","rcdb","ccdb","jana","hdds","sim-recon","gluex_root_analysis","gluex_workshops"]
            if name(pkg) == "sim-recon" cd("src") end
            if uses_cmake(pkg)
                mk_cd("../$(name(pkg))-build")
                run(`mv ../$(version(pkg)) ../$(name(pkg))`)
                mkpath("../$(version(pkg))")
            end
            for cmd in cmds(pkg)
                run(`sh -c $cmd`)
            end
            if uses_cmake(pkg)
                cd("../"); run(`rm -rf $(name(pkg))-build $(name(pkg))`) end
        elseif name(pkg) == "cernlib"
            run(`cp -pr $home/patches $(path(pkg))`)
            run(`sh -c "patch < $(path(pkg))/patches/cernlib/Install_cernlib.patch"`)
            run(`./Install_cernlib`)
            run(`rm -rf 2005/build 2005/src`)
            run(`mv 2005 ../`); cd("../"); run(`rm -rf cernlib`)
            mk_cd(path(pkg)); run(`mv ../2005 .`)
        end # stop timer and write success file
        du_f = split(readchomp(`du -sh $(path(pkg))`))[1]
        success_file = open(path_to_success,"w")
        println(success_file,string("$(name(pkg))-$(git_version(pkg))","\n$(readchomp(`date "+%Y-%m-%d_%H:%M:%S"`))","\n# build time (seconds)\n",round(Int,toc()),
        "\n# disk use, final minus initial\n","\"$(du_f)B - $(du)B\"","\n# compiled against\n",tagged_deps(pkg)))
        close(success_file)
    elseif !is_external(pkg) && ispath(path_to_success)
        d = readdlm(path_to_success,use_mmap=false); w = 22
        if first_success == name(pkg)
            print(Base.text_colors[:bold]); hz("-")
            print(string(rpad("package",w," "),rpad("build time",w-6," "),rpad("disk use",w-3," "),"timestamp\n"))
            hz("-"); print(Base.text_colors[:normal])
        end
        println(rpad(d[1],w," "),rpad(string(d[3]," s"),w-6," "),rpad(d[4],w-3," "),d[2])
    end
end
