module Environs
#
export putenv,printenv
# set environment variables
using Packages
GLUEX_TOP = gettop()
home = Dict{ASCIIString,ASCIIString}()
vers = Dict{ASCIIString,ASCIIString}()
for pkg in get_packages()
    home[name(pkg)] = path(pkg)
    vers[name(pkg)] = version(pkg)
end
if int(vers["cernlib"]) != 2005 && int(vers["cernlib"]) != 2006 warn("using an old CERN_LEVEL (not 2005 or 2006)") end
BMS_OSNAME_BASE = osrelease()
if gettag() != "" 
    BMS_OSNAME = string(osrelease(),"_",gettag())
else
    BMS_OSNAME = osrelease()
end
ROOTSYS = joinpath(home["root"],vers["root"])
if split(home["root"],"/")[end] != "root" ROOTSYS = home["root"] end
CCDB_CONNECTION = "mysql://ccdb_user@hallddb.jlab.org/ccdb"
if haskey(ENV,"USER") USER = ENV["USER"]
else USER = readchomp(`whoami`) end
myenv = [
         "GLUEX_TOP" => "$GLUEX_TOP",
         "BMS_OSNAME" => "$BMS_OSNAME",
         "CERN" => home["cernlib"],
         "CERN_LEVEL" => vers["cernlib"],
         "ROOTSYS" => "$ROOTSYS",
         "CLHEP" => joinpath(home["clhep"],vers["clhep"]),
         "CLHEP_LIB_DIR" => joinpath(home["clhep"],vers["clhep"],"lib"),
         "CLHEP_INCLUDE_DIR" => joinpath(home["clhep"],vers["clhep"],"include"),
         "AMPTOOLS" => home["amptools"],
         "XERCESCROOT" => home["xerces-c"],
         "EVIOROOT" => string(home["evio"],"/",readchomp(`uname -s`),"-",readchomp(`uname -m`)),
         "CCDB_HOME" => home["ccdb"],
         "CCDB_CONNECTION" => "$CCDB_CONNECTION",
         "CCDB_USER" => "$USER",
         "HDDS_HOME" => home["hdds"],
         "JANA_HOME" => string(home["jana"],"/$BMS_OSNAME_BASE"),
         "JANA_CALIB_URL" => "$CCDB_CONNECTION",
         "JANA_GEOMETRY_URL" => string("xmlfile://",home["hdds"],"/main_HDDS.xml"),
         "HALLD_HOME" => home["sim-recon"],
         "MONIT_HOME" => string(home["online-monitoring"],"/$BMS_OSNAME")]
#         
myoptenv = ["JANA_CALIB_CONTEXT" => "\"variation=mc\""]
#
function putenv() 
    # put myenv variables into global ENV dictionary
    for (k,v) in myenv
        ENV[k] = v 
    end
    function add_to_path(path,new_path)
        if !contains(path,new_path) && !contains(new_path,"NA") 
            if path == ""
                return new_path 
            else
                return string(new_path,":",path) 
            end
        end
        path
    end
    # check that PATH variables exist, set to empty string if not
    if !haskey(ENV,"PATH") ENV["PATH"] = "" end
    if !haskey(ENV,"LD_LIBRARY_PATH") ENV["LD_LIBRARY_PATH"] = "" end
    if !haskey(ENV,"PYTHONPATH") ENV["PYTHONPATH"] = "" end
    if !haskey(ENV,"JANA_PLUGIN_PATH") ENV["JANA_PLUGIN_PATH"] = "" end
    # do PATH
    paths = [home["python"],joinpath(ENV["CERN"],ENV["CERN_LEVEL"]),ENV["ROOTSYS"],ENV["XERCESCROOT"],ENV["EVIOROOT"],ENV["CCDB_HOME"],ENV["HDDS_HOME"],ENV["JANA_HOME"],joinpath(ENV["HALLD_HOME"],ENV["BMS_OSNAME"]),ENV["CLHEP"]]
    for p in paths
        ENV["PATH"] = add_to_path(ENV["PATH"],string(p,"/bin"))
    end
    # do LD_LIBRARY_PATH
    for ldp in paths
        ENV["LD_LIBRARY_PATH"] = add_to_path(ENV["LD_LIBRARY_PATH"],string(ldp,"/lib"))
    end
    # do PYTHONPATH
    pypaths = [string(ENV["ROOTSYS"],"/lib"),string(joinpath(ENV["CCDB_HOME"],"python"),":",joinpath(ENV["CCDB_HOME"],"python","ccdb","ccdb_pyllapi/")),string(joinpath(ENV["HALLD_HOME"],ENV["BMS_OSNAME"]),"/lib/python")]
    for pyp in pypaths
        ENV["PYTHONPATH"] = add_to_path(ENV["PYTHONPATH"],pyp)
    end
    plugin_paths = [string(ENV["JANA_HOME"],"/plugins"),string(joinpath(ENV["HALLD_HOME"],ENV["BMS_OSNAME"]),"/plugins"),string(ENV["MONIT_HOME"],"/plugins")]
    # do JANA_PLUGIN_PATH
    for plugin_path in plugin_paths
        ENV["JANA_PLUGIN_PATH"] = add_to_path(ENV["JANA_PLUGIN_PATH"],plugin_path)
    end
    # remove items with Non-Applicable (NA) paths
    for (k,v) in ENV
        if contains(v,"NA")
            pop!(ENV,k)
        end
    end
    #
    ENV
end

function printenv()
    putenv() # set the environment variables before printing them to C-shell and bash scripts
    if !ispath("$GLUEX_TOP") mkdir("$GLUEX_TOP") end
    if !ispath("$GLUEX_TOP/scripts") mkdir("$GLUEX_TOP/scripts") end
    if !ispath("$GLUEX_TOP/scripts/env") mkdir("$GLUEX_TOP/scripts/env") end
    id = gettag()
    if id == ""
        file = open("$GLUEX_TOP/scripts/env/env_halld.csh","w")
    else
        file = open("$GLUEX_TOP/scripts/env/env_halld_$id.csh","w")
    end
    println(file,"#!/bin/tcsh\n#")
    for (k,v) in myenv
         if !contains(v,"NA") println(file,"setenv $k $v") end
    end
    for (k,v) in myoptenv
        println(file,"#setenv $k $v")
    end
    for path_name in ["PATH","LD_LIBRARY_PATH","PYTHONPATH","JANA_PLUGIN_PATH"]
        path = ENV[path_name]
        println(file,"\nsetenv $path_name $path")
    end
    close(file)
    if id == ""
        file = open("$GLUEX_TOP/scripts/env/env_halld.sh","w")
    else
        file = open("$GLUEX_TOP/scripts/env/env_halld_$id.sh","w")
    end
    println(file,"#!/bin/bash\n#")
    for (k,v) in myenv
         if !contains(v,"NA") println(file,"export $k=$v") end
    end
    for (k,v) in myoptenv
        println(file,"#export $k=$v")
    end
    for path_name in ["PATH","LD_LIBRARY_PATH","PYTHONPATH","JANA_PLUGIN_PATH"]
        path = ENV[path_name]
        println(file,"\nexport $path_name=$path")
    end
    close(file)
end

end
