# bash: Get julia binary and put it in PATH
#       Make alias for running hdpm.jl
# usage: source setup.sh
echo "Hall-D Package Manager setup"
SCRIPT_PATH=${BASH_SOURCE[0]}
if [[ -h $SCRIPT_PATH ]]; then
    SCRIPT_PATH=$(readlink $SCRIPT_PATH)
fi
HDPM_PATH=$(cd $(dirname $SCRIPT_PATH); pwd)
echo "Run 'hdpm' to see available commands."
alias hdpm="julia $HDPM_PATH/src/hdpm.jl"
export JULIA_LOAD_PATH=$HDPM_PATH/src
uname=$(uname)
if [[ $uname == "Linux" ]]; then
    JLPATH=/group/halld/Software/ExternalPackages/julia-latest/bin
    if [[ -f $JLPATH/julia ]]; then
        echo "You appear to be on the JLab CUE; Will try to use group installation of julia."
        echo $PATH | grep -q $JLPATH
        if [[ $? -ne 0 ]]; then
            echo "Putting julia in your PATH."
            export PATH=$JLPATH:$PATH; echo "Good to go!"; return
        else
            echo "You already have julia in your PATH."; echo "Good to go!"; return
        fi
    fi
fi
VER=0.4.5
JLPATH=$HDPM_PATH/pkgs/julia-$VER/bin
if [[ -f $JLPATH/julia ]]; then
    echo "julia-$VER directory already exists; nothing to download."
    echo $PATH | grep -q $JLPATH
    if [[ $? -ne 0 ]]; then
        echo "Putting julia in your PATH."
        export PATH=$JLPATH:$PATH; echo "Good to go!"; return
    else
        echo "You already have julia in your PATH."; echo "Good to go!"; return
    fi
fi
echo "Downloading julia-$VER."
if [[ $uname == "Linux" ]]; then
    curl -OL https://julialang.s3.amazonaws.com/bin/linux/x64/0.4/julia-$VER-linux-x86_64.tar.gz
    mkdir -p $HDPM_PATH/pkgs/julia-$VER
    tar -xzf julia-$VER-linux-x86_64.tar.gz -C $HDPM_PATH/pkgs/julia-$VER --strip-components=1
    rm -f julia-$VER-linux-x86_64.tar.gz
fi
if [[ $uname == "Darwin" ]]; then
    curl -OL https://s3.amazonaws.com/julialang/bin/osx/x64/0.4/julia-$VER-osx10.7+.dmg
    hdiutil attach -quiet julia-$VER-osx10.7+.dmg
    mkdir -p $HDPM_PATH/pkgs
    cp -pr /Volumes/Julia/Julia-$VER.app/Contents/Resources/julia $HDPM_PATH/pkgs/julia-$VER
    hdiutil detach -quiet /Volumes/Julia
    rm -f $HDPM_PATH/pkgs/julia-$VER/etc/julia/juliarc.jl
    rm -f julia-$VER-osx10.7+.dmg
fi
if [[ -f $JLPATH/julia ]]; then
    echo "Putting julia in your PATH."
    export PATH=$JLPATH:$PATH
    echo "Good to go!"
else
    echo "julia download failed: Source this setup script to try again."
    echo "If the problem persists, please check your internet connection."
fi
