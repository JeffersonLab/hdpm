using Packages
# select build settings template
if length(ARGS) == 0
    select_template()
elseif length(ARGS) == 1
    select_template(ARGS[1])
else
    error("too many arguments. Run 'hdpm help select' to see available arguments.")
end
