using Environs,Packages
home = pwd()
printenv() # save env variables in shell script
# build packages
for pkg in get_packages()
    if tobuild(pkg) && !script(pkg)
        if !ispath(path(pkg)) throw("Error: Path does not exist, ",path) end
        cd(path(pkg))
        if name(pkg) == "cernlib"
            run(`cp -pr $home/patches $(path(pkg))`)
            run(`sh -c "patch < $(path(pkg))/patches/cernlib/Install_cernlib.patch"`)
            run(`./Install_cernlib`)
        end
        if name(pkg) == "root"
            PREFIX=joinpath(path(pkg),version(pkg))
            if ispath(PREFIX) run(`rm -rf $PREFIX`) end
            MYENV = putenv(); delete!(MYENV,"ROOTSYS")
            run(setenv(`sh -c "configure --prefix=$PREFIX --libdir=$PREFIX/lib --incdir=$PREFIX/include --etcdir=$PREFIX/etc --enable-roofit"`,MYENV))
            run(setenv(`make -j $(nthreads(pkg))`,MYENV))
            run(`make install`)
            run(`cp -p config.log $PREFIX`)
            run(`cp -p config.status $PREFIX`)
            run(`make maintainer-clean`)
        end
        if name(pkg) == "xerces-c"
            run(`sh -c "configure --prefix=$(path(pkg))"`)
            run(`make`); run(`make install`)
        end
        if name(pkg) == "clhep"
            mk_cd("../clhep_build")
            run(`cmake -DCMAKE_INSTALL_PREFIX=$(joinpath(path(pkg),version(pkg))) $(path(pkg))`)
            run(`make -j $(nthreads(pkg))`); run(`make test`); run(`make install`)
            cd("../"); run(`rm -rf clhep_build`)
        end
        if name(pkg) == "amptools" run(setenv(`make`,putenv())) end
        if name(pkg) == "ccdb" run(setenv(`scons`,putenv())) end
        if name(pkg) == "evio"  run(`scons --prefix=$(path(pkg)) install`) end
        if name(pkg) in ["jana","hdds","sim-recon","online-monitoring"]
            if ispath("src") cd("src") end
            if name(pkg) ==  "online-monitoring" cd("plugins") end
            run(setenv(`scons -u -j$(nthreads(pkg)) install`,putenv()))
        end
    end
end
