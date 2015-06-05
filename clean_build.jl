using Environs,Packages
# clean and then do fresh build
# set halld environment for scons command
setenv(`scons`,putenv())
printenv() # save env variables in shell script
BMS_OSNAME = ENV["BMS_OSNAME"]
# build packages
for pkg in get_packages()
    if tobuild(pkg) && !script(pkg)
        cd(path(pkg))
        if name(pkg) == "ccdb"
            run(`scons -c`); run(`scons`)    
        end
        if name(pkg) in ["jana","hdds","sim-recon","online-monitoring"]
            run(`rm -rf $BMS_OSNAME`)
            if ispath("src") cd("src") end; run(`rm -rf .$BMS_OSNAME`)
            if name(pkg) ==  "online-monitoring" cd("plugins") end
            #run(`scons -u -c install`)
            run(`scons -u -j$(nthreads(pkg)) install`)
        end
    end
end
