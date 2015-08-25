module Environs
# set environment variables
export getenv,printenv
using Packages
function getenv()
    GLUEX_TOP = gettop()
    home = Dict{ASCIIString,ASCIIString}()
    vers = Dict{ASCIIString,ASCIIString}()
    for pkg in get_packages()
        home[name(pkg)] = path(pkg)
        vers[name(pkg)] = version(pkg)
    end
    if int(vers["cernlib"]) != 2005 && int(vers["cernlib"]) != 2006 println();warn("using an old CERN_LEVEL (not 2005 or 2006)") end
    BMS_OSNAME = install_dirname()
    if is_external(get_package("hdds")) || is_external(get_package("sim-recon")) BMS_OSNAME = osrelease() end
    JANA_HOME = is_external(get_package("jana")) ? string(home["jana"],"/$(osrelease())") : string(home["jana"],"/$BMS_OSNAME")
    @osx_only begin home["cernlib"] = "NA";vers["cernlib"] = "NA";println();info("Mac OS X detected: disabling cernlib") end
    home_python = "NA"
    if contains(osrelease(),"CentOS6") && ispath("/apps/python/PRO")
        home_python = "/apps/python/PRO"
    end
    CCDB_CONNECTION = "mysql://ccdb_user@hallddb.jlab.org/ccdb"
    if haskey(ENV,"USER") USER = ENV["USER"]
    else USER = readchomp(`whoami`) end
    env = [
             "GLUEX_TOP" => "$GLUEX_TOP",
             "BMS_OSNAME" => "$BMS_OSNAME",
             "CERN" => home["cernlib"],
             "CERN_LEVEL" => vers["cernlib"],
             "ROOTSYS" => home["root"],
             "AMPTOOLS" => home["amptools"],
             "XERCESCROOT" => home["xerces-c"],
             "EVIOROOT" => string(home["evio"],"/",readchomp(`uname -s`),"-",readchomp(`uname -m`)),
             "CCDB_HOME" => home["ccdb"],
             "CCDB_CONNECTION" => "$CCDB_CONNECTION",
             "CCDB_USER" => "$USER",
             "HDDS_HOME" => home["hdds"],
             "JANA_HOME" => "$JANA_HOME",
             "JANA_CALIB_URL" => "$CCDB_CONNECTION",
             "JANA_GEOMETRY_URL" => string("xmlfile://",home["hdds"],"/main_HDDS.xml"),
             "HALLD_HOME" => home["sim-recon"]]
    #
    function add_to_path(path,new_path)
        if !contains(path,new_path) && !contains(new_path,"/NA/")
            return (path == "") ? new_path : string(new_path,":",path)
        end
        path
    end
    # check if PATH variables exist in ENV, set to empty string if not
    if !haskey(ENV,"PATH") env["PATH"] = ""
    else env["PATH"] = ENV["PATH"] end
    @linux_only if !haskey(ENV,"LD_LIBRARY_PATH") env["LD_LIBRARY_PATH"] = ""
    else env["LD_LIBRARY_PATH"] = ENV["LD_LIBRARY_PATH"] end
    @osx_only if !haskey(ENV,"DYLD_LIBRARY_PATH") env["DYLD_LIBRARY_PATH"] = ""
    else env["DYLD_LIBRARY_PATH"] = ENV["DYLD_LIBRARY_PATH"] end
    if !haskey(ENV,"PYTHONPATH") env["PYTHONPATH"] = ""
    else env["PYTHONPATH"] = ENV["PYTHONPATH"] end
    if !haskey(ENV,"JANA_PLUGIN_PATH") env["JANA_PLUGIN_PATH"] = ""
    else env["JANA_PLUGIN_PATH"] = ENV["JANA_PLUGIN_PATH"] end
    # do PATH
    paths = [joinpath(env["CERN"],env["CERN_LEVEL"]),env["ROOTSYS"],env["XERCESCROOT"],env["EVIOROOT"],env["CCDB_HOME"],env["HDDS_HOME"],env["JANA_HOME"],joinpath(env["HALLD_HOME"],env["BMS_OSNAME"])]
    if home_python != "NA" push!(paths,home_python) end
    for p in paths
        env["PATH"] = add_to_path(env["PATH"],string(p,"/bin"))
    end
    # do LD_LIBRARY_PATH
    for ldp in paths
        @linux_only env["LD_LIBRARY_PATH"] = add_to_path(env["LD_LIBRARY_PATH"],string(ldp,"/lib"))
        @osx_only env["DYLD_LIBRARY_PATH"] = add_to_path(env["DYLD_LIBRARY_PATH"],string(ldp,"/lib"))
    end
    # do PYTHONPATH
    pypaths = [string(env["ROOTSYS"],"/lib"),string(joinpath(env["CCDB_HOME"],"python"),":",joinpath(env["CCDB_HOME"],"python","ccdb","ccdb_pyllapi/")),string(joinpath(env["HALLD_HOME"],env["BMS_OSNAME"]),"/lib/python")]
    for pyp in pypaths
        env["PYTHONPATH"] = add_to_path(env["PYTHONPATH"],pyp)
    end
    plugin_paths = [string(env["JANA_HOME"],"/plugins"),string(joinpath(env["HALLD_HOME"],env["BMS_OSNAME"]),"/plugins")]
    # do JANA_PLUGIN_PATH
    for plugin_path in plugin_paths
        env["JANA_PLUGIN_PATH"] = add_to_path(env["JANA_PLUGIN_PATH"],plugin_path)
    end
    # remove items with Non-Applicable (NA) paths
    for (k,v) in env
        if v == "NA" pop!(env,k) end
    end
    env
