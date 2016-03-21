using Packages
# use versions from version XML file
if length(ARGS) == 0
    versions_from_xml()
elseif length(ARGS) == 1
    versions_from_xml(ARGS[1])
else
    usage_error("Too many arguments: Run 'hdpm help v-xml' to see available arguments.")
end
