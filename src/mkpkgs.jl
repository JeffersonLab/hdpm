using Environs,Packages
home = pwd()
printenv() # save env variables in shell script
# build packages
for pkg in get_packages()
    if tobuild(pkg)
        if !ispath(path(pkg)) error("path does not exist, ",path) end
        cd(path(pkg)); println(); info("building $(name(pkg))")
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
        if name(pkg) == "ccdb" run(setenv(`scons`,putenv())) end
        if name(pkg) == "evio"  run(`scons --prefix=$(path(pkg)) install`) end
        if name(pkg) in ["jana","hdds","sim-recon"]
            if ispath("src") cd("src") end
            run(setenv(`scons -u -j$(nthreads(pkg)) install`,putenv()))
        end
    end
end
