using Base.Dates
# keep latest-dist. tarfile/day and today's previous file (Nkeep files in total)
target = joinpath(pwd(),".pkgs")
fbad = filter(r"^sim-recon--.{5}-.{2,3}.tar.gz$",readdir(target))
for file in fbad rm(joinpath(target,file)) end
if length(ARGS) != 2 error("Wrong number of arguments. First argument is number of files to keep; second is OS tag.") end
Nkeep = parse(Int,ARGS[1]); assert(Nkeep>=1)
os_tag = ARGS[2]
todays_date = today()
function get_file_dicts()
    files = Dict{DateTime,ASCIIString}()
    latest_file_each_day = Dict{Date,ASCIIString}()
    dts = Dict{Date,DateTime}()
    for line in readlines(`curl -s https://halldweb.jlab.org/dist/`)
        r = search(line,r"href=\".{25,50}\"")
        if start(r) == 0 continue end
        file = line[start(r)+6:last(r)-1]
        if contains(file,string("-",os_tag,".tar.gz")) && !contains(file,"deps")
            r = search(line,r"(\d{4})-(\d{2})-(\d{2})")
            s = split(line[r],"-")
            y = parse(Int,s[1]); mo = parse(Int,s[2]); d = parse(Int,s[3])
            r = search(line,r"(\d{2}):(\d{2})")
            s = split(line[r],":"); h = parse(Int,s[1]); mi = parse(Int,s[2])
            dt = DateTime(y,mo,d,h,mi)
            files[dt] = file
            d = Date(y,mo,d)
            if !haskey(dts,d) dts[d] = DateTime() end
            if dt > dts[d] dts[d] = dt; latest_file_each_day[d] = file end
        end
    end
    files,latest_file_each_day
end
d1,d2 = get_file_dicts()
sorted_d1_keys = sort(collect(keys(d1)))
sorted_d2_keys = sort(collect(keys(d2)))
for i=1:(length(d2)-Nkeep) delete!(d2,sorted_d2_keys[i]) end
prev_dt = sorted_d1_keys[end-1]
prev_file = (Date(year(prev_dt),month(prev_dt),day(prev_dt)) == todays_date) ? d1[prev_dt] : ""
old_files = Array(ASCIIString,0)
for k in sort(collect(keys(d1)))
    if !(d1[k] in values(d2)) && d1[k] != prev_file
        push!(old_files,d1[k])
    end
end
for file in filter(r"^sim-recon-.{7}-.{5}-.{2,3}.tar.gz$",old_files)
    rm(joinpath(target,file))
end
