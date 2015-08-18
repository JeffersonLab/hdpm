using Packages
if length(ARGS) == 1
    mk_template(ARGS[1])
else
    error("takes only one argument: a new template id")
end
