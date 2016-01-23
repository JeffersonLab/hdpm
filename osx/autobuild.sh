# trigger builds on merge into master
cwd=`pwd`
if ! test -d .sim-recon; then
    git clone https://github.com/JeffersonLab/sim-recon .sim-recon
fi
cd .sim-recon
br=`git symbolic-ref --short HEAD`
if test $br != master; then
    echo "On branch $br, not master!"
    echo "Switching to master!"; git checkout master
fi
git pull
if ! test -f $cwd/.commit; then
    echo "0" > $cwd/.commit
fi
if ! test -f $cwd/.active; then
    echo "0" > $cwd/.active
fi
c_old=`cat $cwd/.commit`
c=`git log -1 --format="%h"`
active=$(cat $cwd/.active)
if test $c != $c_old && test $active == "0"; then
    echo $c > $cwd/.commit
    echo "1" > $cwd/.active
    echo "Triggering build of sim-recon-$c (prev: $c_old)."
    cd $cwd; cd ../
    bash -c "source setup-osx.sh && julia src/hdpm.jl update hdds sim-recon \
        && julia src/hdpm.jl clean-build hdds sim-recon \
        && julia src/distclean.jl yes && source setup-osx.sh && cd $cwd \
        && julia write_deps_id.jl && bash pack.sh && julia cull.jl 10 osx"
    echo "0" > $cwd/.active
else
    echo "sim-recon is up-to-date."
fi
cd $cwd
