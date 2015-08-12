using Environs,Packages
home = pwd()
printenv() # save env variables in shell script
# build packages
for pkg in get_packages()
    path_to_success = joinpath(path(pkg),"success.hdpm")
    if tobuild(pkg) && !ispath(path_to_success)
        @osx_only if name(pkg) == "cernlib" info("Mac OS X detected: skipping cernlib");continue end
        if !ispath(path(pkg)) println();error("path does not exist, ",path,"; Run 'hdpm co' first.") end
        cd(path(pkg)); println(); info("building $(name(pkg))")
        tic() # start timer
        if name(pkg) == "cernlib"
            run(`cp -pr $home/patches $(path(pkg))`)
            run(`sh -c "patch < $(path(pkg))/patches/cernlib/Install_cernlib.patch"`)
            run(`./Install_cernlib`)
        end
        if name(pkg) == "root"
            run(setenv(`sh -c "./configure --enable-roofit"`,putenv()))
            run(setenv(`make -j $(nthreads(pkg))`,putenv()))
        end
        if name(pkg) == "xerces-c"
            run(`sh -c "./configure --prefix=$(path(pkg))"`)
            run(`make`); run(`make install`)
        end
        if name(pkg) == "geant4"
            mk_cd("../$(name(pkg))_build")
            run(`cmake -DCMAKE_INSTALL_PREFIX=$(path(pkg)) $(path(pkg))`)
            run(`make -j $(nthreads(pkg))`)
            run(`make install`)
            cd("../"); run(`rm -rf $(name(pkg))_build`)
        end
        if name(pkg) == "amptools" run(setenv(`make`,putenv())) end
        @osx_only if name(pkg) == "ccdb" run(setenv(`scons with-mysql=false`,putenv())) end
        @linux_only if name(pkg) == "ccdb" run(setenv(`scons`,putenv())) end
        if name(pkg) == "evio" run(`scons --prefix=$(path(pkg)) install`) end
        if name(pkg) in ["jana","hdds","sim-recon"]
            if ispath("src") cd("src") end
            if name(pkg) == "sim-recon"
                run(setenv(`scons -u -j$(nthreads(pkg)) install DEBUG=0`,putenv()))
            else
                run(setenv(`scons -u -j$(nthreads(pkg)) install`,putenv()))
            end
        end # stop timer and write success file
        cd(path(pkg)); success_file = open("success.hdpm","w")
        println(success_file,string("# build time (seconds)\n",round(toc(),1),
        "\n# disk use (Bytes), including src code\n",split(readchomp(`du -sh $(path(pkg))`))[1]))
        close(success_file)
    elseif ispath(path_to_success)
        d = readdlm(path_to_success)
        info(string(name(pkg),": compile time = ",d[1]," seconds, disk usage = ",d[2],"B"))
    end
end
