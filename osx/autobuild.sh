# trigger builds on merge into master
# requires template named 'osx', with top directory set to 'osx'
# to prebuild dependencies, use 'hdpm build amptools evio jana'
# cron: PATH must be set in crontab
top=../pkgs/osx
cwd=$(pwd)
if ! test -d .sim-recon; then
    git clone https://github.com/JeffersonLab/sim-recon .sim-recon
fi
cd .sim-recon
br=$(git symbolic-ref --short HEAD)
if test $br != master; then
    echo "On branch $br, not master!"
    echo "Switching to master!"; git checkout master
fi
if ! test -f $cwd/.commit; then
    echo "0" > $cwd/.commit
fi
if ! test -f $cwd/.active; then
    echo "0" > $cwd/.active
fi
active=$(cat $cwd/.active)
if test $active == "0"; then
    git pull
    c=$(git log -1 --format="%h")
    c_old=$(cat $cwd/.commit)
    if test $c != $c_old; then
        echo $c > $cwd/.commit
        echo "1" > $cwd/.active
        echo "Triggering build of sim-recon-$c (prev: $c_old)."
        cd $cwd; cd ../
        bash -c "source setup.sh && julia src/hdpm.jl select osx \
            && julia src/hdpm.jl build sim-recon && julia src/distclean.jl yes \
            && cd $cwd && julia write_deps_id.jl && bash pack.sh && bash publish.sh \
            && julia cull.jl 10 osx && rm -rf $top/sim-recon $top/hdds"
        echo "0" > $cwd/.active
    else
        echo "sim-recon is up-to-date."
    fi
fi
cd $cwd
