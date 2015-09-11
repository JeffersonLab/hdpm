module Packages
# organize package information
home = pwd()
export Package,name,version,url,path,cmds,is_external,get_packages,get_package,gettop,osrelease,gettag,select_template,show_settings
export get_unpack_file,mk_cd,get_template_ids,get_pkg_names,get_deps,tagged_deps,git_version,check_deps,mk_template,install_dirname
export versions_from_xml
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
function write_id(id)
    fid = open("settings/id.txt","w"); println(fid,id); close(fid)
end
#
function select_template(id="master")
    run(`rm -rf settings`)
    if id == "master" run(`cp -pr templates/$id settings`)
    else run(`cp -pr templates/settings-$id settings`) end
    write_id(id)
end
#
function get_template_ids()
    if !ispath("settings") run(`cp -pr templates/master settings`)
        write_id("master") end
    list = Array(ASCIIString,0)
    push!(list,"master")
    for dir in readdir("templates")
        if contains(dir,"settings") push!(list,split(dir,"settings-")[2]) end
    end
    list
end
function disable_cmds()
    run(`mv settings/commands.txt settings/commands.txt.old`)
    file = ["cmds-old"=>open("settings/commands.txt.old"),"cmds"=>open("settings/commands.txt","w")]
    for line in readlines(file["cmds-old"])
        println(file["cmds"],string("#",chomp(line)))
    end
    for (k,v) in file close(v) end
    rm("settings/commands.txt.old")
end
function mk_template(id)
    if id == "master" error("not able to save template named 'master'. This id is reserved.\n") end
    if ispath("templates/settings-$id") warn("renaming template with same id as old-$id");run(`mv templates/settings-$id templates/settings-old-$id`) end
    if id == "jlab" info("saving 'jlab' template: All build commands are disabled.") end
    if id == "jlab" disable_cmds() end
    write_settings()
    run(`rm -f settings/*.txt~`)
    run(`cp -pr settings templates/settings-$id`)
    write_id(id)
end
#
function mk_cd(path)
    mkpath(path); cd(path)
