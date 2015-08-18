module Packages
# organize package information
home = pwd()
export Package,name,version,url,path,cmds,is_external,get_packages,get_package,gettop,osrelease,gettag,select_template,show_settings
export get_unpack_file,mk_cd,get_template_ids,get_pkg_names,get_deps,tagged_deps,git_version,check_deps,mk_template,install_dirname
immutable Package
    name::ASCIIString
    version::ASCIIString
    url::ASCIIString
    path::ASCIIString
    cmds::Array{ASCIIString,1}
    deps::ASCIIString
end
name(a::Package) = a.name
version(a::Package) = a.version
url(a::Package) = a.url
path(a::Package) = a.path
cmds(a::Package) = a.cmds
deps(a::Package) = a.deps
is_external(a::Package) = length(cmds(a)) == 0
#
function select_template(id)
    if !ispath("templates") run(`cp -pr example-templates templates`) end
    run(`rm -rf settings`)
    run(`cp -pr templates/settings-$id settings`)
    println(open("settings/id.txt","w"),id)
end
#
function get_template_ids()
    if !ispath("templates") run(`cp -pr example-templates templates`) end
    list = Array(ASCIIString,0)
    for dir in readdir("templates")
        push!(list,split(dir,"-")[2])
    end
    list
end
function mk_template(id)
    if ispath("templates/settings-$id") warn("renaming template with same id as old-$id");run(`mv templates/settings-$id templates/settings-old-$id`) end
    run(`cp -pr settings templates/settings-$id`)
    run(`rm templates/settings-$id/id.txt`)
end
#
function mk_cd(path)
    mkpath(path); cd(path)
end
#
function gettop()
    top = string(pwd(),"/pkgs")
    custom_top = readdlm("settings/top.txt",ASCIIString)
    if size(custom_top,1) != 1 || size(custom_top,2) != 2 error("problem reading in custom top directory name; top.txt has wrong number of rows or columns.") end
    if custom_top[1,1] != "default"
        top = custom_top[1,1]
        if !isabspath(top) top = string(pwd(),"/pkgs/",top) end
    end
    top
end
#
osrelease() = readchomp(`perl src/osrelease.pl`)
#
function gettag()
    tag = ""
    custom_tag = readdlm("settings/top.txt",ASCIIString)
    if size(custom_tag,1) != 1 || size(custom_tag,2) != 2 error("problem reading in custom tag name; top.txt has wrong number of rows or columns.") end
    if custom_tag[1,2] != "default" tag = custom_tag[1,2] end
    tag
