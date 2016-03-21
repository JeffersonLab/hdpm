using Packages
if length(ARGS) == 1
    mk_template(ARGS[1])
else
    usage_error("'hdpm save' takes only one argument, a new template id.")
end
