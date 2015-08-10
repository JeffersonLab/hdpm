# tcsh: Get julia binary for Mac OS X and put it in PATH
#       Make alias for running hdpm.jl
# usage: source setup-osx.csh
echo "Mac OS X (10.7+): Hall-D Package Manager setup"
echo "Run the 'hdpm' command in the current working directory."
alias hdpm 'julia src/hdpm.jl'
setenv JULIA_LOAD_PATH `pwd`/src
set VER=0.3.11
set JLPATH=`pwd`/pkgs/deps/julia-$VER/bin
if ( -e ${JLPATH}/julia ) then
    echo "julia-$VER directory already exists; nothing to download."
    echo $PATH | grep -q $JLPATH
    if ( $? != 0 ) then
	echo "Putting julia in your PATH."
	setenv PATH ${JLPATH}:$PATH; echo "Good to go!"; goto end
    else
        echo "You already have julia in your PATH."; echo "Good to go!"; goto end
    endif
endif
echo "Downloading julia-$VER."
curl -OL https://s3.amazonaws.com/julialang/bin/osx/x64/0.3/julia-$VER-osx10.7+.dmg
hdiutil attach -quiet julia-$VER-osx10.7+.dmg
mkdir -p pkgs/deps
cp -pr /Volumes/Julia/Julia-$VER.app/Contents/Resources/julia pkgs/deps/julia-$VER
hdiutil detach -quiet /Volumes/Julia
rm -f julia-$VER-osx10.7+.dmg
echo "Putting julia in your PATH."
setenv PATH ${JLPATH}:$PATH
end:
    echo "Good to go!"
