module Packages
#
export Package,name,version,url,path,nthreads,tobuild,get_packages,get_package,gettop,osrelease,gettag,select_template,getbase,script,show_settings,get_unpack_file,mk_cd
#
immutable Package
    name::ASCIIString
    version::ASCIIString
    url::ASCIIString
    path::ASCIIString
    nthreads::ASCIIString
    tobuild::Bool
end
#
name(a::Package) = a.name
version(a::Package) = a.version
url(a::Package) = a.url
path(a::Package) = a.path
nthreads(a::Package) = a.nthreads
tobuild(a::Package) = a.tobuild
#
Package() = Package("","","","","",false)
#
function select_template(id)
    run(`rm -rf settings`)
    run(`cp -pr templates/settings_$id settings`)
    println(open("settings/id.txt","w"),id)
end
#
function mkbool(a::ASCIIString)
    if a != "true" && a != "false" error("tobuild() must be 'true' or 'false'; Please check for typos.") end
    if a == "true" return true
    else return false
    end
end
#
function getbase(path) 
    s = "/"
    a = split(path,"/")
    for i=1:length(a) - 1
        s = joinpath(s,a[i])
    end
    s
end
function mk_cd(path)
    if !ispath(path) mkdir(path) end 
    cd(path)
end
#
function gettop()
    top = string(pwd(),"/build_",readchomp(`date "+%Y-%m-%d"`))
    custom_top = readdlm("settings/top.txt",ASCIIString)
    if size(custom_top,1) != 1 || size(custom_top,2) != 2 error("problem reading in custom top directory name; top.txt has wrong number of rows or columns.") end
    if custom_top[1,1] != "default" 
        top = custom_top[1,1] 
        if !ispath(getbase(top)) error("base directory of custom top does not exist.") end
    end
    return top
end
#
osrelease() = readchomp(`osrelease.pl`)
#
function gettag()
    tag = ""
    custom_tag = readdlm("settings/top.txt",ASCIIString)
    if size(custom_tag,1) != 1 || size(custom_tag,2) != 2 error("problem reading in custom tag name; top.txt has wrong number of rows or columns.") end
    if custom_tag[1,2] != "default" tag = custom_tag[1,2] end
    return tag
end
#
getid() = readchomp("settings/id.txt")
#
function script(pkg)
    if name(pkg) == "online-sbms" || name(pkg) == "scripts"
        return true
    else
        return false
    end
end
#
function get_packages()
    vers = readdlm("settings/vers.txt",ASCIIString)
    urls = readdlm("settings/urls.txt",ASCIIString)
    paths = readdlm("settings/paths.txt",ASCIIString)
    nthreads = readdlm("settings/nthreads.txt",ASCIIString)
    tobuild = readdlm("settings/tobuild.txt",ASCIIString)
    pkgs = Array(Package,0)
    for i=1:size(paths,1)
        name = paths[i,1]
        if isabspath(paths[i,2])
            path = paths[i,2]
        else
            path = joinpath(gettop(),paths[i,2])
        end
        push!(pkgs,Package(name,vers[i,2],urls[i,2],path,nthreads[i,2],mkbool(tobuild[i,2])))
    end
    pkgs
end
function get_package(a::ASCIIString)
    for pkg in get_packages()
        if name(pkg) == a return pkg end
    end 
end
function get_unpack_file(URL)
    try
        run(`wget $URL`)
    catch
        run(`curl -O $URL`)
    end
    file = split(URL,"/")[end]
    run(`tar -xzvf $file`); rm(file)
end
function max_sizes()
    sizes = [:name=>0,:version=>0,:url=>0,:path=>0,:nthreads=>0,:tobuild=>0]
    for pkg in get_packages()
        for n in names(pkg)
            if length(pkg.(n)) > sizes[n] sizes[n] = length(pkg.(n)) end
        end
    end 
    return sizes
end
function show_settings()
    if !ispath("settings/") 
        error("no build settings to show. Please select a build settings template by running:\n'julia select_template.jl <id>'") 
    end
    print("\n",Base.text_colors[:bold])
    println("Current build settings",Base.text_colors[:bold])
    println("ID: ",getid())
    println("TOP: ",gettop())
    println("TAG: ",gettag())
    sizes = max_sizes()
    print("\n",Base.text_colors[:bold])
    for n in names(Package)
        print(n,Base.text_colors[:bold])
        spaces = sizes[n] - length(string(n)) + 8
        for i=1:spaces print(" ") end
    end
    print("\n",Base.text_colors[:normal])
    for pkg in get_packages()
        for n in names(pkg)
            print(pkg.(n),Base.text_colors[:normal])
            spaces = sizes[n] - length(pkg.(n)) + 8
            for i=1:spaces print(" ") end
        end
        print("\n")
    end 
end
#
end
