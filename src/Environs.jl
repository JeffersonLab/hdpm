module Environs
# set environment variables
export putenv,printenv
using Packages
GLUEX_TOP = gettop()
home = Dict{ASCIIString,ASCIIString}()
vers = Dict{ASCIIString,ASCIIString}()
for pkg in get_packages()
    home[name(pkg)] = path(pkg)
    vers[name(pkg)] = version(pkg)
end
if int(vers["cernlib"]) != 2005 && int(vers["cernlib"]) != 2006 println();warn("using an old CERN_LEVEL (not 2005 or 2006)") end
BMS_OSNAME = install_dirname()
JANA_HOME = is_external(get_package("jana")) ? string(home["jana"],"/$(osrelease())") : string(home["jana"],"/$BMS_OSNAME")
@osx_only begin home["cernlib"] = "NA";vers["cernlib"] = "NA";println();info("Mac OS X detected: disabling cernlib") end
home_python = "NA"
if contains(osrelease(),"CentOS6") && ispath("/apps/python/PRO")
    home_python = "/apps/python/PRO"
end
CCDB_CONNECTION = "mysql://ccdb_user@hallddb.jlab.org/ccdb"
if haskey(ENV,"USER") USER = ENV["USER"]
else USER = readchomp(`whoami`) end
myenv = [
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
myoptenv = ["JANA_CALIB_CONTEXT" => "\"variation=mc\""]
#
function putenv()
    # put myenv variables into global ENV dictionary
    for (k,v) in myenv
        ENV[k] = v
    end
    function add_to_path(path,new_path)
        if !contains(path,new_path) && !contains(new_path,"/NA/")
            return (path == "") ? new_path : string(new_path,":",path)
        end
        path
    end
    # check that PATH variables exist, set to empty string if not
    if !haskey(ENV,"PATH") ENV["PATH"] = "" end
    @linux_only if !haskey(ENV,"LD_LIBRARY_PATH") ENV["LD_LIBRARY_PATH"] = "" end
    @osx_only if !haskey(ENV,"DYLD_LIBRARY_PATH") ENV["DYLD_LIBRARY_PATH"] = "" end
    if !haskey(ENV,"PYTHONPATH") ENV["PYTHONPATH"] = "" end
    if !haskey(ENV,"JANA_PLUGIN_PATH") ENV["JANA_PLUGIN_PATH"] = "" end
    # do PATH
    paths = [joinpath(ENV["CERN"],ENV["CERN_LEVEL"]),ENV["ROOTSYS"],ENV["XERCESCROOT"],ENV["EVIOROOT"],ENV["CCDB_HOME"],ENV["HDDS_HOME"],ENV["JANA_HOME"],joinpath(ENV["HALLD_HOME"],ENV["BMS_OSNAME"])]
    if home_python != "NA" push!(paths,home_python) end
    for p in paths
        ENV["PATH"] = add_to_path(ENV["PATH"],string(p,"/bin"))
    end
    # do LD_LIBRARY_PATH
    for ldp in paths
        @linux_only ENV["LD_LIBRARY_PATH"] = add_to_path(ENV["LD_LIBRARY_PATH"],string(ldp,"/lib"))
        @osx_only ENV["DYLD_LIBRARY_PATH"] = add_to_path(ENV["DYLD_LIBRARY_PATH"],string(ldp,"/lib"))
    end
    # do PYTHONPATH
    pypaths = [string(ENV["ROOTSYS"],"/lib"),string(joinpath(ENV["CCDB_HOME"],"python"),":",joinpath(ENV["CCDB_HOME"],"python","ccdb","ccdb_pyllapi/")),string(joinpath(ENV["HALLD_HOME"],ENV["BMS_OSNAME"]),"/lib/python")]
    for pyp in pypaths
        ENV["PYTHONPATH"] = add_to_path(ENV["PYTHONPATH"],pyp)
    end
    plugin_paths = [string(ENV["JANA_HOME"],"/plugins"),string(joinpath(ENV["HALLD_HOME"],ENV["BMS_OSNAME"]),"/plugins")]
    # do JANA_PLUGIN_PATH
    for plugin_path in plugin_paths
        ENV["JANA_PLUGIN_PATH"] = add_to_path(ENV["JANA_PLUGIN_PATH"],plugin_path)
    end
    # remove items with Non-Applicable (NA) paths
    for (k,v) in ENV
        if v == "NA" pop!(ENV,k) end
    end
    ENV
end
function printenv()
    putenv() # set env variables before printing them for tcsh and bash shells
    mkpath("$GLUEX_TOP/env-setup")
    id = gettag()
    function myprint(sh)
        if sh == "sh" n = "bash"; set = "export"; eq = "="
        elseif sh == "csh" n = "tcsh"; set = "setenv"; eq = " "
        else error("unknown shell type") end
        file = (id == "") ? open("$GLUEX_TOP/env-setup/env_halld.$sh","w") : open("$GLUEX_TOP/env-setup/env_halld_$id.$sh","w")
        println(file,"# $n\n#")
        println(file,string("$set GLUEX_TOP$eq",ENV["GLUEX_TOP"]))
        println(file,string("$set BMS_OSNAME$eq",ENV["BMS_OSNAME"]))
        for (k,v) in myenv; if v == "NA" || k == "GLUEX_TOP" || k == "BMS_OSNAME" continue end
            v = replace(v,ENV["GLUEX_TOP"],"\$GLUEX_TOP");v = replace(v,ENV["BMS_OSNAME"],"\$BMS_OSNAME")
            println(file,"$set $k$(eq)$v")
        end
        for (k,v) in myoptenv
            println(file,"#$set $k$(eq)$v")
        end
        @linux_only ldlp = "LD_LIBRARY_PATH"
        @osx_only ldlp = "DYLD_LIBRARY_PATH"
        for path_name in ["PATH",ldlp,"PYTHONPATH","JANA_PLUGIN_PATH"]
            path = ENV[path_name]
            for (k,v) in myenv; if v == "NA" || k == "GLUEX_TOP" || k == "CCDB_USER" continue end
                path = replace(path,v,string("\$",k))
            end
            path = replace(path,ENV["GLUEX_TOP"],"\$GLUEX_TOP")
            println(file,"\n$set $path_name$(eq)$path")
        end
        close(file)
    end
    myprint("sh"); myprint("csh")
end

end
