# unified interface
if length(ARGS) == 0 || (length(ARGS) == 1 && ARGS[1] == "help")
    println("usage: hdpm <command> [<args>]
    commands:
    \t help        Show available commands
    \t select      Select a 'build template' from the 'templates' directory
    \t co          Checkout/download selected packages
    \t build       Build selected packages
    \t install     Checkout/download and build selected packages
    \t show        Show selected packages
    \t update      Update selected packages
    \t clean-build Clean build of selected packages") 
end
if length(ARGS) == 1 && ARGS[1] != "help"
    if ARGS[1] == "co"
        run(`julia src/copkgs.jl`) 
    elseif ARGS[1] == "build"
        run(`julia src/mkpkgs.jl`) 
    elseif ARGS[1] == "install"
        run(`julia src/copkgs.jl`) 
        run(`julia src/mkpkgs.jl`) 
    elseif ARGS[1] == "update"
        run(`julia src/update.jl`) 
    elseif ARGS[1] == "clean-build"
        run(`julia src/clean_build.jl`) 
    elseif ARGS[1] == "show"
        run(`julia src/show_settings.jl`) 
    else
        error("Unknown command. Use 'hdpm help' to see available commands.")
    end
end
if length(ARGS) == 2
    if ARGS[1] == "select"
        id = ARGS[2]
        run(`julia src/select_template.jl $id`) 
    elseif ARGS[1] == "install"
        id = ARGS[2]
        run(`julia src/select_template.jl $id`) 
        run(`julia src/copkgs.jl`)
        run(`julia src/mkpkgs.jl`) 
    elseif ARGS[1] == "show"
        run(`julia src/show_settings.jl $(ARGS[2])`) 
    else
        error("Unknown command. Use 'hdpm help' to see available commands.")
    end
end
if length(ARGS) == 3
    if ARGS[1] == "show"
        run(`julia src/show_settings.jl $(ARGS[2]) $(ARGS[3])`) 
    else
        error("Unknown command. Use 'hdpm help' to see available commands.")
    end
end
if length(ARGS) > 3
    error("Unknown command. Use 'hdpm help' to see available commands.")
end
