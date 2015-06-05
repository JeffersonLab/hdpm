using Packages
# select build settings template
if length(ARGS) == 1 
    select_template(ARGS[1]) 
else 
    throw("Error: Requires 1 argument specifying the id of the build settings template.")
end
