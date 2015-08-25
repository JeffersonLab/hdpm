# unified interface
using Packages
template_ids = get_template_ids()
pkg_names = get_pkg_names()
pkg_cols = ["version", "url", "path", "deps", "cmds"]
if length(ARGS) == 0 || (length(ARGS) == 1 && ARGS[1] == "help")
    println("usage: hdpm <command> [<args>]
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
Use 'hdpm help <command>' to see available arguments.")
end
if length(ARGS) == 1 && ARGS[1] != "help"
    if ARGS[1] == "select"
        run(`julia src/select_template.jl`)
    elseif ARGS[1] == "fetch"
        run(`julia src/copkgs.jl`)
    elseif ARGS[1] == "build"
        run(`julia src/copkgs.jl`)
        run(`julia src/mkpkgs.jl`)
    elseif ARGS[1] == "update"
        run(`julia src/update.jl`)
    elseif ARGS[1] == "clean"
        run(`julia src/clean.jl`)
    elseif ARGS[1] == "clean-build"
        run(`julia src/clean.jl`)
        run(`julia src/mkpkgs.jl`)
    elseif ARGS[1] == "show"
        run(`julia src/show_settings.jl`)
    elseif ARGS[1] == "v-xml"
        run(`julia src/versions_from_xml.jl`)
    elseif ARGS[1] == "select" || ARGS[1] == "save"
        error("Requires one argument. Use 'hdpm help $(ARGS[1])' to see available arguments.\n")
    else
        error("Unknown command. Use 'hdpm help' to see available commands.\n")
    end
end
if length(ARGS) == 2 && ARGS[1] == "help"
    if ARGS[2] == "select"
        println("Select the desired build template")
        println("usage: hdpm select |<template id>|")
        println("ids: ",string(template_ids))
    elseif ARGS[2] == "save"
        println("Save the current settings as a new build template")
        println("usage: hdpm save <new template id>")
    elseif ARGS[2] == "show"
        println("Show the current build settings")
        println("usage: hdpm show |<column>| |<column spacing>|")
        println("columns: version, url, path, deps, cmds")
    elseif ARGS[2] == "fetch"
        println("Checkout/download/clone the selected packages")
        println("usage: hdpm fetch |<pkgs>...|")
        println("pkgs: ",string(pkg_names))
    elseif ARGS[2] == "build"
        println("Build the selected packages (fetch if needed)")
        println("usage: hdpm build |<template id>|")
        println("ids: ",string(template_ids))
        println("usage: hdpm build |<xml file url or path>|")
        println("usage: hdpm build |<pkgs>...|")
        println("pkgs: ",string(pkg_names))
    elseif ARGS[2] == "update"
        println("Update selected Git/SVN packages")
        println("usage: hdpm update |<pkgs>...|")
        println("pkgs: ",string(pkg_names))
    elseif ARGS[2] == "clean"
        println("Remove build products of the selected packages")
        println("usage: hdpm clean |<pkgs>...|")
        println("pkgs: ",string(pkg_names))
    elseif ARGS[2] == "clean-build"
        println("Do a clean build of the selected packages")
        println("usage: hdpm clean-build |<pkgs>...|")
        println("pkgs: ",string(pkg_names))
    elseif ARGS[2] == "v-xml"
        println("Replace versions with versions from a version XML file")
        println("usage: hdpm v-xml |<url or path>|")
    else
        error("Unknown command. Use 'hdpm help' to see available commands.\n")
    end
end
if length(ARGS) == 2 && ARGS[1] != "help"
    if ARGS[1] == "select"
        if ARGS[2] in template_ids
            run(`julia src/select_template.jl $(ARGS[2])`)
        else error("Unknown argument. Use 'hdpm help $(ARGS[1])' to see available arguments.\n") end
    elseif ARGS[1] == "save"
        run(`julia src/mk_template.jl $(ARGS[2])`)
    elseif ARGS[1] == "build"
        if ARGS[2] in template_ids && ARGS[2] in pkg_names
            error("template id cannot be the same as a package name. Please rename the template id.\n")
        end
        if ARGS[2] in template_ids
            run(`julia src/select_template.jl $(ARGS[2])`)
            run(`julia src/copkgs.jl`)
            run(`julia src/mkpkgs.jl`)
        elseif ARGS[2] in pkg_names
            run(`julia src/copkgs.jl $(ARGS[2])`)
            run(`julia src/mkpkgs.jl $(ARGS[2])`)
        elseif contains(ARGS[2],".xml")
            run(`julia src/select_template.jl`)
            run(`julia src/versions_from_xml.jl $(ARGS[2])`)
            run(`julia src/copkgs.jl`)
            run(`julia src/mkpkgs.jl`)
        else
            error("Unknown argument. Use 'hdpm help $(ARGS[1])' to see available arguments.\n")
        end
    elseif ARGS[1] == "show"
        if ARGS[2] in pkg_cols
            run(`julia src/show_settings.jl $(ARGS[2])`)
        else
            error("Unknown argument. Use 'hdpm help $(ARGS[1])' to see available arguments.\n")
        end
    elseif ARGS[1] == "fetch" && ARGS[2] in pkg_names
        run(`julia src/copkgs.jl $(ARGS[2])`)
    elseif ARGS[1] == "clean" && ARGS[2] in pkg_names
        run(`julia src/clean.jl $(ARGS[2])`)
    elseif ARGS[1] == "clean-build" && ARGS[2] in pkg_names
        run(`julia src/clean.jl $(ARGS[2])`)
        run(`julia src/mkpkgs.jl $(ARGS[2])`)
    elseif ARGS[1] == "update" && ARGS[2] in pkg_names
        run(`julia src/update.jl $(ARGS[2])`)
    elseif ARGS[1] == "v-xml"
        run(`julia src/versions_from_xml.jl $(ARGS[2])`)
    else
        error("Unknown command. Use 'hdpm help' to see available commands.\n")
    end
end
if length(ARGS) == 3 && ARGS[1] == "show"
    if ARGS[2] in pkg_cols || ARGS[3] in pkg_cols
        run(`julia src/show_settings.jl $(ARGS[2]) $(ARGS[3])`)
    else
        error("Unknown argument. Use 'hdpm help $(ARGS[1])' to see available arguments.\n")
    end
end
if length(ARGS) >= 3 && (length(ARGS) <= length(pkg_names) + 1) && ARGS[1] != "show"
    trouble = false
    for i=2:length(ARGS)
        if !(ARGS[i] in pkg_names) warn("Unknown argument: ",ARGS[i]); trouble = true end
    end
    if trouble error("Unknown argument (typo?). Use 'hdpm help $(ARGS[1])' to see available arguments.\n") end
    #
    nargs = ``
    for i=2:length(ARGS)
        nargs = `$nargs $(ARGS[i])`
    end
    if ARGS[1] == "fetch"
        run(`julia src/copkgs.jl $nargs`)
    elseif ARGS[1] == "build"
        run(`julia src/copkgs.jl $nargs`)
        run(`julia src/mkpkgs.jl $nargs`)
    elseif ARGS[1] == "update"
        run(`julia src/update.jl $nargs`)
    elseif ARGS[1] == "clean"
        run(`julia src/clean.jl $nargs`)
    elseif ARGS[1] == "clean-build"
        run(`julia src/clean.jl $nargs`)
        run(`julia src/mkpkgs.jl $nargs`)
    else
        error("Unknown command. Use 'hdpm help' to see available commands.\n")
    end
elseif length(ARGS) > length(pkg_names) + 1
    error("Too many arguments. Use 'hdpm help $(ARGS[1])' to see available arguments.\n")
end