end
#
function check_for_settings()
    if !ispath("settings")
        error("Please select a 'build template'.
        \t Use 'hdpm select <id>'
        \t ids: ",get_template_ids(),"\n") end
end
function gettop()
    check_for_settings()
    top = string(pwd(),"/pkgs")
    custom_top = readdlm("settings/top.txt",ASCIIString,use_mmap=false)
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
    custom_tag = readdlm("settings/top.txt",ASCIIString,use_mmap=false)
    if size(custom_tag,1) != 1 || size(custom_tag,2) != 2 error("problem reading in custom tag name; top.txt has wrong number of rows or columns.") end
    if custom_tag[1,2] != "default" tag = custom_tag[1,2] end
    tag
end
install_dirname() = (gettag() == "") ? osrelease() : string("build-",gettag())
get_pkg_names() = ["xerces-c","cernlib","root","amptools","geant4","evio","ccdb","jana","hdds","sim-recon"]
jlab_top() = string("/group/halld/Software/builds/",osrelease())
#
function major_minor(ver)
    for v in split(ver,"-")
        if contains(v,".") return split(v,".")[1],split(v,".")[2] end
    end
    "0","0"
end
function get_packages()
    check_for_settings()
    vers = readdlm("settings/versions.txt",ASCIIString,use_mmap=false)
    urls = readdlm("settings/urls.txt",ASCIIString,use_mmap=false)
    paths = readdlm("settings/paths.txt",ASCIIString,use_mmap=false)
    pkg_names = get_pkg_names()
    @assert(vers[:,1] == pkg_names,string("'versions.txt' has wrong number of packages, names, or order.\nNeed to match ",pkg_names,"\n"))
    @assert(urls[:,1] == pkg_names,string("'urls.txt' has wrong number of packages, names, or order.\nNeed to match ",pkg_names,"\n"))
    @assert(paths[:,1] == pkg_names,string("'paths.txt' has wrong number of packages, names, or order.\nNeed to match ",pkg_names,"\n"))
    #
    commands = [[] []]
    try
        commands = readdlm("settings/commands.txt",ASCIIString,use_mmap=false)
    catch
        info("No packages to build. Using external installations.")
    end
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
    jsep = ["xerces-c"=>"-","cernlib"=>"","root"=>"_","amptools"=>"_","geant4"=>"-","evio"=>"-","ccdb"=>"_","jana"=>"_","hdds"=>"-","sim-recon"=>"-"]
    pkgs = Array(Package,0)
    for i=1:size(paths,1)
        name = paths[i,1]
        path = paths[i,2]; path = replace(path,"[OS]",osrelease())
        path = (vers[i,2] != "latest") ? replace(path,"[VER]",vers[i,2]) : replace(replace(path,"-[VER]",""),"_[VER]","")
        url = urls[i,2]
        if name == "evio"
            evio_major_minor = join(major_minor(vers[i,2]),".")
            if !contains(url,evio_major_minor) url = replace(url,r"4.[0-9]",evio_major_minor) end
        end
        url = replace(url,"[VER]",vers[i,2])
        if !isabspath(path) && path != "NA"
            path = joinpath(gettop(),path)
        end
        if vers[i,2] == "NA" url = "NA"; path = "NA" end
        core = ["xerces-c","root","evio","ccdb","jana","hdds","sim-recon"]
        if path == "NA" && name in core
            error("core packages cannot be disabled. Please replace 'NA' with a valid path in 'paths.txt'.
            core: ",core,"\n") end
        for cmd in tmp_cmds[name]; if path == "NA" continue end
            push!(cmds[name],replace(cmd,"[PATH]",path))
        end
        if name == "ccdb" && length(cmds[name]) == 0 && ispath(jlab_top()) vers[i,2] = join(major_minor(vers[i,2]),".") end
        jpath = joinpath(jlab_top(),name,string(name,jsep[name],vers[i,2]))
        if length(cmds[name]) == 0 && ispath(jpath) path = jpath end
        if name == "cernlib" && length(cmds[name]) == 0 && ispath(joinpath(jlab_top(),name)) path = joinpath(jlab_top(),name) end
        if length(cmds[name]) > 0 && !contains(path,gettop()) path = joinpath(gettop(),basename(path)) end
        if (name == "hdds" || name == "sim-recon") && length(cmds[name]) > 0 && vers[i,2] != "latest"
            vmm = major_minor(vers[i,2])
            url_alt = "https://github.com/JeffersonLab/$name/archive/$name-$(vers[i,2]).tar.gz"
            if name == "hdds"
                if int(vmm[1]) <= 3 && int(vmm[2]) <= 2 || int(vmm[1]) <= 2 url = url_alt end
            elseif name == "sim-recon"
                if int(vmm[1]) <= 1 && int(vmm[2]) <= 3 || int(vmm[1]) == 0 || contains(vers[i,2],"dc") url = url_alt end
            end
        end
        if length(cmds[name]) > 0 && vers[i,2] == "latest" && contains(url,"https://github.com/JeffersonLab/$name/archive/")
            url = "https://github.com/JeffersonLab/$name" end
        if name == "jana" && length(cmds[name]) > 0 && vers[i,2] == "latest" url = "https://phys12svn.jlab.org/repos/JANA" end
        push!(pkgs,Package(name,vers[i,2],url,path,cmds[name],mydeps[name]))
    end
    pkgs
end
function write_settings()
    mkdir("settings-tmp")
    run(`cp -p settings/top.txt settings-tmp`); run(`cp -p settings/commands.txt settings-tmp`)
    file = ["vers"=>open("settings-tmp/versions.txt","w"),"urls"=>open("settings-tmp/urls.txt","w"),"paths"=>open("settings-tmp/paths.txt","w")]
    w = 10
    for pkg in get_packages()
        println(file["vers"],rpad(name(pkg),w," "),version(pkg))
        if version(pkg) != "NA"
            PATH = contains(path(pkg),gettop()) ? replace(basename(path(pkg)),version(pkg),"[VER]") : replace(replace(path(pkg),osrelease(),"[OS]"),version(pkg),"[VER]")
            println(file["urls"],rpad(name(pkg),w," "),replace(url(pkg),version(pkg),"[VER]"))
            println(file["paths"],rpad(name(pkg),w," "),PATH)
        else
            println(file["urls"],rpad(name(pkg),w," "),"NA")
            println(file["paths"],rpad(name(pkg),w," "),"NA")
        end
    end
    for (k,v) in file close(v) end
    run(`rm -rf settings`); run(`mv settings-tmp settings`)
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
    check_for_settings()
    if sep <= 1 sep = 1; info("Using min. column spacing of ",string(sep)," spaces.") end
    if sep >= 24 sep = 24; info("Using max. column spacing of ",string(sep)," spaces.") end
    print("\n",Base.text_colors[:bold])
    println("Current build settings",Base.text_colors[:bold])
    try
        println("ID: ",readchomp("settings/id.txt"))
    catch
        println("ID: ","id file not found; This will not affect build.")
    end
    println("TOP: ",gettop())
    println("TAG: ",gettag())
    sizes = [:name=>0,:version=>0,:url=>0]
    for pkg in get_packages()
        for s in [:name,:version,:url]
            sizes[s] = max(sizes[s],length(pkg.(s)))
        end
    end
    w1 = sizes[:name] + sep; w2 = sizes[:version] + sep; w3 = sizes[:url] + sep
    print("\n",Base.text_colors[:bold])
    for k in [:name,:version,:url,:path]; if col != :all && !(k in [:name,col]) continue end
        if k != :path print(rpad(k,sizes[k]+sep," "),Base.text_colors[:bold])
        else print(k,Base.text_colors[:bold]) end
    end
    for k in [:cmds,:deps]; if col == :all || k != col continue end
        print(k,Base.text_colors[:bold])
    end
    print("\n",Base.text_colors[:normal])
    for pkg in get_packages()
        p = replace(path(pkg),string(gettop(),"/"),"")
        if col==:all
            println(rpad(name(pkg),w1," "),rpad(git_version(pkg),w2," "),rpad(url(pkg),w3," "),p)
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
    install_dir = is_external(get_package("hdds")) ? osrelease() : install_dirname()
    test_cmds = [
        "xerces-c" => `$LDD $(path(get_package("xerces-c")))/lib/libxerces-c.$OE`,
        "cernlib" => `ls -lh $(path(get_package("cernlib")))/$(version(get_package("cernlib")))/lib/libgeant321.a`,
        "root" => `root -b -q -l`,
        "evio" => `evio2xml`,
        "ccdb" => `ccdb`,
        "jana" => `jana`,
        "hdds" => `$LDD $(path(get_package("hdds")))/$install_dir/lib/libhdds.so` |> `grep libxerces-c`,
        "sim-recon" => `hd_root`]
    for dep in get_deps([name(pkg)])
        if !success(test_cmds[dep])
            error("'$dep' does not appear to be installed. Please check path if using external installation.
            To build all dependencies, run 'hdpm build' with all packages enabled in 'commands.txt'.\n")
        end
    end
end
function versions_from_xml(path="https://halldweb.jlab.org/dist/version.xml")
    check_for_settings()
    file = path; wasurl = false
    if contains(path,"https://") || contains(path,"http://")
        wasurl = true
        println(); info("downloading $file")
        file = basename(path)
        run(`curl -OL $path`)
    end
    println()
    if !ispath(jlab_top()) info("Browse version xml files at https://halldweb.jlab.org/dist") end
    if ispath(jlab_top()) info("Browse version xml files at /group/halld/www/halldweb/html/dist
Problems? Try ",joinpath(jlab_top(),"version.xml")) end
    if !ispath(file) error(file," does not exist!\n") end
    if !contains(file,".xml") error(file," does not appear to be an xml file!\n") end
    d = readdlm(file,use_mmap=false)
    a = Dict{ASCIIString,ASCIIString}()
    for i=1:size(d,1)
        a[replace(replace(d[i,2],"name=",""),"\"","")] = replace(replace(replace(d[i,3],"version=",""),"/>",""),"\"","")
    end
    a["amptools"] = "NA"; a["geant4"] = "NA"
    a["ccdb"] = is_external(get_package("ccdb")) ? a["ccdb"] : replace(a["ccdb"],a["ccdb"],string(a["ccdb"],".00"))
    vers = readdlm("settings/versions.txt",ASCIIString,use_mmap=false)
    output = open("settings/versions.txt","w")
    for i=1:size(vers,1)
        for (k,v) in a
            if vers[i,1] == k println(output,rpad(k,10," "),v) end
        end
        if !haskey(a,"evio") && vers[i,1] == "evio" println(output,rpad("evio",10," "),vers[i,2]) end
     end
     close(output)
     if wasurl rm(file) end
end
#
end