end
function printenv()
    function myprint(sh,env) # print env-setup scripts for tcsh and bash shells
        myoptenv = ["JANA_CALIB_CONTEXT" => "\"variation=mc\""]
        mkpath("$(env["GLUEX_TOP"])/env-setup")
        id = gettag()
        if sh == "sh" n = "bash"; set = "export"; eq = "="
        elseif sh == "csh" n = "tcsh"; set = "setenv"; eq = " "
        else error("unknown shell type") end
        file = (id == "") ? open("$(env["GLUEX_TOP"])/env-setup/hdenv.$sh","w") : open("$(env["GLUEX_TOP"])/env-setup/hdenv-$id.$sh","w")
        println(file,"# $n\n#")
        println(file,string("$set GLUEX_TOP$eq",env["GLUEX_TOP"]))
        println(file,string("$set BMS_OSNAME$eq",env["BMS_OSNAME"]))
        for (k,v) in env; if k == "GLUEX_TOP" || k == "BMS_OSNAME" || contains(k,"PATH") continue end
            v = replace(v,env["GLUEX_TOP"],"\$GLUEX_TOP");v = replace(v,env["BMS_OSNAME"],"\$BMS_OSNAME")
            println(file,"$set $k$(eq)$v")
        end
        for (k,v) in myoptenv
            println(file,"#$set $k$(eq)$v")
        end
        @linux_only ldlp = "LD_LIBRARY_PATH"
        @osx_only ldlp = "DYLD_LIBRARY_PATH"
        for path_name in ["PATH",ldlp,"PYTHONPATH","JANA_PLUGIN_PATH"]
            path = env[path_name]
            for (k,v) in env; if k == "GLUEX_TOP" || k == "CCDB_USER" || contains(k,"PATH") continue end
                path = replace(path,v,string("\$",k))
            end
            path = replace(path,env["GLUEX_TOP"],"\$GLUEX_TOP")
            println(file,"\n$set $path_name$(eq)$path")
        end
        close(file)
    end
    const env = getenv()
    myprint("sh",env); myprint("csh",env)
    for (k,v) in env # expose env variables through global ENV
        ENV[k] = v
    end
end

end
