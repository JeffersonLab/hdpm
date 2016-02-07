# trigger builds on merge into master
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
JL=/group/halld/Software/ExternalPackages/julia-latest/bin/julia
active=$(cat $cwd/.active)
if test $active == "0"; then
    git pull
    c=$(git log -1 --format="%h")
    c_old=$(cat $cwd/.commit)
    if test $c != $c_old; then
        echo $c > $cwd/.commit
        echo "1" > $cwd/.active
        echo "Triggering build of sim-recon-$c (prev: $c_old)."
        cd $cwd
        $JL build.jl sim-recon c6; $JL cull.jl 10 c6
        $JL write_deps_id.jl nathansparks c6; bash pack.sh c6
        $JL build.jl sim-recon c7; $JL cull.jl 10 c7
        $JL write_deps_id.jl nathansparks c7; bash pack.sh c7
        cd ../osx && bash ship.sh && cd $cwd; $JL cull.jl 10 osx
        $JL build.jl sim-recon u14 f22; $JL cull.jl 5 u14
        $JL write_deps_id.jl nathansparks u14 f22; $JL cull.jl 5 f22
        bash pack.sh u14; bash pack.sh f22
        echo "0" > $cwd/.active
    else
        echo "sim-recon is up-to-date."
    fi
fi
cd $cwd
