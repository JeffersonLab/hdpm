# tcsh: Get julia binary for 64-bit Linux and put it in PATH
#       Make alias for running hdpm.jl 	
# usage: source setup.csh
echo "Hall-D Package Manager setup"
echo "Run the 'hdpm' command in the current working directory."
alias hdpm 'julia src/hdpm.jl'
setenv JULIA_LOAD_PATH `pwd`/src
if ( -d /group/halld/Software/ExternalPackages/julia-latest ) then
    echo "You appear to be on the JLab CUE; Will try to use group installation of julia."
    if ( ! -e ~/bin/julia ) then
	echo "Making a link to group installation."
	ln -s /u/group/halld/Software/ExternalPackages/julia-latest/bin/julia ~/bin/julia
	rehash; goto end
    endif
    echo "You already have julia in your PATH."; goto end
endif
set VER=0.3.10
set JLPATH=`pwd`/pkgs/deps/julia-$VER/bin
if ( -d pkgs/deps/julia-$VER ) then
    echo "julia directory already exists; nothing to download. If not in PATH, use:" 
    echo "setenv PATH ${JLPATH}:"\$PATH""; goto end
endif
echo "Downloading julia."
curl -OL https://julialang.s3.amazonaws.com/bin/linux/x64/0.3/julia-$VER-linux-x86_64.tar.gz
mkdir -p pkgs/deps/julia-$VER
tar -xzf julia-$VER-linux-x86_64.tar.gz -C pkgs/deps/julia-$VER --strip-components=1
rm -f julia-$VER-linux-x86_64.tar.gz
echo "Putting julia in your PATH."
setenv PATH ${JLPATH}:$PATH
end:
    echo "Good to go!"
