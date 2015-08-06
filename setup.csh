# tcsh: Get julia binary for 64-bit Linux and put it in PATH
#       Make alias for running hdpm.jl
# usage: source setup.csh
echo "Linux (64-bit): Hall-D Package Manager setup"
echo "Run the 'hdpm' command in the current working directory."
alias hdpm 'julia src/hdpm.jl'
setenv JULIA_LOAD_PATH `pwd`/src
set JLPATH=/group/halld/Software/ExternalPackages/julia-latest/bin
if ( -e ${JLPATH}/julia ) then
    echo "You appear to be on the JLab CUE; Will try to use group installation of julia."
    echo $PATH | grep -q $JLPATH
    if ( $? != 0 ) then
	echo "Putting julia in your PATH."
	setenv PATH ${JLPATH}:$PATH; goto end
    else
        echo "You already have julia in your PATH."; goto end
    endif
endif
set VER=0.3.11
set JLPATH=`pwd`/pkgs/deps/julia-$VER/bin
if ( -e ${JLPATH}/julia ) then
    echo "julia-$VER directory already exists; nothing to download."
    echo $PATH | grep -q $JLPATH
    if ( $? != 0 ) then
	echo "Putting julia in your PATH."
	setenv PATH ${JLPATH}:$PATH; goto end
    else
        echo "You already have julia in your PATH."; goto end
    endif
endif
echo "Downloading julia-$VER."
curl -OL https://julialang.s3.amazonaws.com/bin/linux/x64/0.3/julia-$VER-linux-x86_64.tar.gz
mkdir -p pkgs/deps/julia-$VER
tar -xzf julia-$VER-linux-x86_64.tar.gz -C pkgs/deps/julia-$VER --strip-components=1
rm -f julia-$VER-linux-x86_64.tar.gz
echo "Putting julia in your PATH."
setenv PATH ${JLPATH}:$PATH
end:
    echo "Good to go!"
