using Packages
# download and unpack binaries
os = osrelease()
if contains(os,"CentOS6") || contains(os,"RHEL6") os_tag = "c6"
elseif contains(os,"CentOS7") || contains(os,"RHEL7") os_tag = "c7"
elseif contains(os,"Ubuntu14") os_tag = "u14"
elseif contains(os,"Fedora22") os_tag = "f22"
else warn("Unsupported operating system"); os_tag = os end
PATH = joinpath(gettop(),".dist")
info("Browse binary-distribution tarfiles at https://halldweb.jlab.org/dist")
info("Path on JLab CUE: /group/halld/www/halldweb/html/dist")
info("URL format: sim-recon-[commit]-[id_deps]-[os_tag].tar.gz")
info("Available OS tags: c6 (CentOS6), c7 (CentOS7), u14 (Ubuntu14), f22 (Fedora22)")
if length(ARGS) != 1 error("Requires 1 argument specifying the URL of the sim-recon binary-distribution.") end
URL = ARGS[1]
if !contains(URL,"https://halldweb.jlab.org/dist")
    warn("$URL does not appear to be a proper URL")
end
parts = split(URL,"-")
if length(parts) != 5 error("$URL does not appear to be a proper URL") end
commit = parts[3]; id_deps = parts[4]; tag = split(parts[5],".")[1]
if tag != os_tag warn("$URL is for $tag distribution, but you are on $os_tag") end
if !contains(URL,"sim-recon") || (!contains(URL,".tar.gz") && !contains(URL,".tgz"))
    error("$URL does not appear to be a proper URL")
end
url_deps = replace(URL,commit,"deps")
update_deps = false
if !ispath(PATH) || (ispath("$PATH/.id-deps-$tag") && id_deps != readchomp(`$PATH/.id-deps-$tag`))
    run(`rm -rf $PATH`); update_deps = true
    mk_cd(PATH); get_unpack_file(url_deps,PATH); get_unpack_file(URL,PATH)
else
    run(`rm -rf $PATH/sim-recon`); run(`rm -rf $PATH/hdds`)
    mk_cd(PATH); get_unpack_file(URL,PATH)
end
rm_regex(r"^version_.+",PATH)
run(`touch $PATH/version_sim-recon-$(commit)_deps-$id_deps`)
function update_env_script(fname)
    f = open(fname,"r")
    data = readall(f); close(f)
    p = dirname(dirname(fname)); set = contains(fname,".sh") ? "export" : "setenv"
    tobe_replaced = ["/home/hdpm/pkgs",r"\$GLUEX_TOP/julia-.+/bin:","/opt/rh/python27"]
    replacement = [p,"","\$GLUEX_TOP/opt/rh/python27"]
    for i=1:length(tobe_replaced)
      data = replace(data,tobe_replaced[i],replacement[i])
    end
    res_path = "/u/group/halld/www/halldweb/html/resources"
    if ispath(res_path) data = replace(data,"#$set JANA_RES","$set JANA_RES")
        data = replace(data,"/path/to/resources",res_path)
    end
    g = open(fname,"w")
    write(g,data); close(g)
end
if update_deps
    update_env_script(joinpath(PATH,"env-setup","hdenv.sh"))
    update_env_script(joinpath(PATH,"env-setup","hdenv.csh"))
end
info("Environment setup: source $(joinpath(PATH,"env-setup","hdenv.(c)sh"))")
