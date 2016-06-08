# unified interface
using Packages
const home = dirname(dirname(@__FILE__))
const template_ids = get_template_ids()
const pkg_names = get_pkg_names()
const pkg_cols = ["version", "url", "path", "deps", "cmds"]
if length(ARGS) == 0 || (length(ARGS) == 1 && ARGS[1] == "help")
    hz("="); println("usage: hdpm <command> [<args>]
    commands:
    \t help        Show available commands
    \t select      Select a 'build template' from the 'templates' directory
    \t save        Save the current settings as a new 'build template'
    \t show        Show selected packages
    \t fetch       Checkout/download/clone selected packages
    \t build       Build selected packages (fetch if needed)
    \t update      Update selected Git/SVN packages
    \t clean       Completely remove build products of selected packages
    \t clean-build Clean build of selected packages
    \t v-xml       Replace versions with versions from a version XML file
    \t install     Install binary distribution of sim-recon and its deps
    \t run         Run a command in the Hall-D offline environment
--------------------------------------------------------------------------------
Use 'hdpm help <command>' to see available arguments."); hz("=")
end
const os = osrelease()
if contains(os,"CentOS6") || contains(os,"RHEL6") && ispath(joinpath(gettop(),".dist"))
    p = joinpath(gettop(),".dist")
    ENV["PATH"] = string("$p/opt/rh/python27/root/usr/bin:$p/opt/rh/devtoolset-3/root/usr/bin:",ENV["PATH"])
    a = "$p/opt/rh/python27/root/usr/lib64:$p/opt/rh/devtoolset-3/root/usr/lib64:$p/opt/rh/devtoolset-3/root/usr/lib"
    if !haskey(ENV,"LD_LIBRARY_PATH")
        ENV["LD_LIBRARY_PATH"] = a
    else
        ENV["LD_LIBRARY_PATH"] = string(a,":",ENV["LD_LIBRARY_PATH"])
    end
end
if length(ARGS) == 1 && ARGS[1] != "help"
    if ARGS[1] == "select"
        run(`julia $home/src/select_template.jl`)
    elseif ARGS[1] == "fetch"
        run(`julia $home/src/copkgs.jl`)
    elseif ARGS[1] == "build"
        run(`julia $home/src/copkgs.jl`)
        run(`julia $home/src/mkpkgs.jl`)
    elseif ARGS[1] == "update"
        run(`julia $home/src/update.jl`)
    elseif ARGS[1] == "clean"
        run(`julia $home/src/clean.jl`)
    elseif ARGS[1] == "clean-build"
        run(`julia $home/src/clean.jl`)
        run(`julia $home/src/mkpkgs.jl`)
    elseif ARGS[1] == "show"
        run(`julia $home/src/show_settings.jl`)
    elseif ARGS[1] == "v-xml"
        run(`julia $home/src/versions_from_xml.jl`)
    elseif ARGS[1] == "install"
        run(`julia $home/src/fetch_dist.jl`)
        run(`julia $home/src/install_dist.jl`)
    elseif ARGS[1] == "run"
        run(`julia $home/src/run.jl`)
    elseif ARGS[1] == "save"
        usage_error("'hdpm save' requires one argument.\n\tUse 'hdpm help $(ARGS[1])' to see available arguments.")
    else
        usage_error("'$(ARGS[1])' is not a hdpm command.\n\tUse 'hdpm help' to see available commands.")
    end
end
if length(ARGS) == 2 && ARGS[1] == "help"
    if ARGS[2] == "select" hz("=")
        println("Select the desired build template"); hz("-")
        println("usage: hdpm select\n\t(select the master template)"); hz("-")
        println("usage: hdpm select [<template-id>]")
        println("ids:   ",join(template_ids,", ")); hz("=")
    elseif ARGS[2] == "save" hz("=")
        println("Save the current settings as a new build template"); hz("-")
        println("usage: hdpm save <template-id>"); hz("-")
        println("If <template-id> is set to 'jlab',\ngenerate a base template with all package builds disabled.")
        println("Create JLab base template to use halld group installations:")
        println("       hdpm save jlab"); hz("=")
    elseif ARGS[2] == "show" hz("=")
        println("Show the current build settings"); hz("-")
        println("usage:   hdpm show\n\t(show version and path settings)"); hz("-")
        println("usage:   hdpm show [<column>] [<column-spacing>]")
        println("columns: version, url, path, deps, cmds"); hz("=")
    elseif ARGS[2] == "fetch" hz("=")
        println("Checkout/download/clone the selected packages"); hz("-")
        println("usage: hdpm fetch\n\t(fetch all enabled packages)"); hz("-")
        println("usage: hdpm fetch [<pkg>...]")
        println("pkgs:  ",join(pkg_names,",")); hz("=")
    elseif ARGS[2] == "build" hz("=")
        println("Build the selected packages (fetch if needed)")
        println("Display build information if a package is already built"); hz("-")
        println("usage: hdpm build\n\t(build all enabled packages)"); hz("-")
        println("usage: hdpm build [<pkg>...]")
        println("pkgs:  ",join(pkg_names,",")); hz("-")
        println("usage: hdpm build [<template-id>]")
        println("ids:   ",join(template_ids,", ")); hz("-")
        println("usage: hdpm build [<xmlfile-url> | <xmlfile-path>]"); hz("=")
        println("Build a sim-recon subdirectory (use '-c' option to clean)"); hz("-")
        println("usage: hdpm build [-c] <subdirectory>"); hz("=")
    elseif ARGS[2] == "update" hz("=")
        println("Update selected Git/SVN packages"); hz("-")
        println("usage: hdpm update\n\t(update all enabled repository packages)"); hz("-")
        println("usage: hdpm update [<pkg>...]")
        println("pkgs:  ",join(pkg_names,",")); hz("=")
    elseif ARGS[2] == "clean" hz("=")
        println("Remove build products of the selected packages"); hz("-")
        println("usage: hdpm clean\n\t(clean all enabled packages)"); hz("-")
        println("usage: hdpm clean [<pkg>...]")
        println("pkgs:  ccdb, jana, hdds, sim-recon"); hz("=")
    elseif ARGS[2] == "clean-build" hz("=")
        println("Do a clean build of the selected packages"); hz("-")
        println("usage: hdpm clean-build\n\t(clean-build of all enabled packages)"); hz("-")
        println("usage: hdpm clean-build [<pkg>...]")
        println("pkgs:  ccdb, jana, hdds, sim-recon"); hz("=")
    elseif ARGS[2] == "v-xml" hz("=")
        println("Replace versions with versions from a version XML file"); hz("-")
        println("usage: hdpm v-xml\n\t(w/ https://halldweb.jlab.org/dist/version.xml)"); hz("-")
        println("usage: hdpm v-xml [<xmlfile-url> | <xmlfile-path>]"); hz("=")
    elseif ARGS[2] == "install" hz("=")
        println("Install binary distribution of sim-recon and its deps"); hz("-")
        println("usage: hdpm install\n\t(install latest binary distribution)"); hz("-")
        println("usage: hdpm install [<commit>]"); hz("-")
        println("usage: hdpm install [<tarfile-url> | <tarfile-path>]"); hz("-")
        println("List available binary distribution tarfiles"); hz("-")
        println("usage: hdpm install -l"); hz("=")
    elseif ARGS[2] == "run" hz("=")
        println("Run a command in the Hall-D offline environment"); hz("-")
        println("usage: hdpm run\n\t(run a bash shell for interactive work)"); hz("-")
        println("usage: hdpm run \"<cmd>\""); hz("=")
    elseif ARGS[2] == "help" hz("=")
        println("Show available commands"); hz("-"); println("usage: hdpm help"); hz("=")
    else
        usage_error("'$(ARGS[2])' is not a hdpm command. Use 'hdpm help' to see available commands.")
    end
end
if length(ARGS) == 2 && ARGS[1] != "help"
    if ARGS[1] == "select"
        if ARGS[2] in template_ids
            run(`julia $home/src/select_template.jl $(ARGS[2])`)
        else usage_error("'$(ARGS[2])' is not a valid template id.\n\tUse 'hdpm help $(ARGS[1])' to see available template ids.") end
    elseif ARGS[1] == "save"
        run(`julia $home/src/mk_template.jl $(ARGS[2])`)
    elseif ARGS[1] == "build"
        if ARGS[2] in template_ids && ARGS[2] in pkg_names
            usage_error("'$(ARGS[2])' template id has the same name as a package.\n\tRename this template id.")
        end
        if ispath(ARGS[2]) && contains(ARGS[2],"sim-recon")
            run(`julia $home/src/build_dir.jl $(ARGS[2])`)
        elseif ARGS[2] in template_ids
            run(`julia $home/src/select_template.jl $(ARGS[2])`)
            run(`julia $home/src/copkgs.jl`)
            run(`julia $home/src/mkpkgs.jl`)
        elseif ARGS[2] in pkg_names
            run(`julia $home/src/copkgs.jl $(ARGS[2])`)
            run(`julia $home/src/mkpkgs.jl $(ARGS[2])`)
        elseif contains(ARGS[2],".xml")
            run(`julia $home/src/select_template.jl`)
            run(`julia $home/src/versions_from_xml.jl $(ARGS[2])`)
            run(`julia $home/src/copkgs.jl`)
            run(`julia $home/src/mkpkgs.jl`)
        else
            usage_error("'$(ARGS[2])' is not a valid argument.\n\tUse 'hdpm help $(ARGS[1])' to see available arguments.")
        end
    elseif ARGS[1] == "show"
        if ARGS[2] in pkg_cols || isinteger(parse(Int,ARGS[2]))
            run(`julia $home/src/show_settings.jl $(ARGS[2])`)
        else
            usage_error("'$(ARGS[2])' is not a valid argument.\n\tUse 'hdpm help $(ARGS[1])' to see available arguments.")
        end
    elseif ARGS[1] == "fetch" && ARGS[2] in pkg_names
        run(`julia $home/src/copkgs.jl $(ARGS[2])`)
    elseif ARGS[1] == "clean" && ARGS[2] in pkg_names
        run(`julia $home/src/clean.jl $(ARGS[2])`)
    elseif ARGS[1] == "clean-build" && ARGS[2] in pkg_names
        run(`julia $home/src/clean.jl $(ARGS[2])`)
        run(`julia $home/src/mkpkgs.jl $(ARGS[2])`)
    elseif ARGS[1] == "update" && ARGS[2] in pkg_names
        run(`julia $home/src/update.jl $(ARGS[2])`)
    elseif ARGS[1] == "v-xml"
        run(`julia $home/src/versions_from_xml.jl $(ARGS[2])`)
    elseif ARGS[1] == "install"
        run(`julia $home/src/fetch_dist.jl $(ARGS[2])`)
        if ARGS[2] != "-l" run(`julia $home/src/install_dist.jl`) end
    elseif ARGS[1] == "run"
        run(`julia $home/src/run.jl $(ARGS[2])`)
    else
        usage_error("'$(ARGS[1])' is not a hdpm command.\n\tUse 'hdpm help' to see available commands.")
    end
end
if length(ARGS) == 3 && ARGS[1] == "build" && ARGS[2] == "-c"
    if ispath(ARGS[3]) && contains(ARGS[3],"sim-recon")
        run(`julia $home/src/build_dir.jl -c $(ARGS[3])`); exit()
    else
        usage_error("'$(ARGS[3])' is not a valid argument.\n\tUse 'hdpm help $(ARGS[1])' to see available arguments.")
    end
end
if length(ARGS) == 3 && ARGS[1] == "show"
    if ARGS[2] in pkg_cols || ARGS[3] in pkg_cols
        run(`julia $home/src/show_settings.jl $(ARGS[2]) $(ARGS[3])`)
    else
        usage_error("'$(ARGS[2])' or '$(ARGS[3])' is not a valid argument.\n\tUse 'hdpm help $(ARGS[1])' to see available arguments.")
    end
end
if length(ARGS) > 3 && ARGS[1] == "show"
    usage_error("Too many arguments: Use 'hdpm help $(ARGS[1])' to see available arguments.")
end
if length(ARGS) >= 3 && (length(ARGS) <= length(pkg_names) + 1) && ARGS[1] != "show"
    trouble = false
    for i=2:length(ARGS)
        if !(ARGS[i] in pkg_names) warn("Unknown argument: ",ARGS[i]); trouble = true end
    end
    if trouble usage_error("Unknown argument(s) (typo?).\n\tUse 'hdpm help $(ARGS[1])' to see available arguments.") end
    #
    nargs = ``
    for i=2:length(ARGS)
        nargs = `$nargs $(ARGS[i])`
    end
    if ARGS[1] == "fetch"
        run(`julia $home/src/copkgs.jl $nargs`)
    elseif ARGS[1] == "build"
        run(`julia $home/src/copkgs.jl $nargs`)
        run(`julia $home/src/mkpkgs.jl $nargs`)
    elseif ARGS[1] == "update"
        run(`julia $home/src/update.jl $nargs`)
    elseif ARGS[1] == "clean"
        run(`julia $home/src/clean.jl $nargs`)
    elseif ARGS[1] == "clean-build"
        run(`julia $home/src/clean.jl $nargs`)
        run(`julia $home/src/mkpkgs.jl $nargs`)
    else
        usage_error("'$(ARGS[1])' is not a hdpm command.\n\tUse 'hdpm help' to see available commands.")
    end
elseif length(ARGS) > length(pkg_names) + 1
    usage_error("Too many arguments: Use 'hdpm help $(ARGS[1])' to see available arguments.")
end