end
install_dirname() = (gettag() == "") ? osrelease() : string("build-",gettag())
getid() = readchomp("settings/id.txt")
#
get_pkg_names() = ["xerces-c","cernlib","root","amptools","geant4","evio","ccdb","jana","hdds","sim-recon"]
#
function get_packages()
    vers = readdlm("settings/vers.txt",ASCIIString)
    urls = readdlm("settings/urls.txt",ASCIIString)
    paths = readdlm("settings/paths.txt",ASCIIString)
    pkg_names = get_pkg_names()
    @assert(vers[:,1] == pkg_names,string("'vers.txt' has wrong number of packages, names, or order.\nNeed to match ",pkg_names,"\n"))
    @assert(urls[:,1] == pkg_names,string("'urls.txt' has wrong number of packages, names, or order.\nNeed to match ",pkg_names,"\n"))
    @assert(paths[:,1] == pkg_names,string("'paths.txt' has wrong number of packages, names, or order.\nNeed to match ",pkg_names,"\n"))
    #
    commands = readdlm("settings/commands.txt",ASCIIString)
    tmp_cmds = Dict{ASCIIString,Array{ASCIIString,1}}()
    cmds = Dict{ASCIIString,Array{ASCIIString,1}}()
    for name in get_pkg_names()
        tmp_cmds[name] = Array(ASCIIString,0)
        cmds[name] = Array(ASCIIString,0)
    end
    for i=1:size(commands,1)
        push!(tmp_cmds[commands[i,1]],commands[i,2])
    end
    mydeps = [
        "xerces-c" => "",
        "cernlib" => "",
        "root" => "",
        "amptools" => "root",
        "geant4" => "",
        "evio" => "",
        "ccdb" => "",
        "jana" => "xerces-c,root,ccdb",
        "hdds" => "xerces-c",
        "sim-recon" => "xerces-c,cernlib,root,evio,ccdb,jana,hdds"]
    @osx_only mydeps["sim-recon"] = "xerces-c,root,evio,ccdb,jana,hdds"
    pkgs = Array(Package,0)
    for i=1:size(paths,1)
        name = paths[i,1]
        path = paths[i,2]; path = replace(path,"[OS]",osrelease()); path = replace(path,"[VER]",vers[i,2])
        url = replace(urls[i,2],"[VER]",vers[i,2])
        if !isabspath(path) && path != "NA"
            path = joinpath(gettop(),path)
        end
        core = ["xerces-c","root","evio","ccdb","jana","hdds","sim-recon"]
        if path == "NA" && name in core
            error("core packages cannot be disabled. Please replace 'NA' with a valid path in 'paths.txt'.
            core: ",core,"\n") end
        for cmd in tmp_cmds[name]; if path == "NA" continue end
            push!(cmds[name],replace(cmd,"[PATH]",path))
        end
        push!(pkgs,Package(name,vers[i,2],url,path,cmds[name],mydeps[name]))
    end
    pkgs
end
function get_package(a::ASCIIString)
    cd(home)
    for pkg in get_packages()
        if name(pkg) == a return pkg end
    end
end # use git hash for git-repo. packages
git_version(a) = ispath(joinpath(path(a),".git")) ? begin dir = pwd(); cd(path(a)); ver = readchomp(`git log -1 --format="%h"`); cd(dir); ver end : version(a)
function get_deps(arguments)
    mydeps = Array(ASCIIString,0)
    for pkg_name in arguments
        pkg_name = convert(ASCIIString,pkg_name)
        for dep in split(deps(get_package(pkg_name)),",")
            dep = convert(ASCIIString,dep)
            if dep != ""  push!(mydeps,dep) end
        end
    end
    unique(mydeps)
end
function tagged_deps(a)
    mydeps = Array(ASCIIString,0)
    for dep in split(deps(a),",")
        dep = convert(ASCIIString,dep)
        if dep == "" continue end
        push!(mydeps,string(dep,"-",git_version(get_package(dep))))
    end
    if length(mydeps) == 0 push!(mydeps,"none listed") end
    string("\"",join(mydeps,","),"\"")
end
function get_unpack_file(URL,PATH="")
    file = basename(URL); info("downloading $file")
    run(`curl -OL $URL`)
    if PATH != ""
        mkpath(PATH); if readchomp(`tar tf $file`|>`head`)[1] != '.' ncomp = 1 else ncomp = 2 end
        run(`tar xf $file -C $PATH --strip-components=$ncomp`)
    else
        run(`tar xf $file`)
    end
    rm(file)
end
function show_settings(;col=:all,sep=2)
    if !ispath("settings/")
        error("no build settings to show. Please select a build settings template by running:\n\t'hdpm select <id>'")
    end
    if !(col in names(Package)) && col != :all
        error("incorrect name: use one of the following ",[string(i) for i in names(Package)])
    end
    print("\n",Base.text_colors[:bold])
    println("Current build settings",Base.text_colors[:bold])
    try
        println("ID: ",getid())
    catch
        println("ID: ","id file not found; This will not affect build.")
    end
    println("TOP: ",gettop())
    println("TAG: ",gettag())
    #
    w1 = 9 + sep; w2 = 87+sep
    print("\n",Base.text_colors[:bold])
    for n in names(Package)
        w = w1
	if (n == :deps || n == :cmds) && col == :all continue end
        if n in [:name,col] || col == :all
            if col == :all && n == :url w = w2 end
            print(rpad(n,w," "),Base.text_colors[:bold])
        end
    end
    print("\n",Base.text_colors[:normal])
    for pkg in get_packages()
        p = replace(path(pkg),string(gettop(),"/"),"")
        if col==:all
            println(rpad(name(pkg),w1," "),rpad(git_version(pkg),w1," "),rpad(url(pkg),w2," "),p)
        elseif col==:version
            println(rpad(name(pkg),w1," "),git_version(pkg))
        elseif col==:url
            println(rpad(name(pkg),w1," "),url(pkg))
        elseif col==:path
            println(rpad(name(pkg),w1," "),p)
        elseif col==:deps
            println(rpad(name(pkg),w1," "),deps(pkg))
        elseif col==:cmds
            for cmd in cmds(pkg)
                println(rpad(name(pkg),w1," "),replace(cmd,string(gettop(),"/"),""))
            end
        end
    end
end
function check_deps(pkg)
    @linux_only begin LDD = `ldd`; OE = `so` end
    @osx_only begin LDD = `otool -L`; OE = `dylib` end
    test_cmds = [
        "xerces-c" => `$LDD $(path(get_package("xerces-c")))/lib/libxerces-c.$OE` |> `grep libcurl`,
        "cernlib" => `ls -lh $(path(get_package("cernlib")))/$(version(get_package("cernlib")))/lib/libgeant321.a`,
        "root" => `root -q -l`,
        "evio" => `evio2xml`,
        "ccdb" => `ccdb`,
        "jana" => `jana`,
        "hdds" => `$LDD $(path(get_package("hdds")))/$(install_dirname())/lib/libhdds.so` |> `grep libxerces-c`,
        "sim-recon" => `hd_root`]
    for dep in get_deps([name(pkg)])
        if !success(test_cmds[dep])
            error("'$dep' does not appear to be installed. Please check paths if using external installations.
            To build all dependencies, run 'hdpm build' with all packages enabled in 'commands.txt'.\n")
        end
    end
end
#
end
