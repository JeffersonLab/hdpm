# bash: Get julia binary and put it in PATH
#       Make alias for running hdpm.jl
# usage: source setup.sh
echo "Hall-D Package Manager setup"
echo "Run the 'hdpm' command in the current working directory."
alias hdpm='julia src/hdpm.jl'
export JULIA_LOAD_PATH=$(pwd)/src
uname=$(uname)
if test $uname == "Linux"; then
    JLPATH=/group/halld/Software/ExternalPackages/julia-latest/bin
    if test -f $JLPATH/julia
    then
        echo "You appear to be on the JLab CUE; Will try to use group installation of julia."
        echo $PATH | grep -q $JLPATH
        if test $? -ne 0; then
            echo "Putting julia in your PATH."
            export PATH=$JLPATH:$PATH; echo "Good to go!"; return
        else
            echo "You already have julia in your PATH."; echo "Good to go!"; return
        fi
    fi
fi
VER=0.4.3
JLPATH=$(pwd)/pkgs/julia-$VER/bin
if test -f $JLPATH/julia
then
    echo "julia-$VER directory already exists; nothing to download."
    echo $PATH | grep -q $JLPATH
    if test $? -ne 0; then
        echo "Putting julia in your PATH."
        export PATH=$JLPATH:$PATH; echo "Good to go!"; return
    else
        echo "You already have julia in your PATH."; echo "Good to go!"; return
    fi
fi
echo "Downloading julia-$VER."
if test $uname == "Linux"; then
    curl -OL https://julialang.s3.amazonaws.com/bin/linux/x64/0.4/julia-$VER-linux-x86_64.tar.gz
    mkdir -p pkgs/julia-$VER
    tar -xzf julia-$VER-linux-x86_64.tar.gz -C pkgs/julia-$VER --strip-components=1
    rm -f julia-$VER-linux-x86_64.tar.gz
fi
if test $uname == "Darwin"; then
    curl -OL https://s3.amazonaws.com/julialang/bin/osx/x64/0.4/julia-$VER-osx10.7+.dmg
    hdiutil attach -quiet julia-$VER-osx10.7+.dmg
    mkdir -p pkgs
    cp -pr /Volumes/Julia/Julia-$VER.app/Contents/Resources/julia pkgs/julia-$VER
    hdiutil detach -quiet /Volumes/Julia
    rm -f pkgs/julia-$VER/etc/julia/juliarc.jl
    rm -f julia-$VER-osx10.7+.dmg
fi
echo "Putting julia in your PATH."
export PATH=$JLPATH:$PATH
echo "Good to go!"
