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
info("Filename format: sim-recon-[commit]-[id_deps]-[os_tag].tar.gz")
info("Available OS tags: c6 (CentOS6), c7 (CentOS7), u14 (Ubuntu14), f22 (Fedora22)")
if length(ARGS) > 1 error("Too many arguments; Use 'hdpm help fetch-dist' to see available arguments.") end
function get_latest_URL(str)
    latest_file = ""; latest_dt = DateTime()
    for line in readlines(`curl -s https://halldweb.jlab.org/dist/`)
        r = search(line,r"href=\".{25,50}\"")
        if start(r) == 0 continue end
        file = line[start(r)+6:last(r)-1]
        if contains(file,os_tag) && (str != "deps" ? contains(file,str):!contains(file,str))
            r = search(line,r"(\d{4})-(\d{2})-(\d{2})")
            s = split(line[r],"-")
            y = parse(Int,s[1]); mo = parse(Int,s[2]); d = parse(Int,s[3])
            r = search(line,r"(\d{2}):(\d{2})")
            s = split(line[r],":"); h = parse(Int,s[1]); mi = parse(Int,s[2])
            dt = DateTime(y,mo,d,h,mi)
            if dt > latest_dt latest_dt = dt; latest_file = file end
        end
    end
    if latest_file == "" error("File not found at https://halldweb.jlab.org/dist for $os_tag OS tag.") end
    URL = string("https://halldweb.jlab.org/dist/",latest_file)
    info("Latest file: $URL"); info("Timestamp: $latest_dt")
    URL
end
if length(ARGS) == 1
    URL = ARGS[1]
    if length(URL) < 4 error("Please provide 4-7 characters to specify a commit hash") end
    if length(URL) >= 4 && length(URL) <= 7 && !ispath(URL)
        URL = get_latest_URL(URL)
    end
end
if length(ARGS) == 0 URL = get_latest_URL("deps") end
isurl = false
if contains(URL,"https://") || contains(URL,"http://") isurl = true end
if isurl && !contains(URL,"https://halldweb.jlab.org/dist")
    warn("$URL is an unfamiliar URL")
end
if !isurl && !contains(URL,"/group/halld/www/halldweb/html/dist")
    warn("$URL is an unfamiliar PATH")
end
parts = split(URL,"-")
if length(parts) != 5 || !contains(URL,"sim-recon") error("Unsupported filename format") end
commit = parts[3]; id_deps = parts[4]; tag = split(parts[end],".")[1]
if tag != os_tag warn("$URL is for $tag distribution, but you are on $os_tag") end
url_deps = replace(URL,commit,"deps")
update_deps = false; update = false
if !ispath(PATH) || (ispath("$PATH/.id-deps-$tag") && id_deps != readchomp("$PATH/.id-deps-$tag"))
    run(`rm -rf $PATH`); update_deps = true; update = true
    mk_cd(PATH); get_unpack_file(url_deps,PATH); get_unpack_file(URL,PATH)
elseif commit != split(split(filter(r"^version_sim-recon-",readdir(PATH))[1],"-")[3],"_")[1]
    run(`rm -rf $PATH/sim-recon`); run(`rm -rf $PATH/hdds`)
    mk_cd(PATH); get_unpack_file(URL,PATH); update = true
else info("Already up-to-date, at commit=$commit") end
if update rm_regex(r"^version_.+",PATH)
    run(`touch $PATH/version_sim-recon-$(commit)_deps-$id_deps`)
end
function update_env_script(fname)
    f = open(fname,"r")
    data = readall(f); close(f)
    p = dirname(dirname(fname)); set = contains(fname,".sh") ? "export" : "setenv"
    tobe_replaced = ["/home/hdpm/pkgs",r"\$GLUEX_TOP/julia-.{5,7}/bin:","/opt/rh/python27"]
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
