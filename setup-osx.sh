# bash: Get julia binary for Mac OS X and put it in PATH
#       Make alias for running hdpm.jl
# usage: source setup-osx.sh
echo "Mac OS X (10.7+): Hall-D Package Manager setup"
echo "Run the 'hdpm' command in the current working directory."
alias hdpm='julia src/hdpm.jl'
export JULIA_LOAD_PATH=`pwd`/src
VER=0.4.3
JLPATH=`pwd`/pkgs/julia-$VER/bin
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
curl -OL https://s3.amazonaws.com/julialang/bin/osx/x64/0.4/julia-$VER-osx10.7+.dmg
hdiutil attach -quiet julia-$VER-osx10.7+.dmg
mkdir -p pkgs
cp -pr /Volumes/Julia/Julia-$VER.app/Contents/Resources/julia pkgs/julia-$VER
hdiutil detach -quiet /Volumes/Julia
rm -f pkgs/julia-$VER/etc/julia/juliarc.jl
rm -f julia-$VER-osx10.7+.dmg
echo "Putting julia in your PATH."
export PATH=$JLPATH:$PATH
echo "Good to go!"
