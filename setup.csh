# tcsh: Add hdpm binary to PATH
# usage: source setup.csh
set ARGS=($_)
if ("$ARGS" != "") then
    set HDPM_PATH="`dirname ${ARGS[2]}`"
    set HDPM_PATH="`cd $HDPM_PATH; pwd`"
else
    if ( -e setup.csh ) then
        set HDPM_PATH="`pwd`"
    else if ( "$1" != "" && -e ${1}/setup.csh ) then
        set HDPM_PATH=${1}
    else
        echo "ERROR: Non-interactive usage requires one of the following lines:"
        echo "1. cd <HDPM_PATH>; source setup.csh"
        echo "2. source <HDPM_PATH>/setup.csh <HDPM_PATH>"
        exit 1
    endif
endif
echo $PATH | grep -q $HDPM_PATH
if ( $? != 0 ) then
    echo "Adding hdpm binary to PATH..."
    setenv PATH ${HDPM_PATH}/bin:$PATH
endif
