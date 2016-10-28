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
        for tag in c6 c7 u14; do
             gxd build --rmi sim-recon $tag; gxd write -u nathansparks $tag
             bash pack.sh $tag; gxd cull -n 5 $tag
        done
        echo "0" > $cwd/.active
    else
        echo "sim-recon is up-to-date."
    fi
fi
cd $cwd
