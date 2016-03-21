using Packages
# display settings
if length(ARGS) == 0
    show_settings()
elseif length(ARGS) == 1
    try
        show_settings(sep=parse(Int,ARGS[1]))
    catch
        show_settings(col=symbol(ARGS[1]))
    end
elseif length(ARGS) == 2
    try
        show_settings(col=symbol(ARGS[1]),sep=parse(Int,ARGS[2]))
    catch
        try
            show_settings(col=symbol(ARGS[2]),sep=parse(Int,ARGS[1]))
        catch
            usage_error("If 2 arguments are given, 1 arg. must be an integer specifying the minimum spacing between columns, the other the column name.")
        end
    end
else
    usage_error("Too many arguments (>2):\n\tUse no args. to show a column for each setting with the default minimum spacing of 8 spaces; use 1 arg. to specify a column name or spacing; use 2 args. to specify both.")
end
