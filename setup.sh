# bash: hdpm setup
# usage: source setup.sh
SCRIPT_PATH=${BASH_SOURCE[0]}
if [[ -h $SCRIPT_PATH ]]; then
    SCRIPT_PATH=$(readlink $SCRIPT_PATH)
fi
HDPM_PATH=$(cd $(dirname $SCRIPT_PATH); pwd)
echo $PATH | grep -q $HDPM_PATH
if [[ $? -ne 0 ]]; then
    echo "Adding hdpm binary to PATH..."
    export PATH=$HDPM_PATH/bin:$PATH
fi
if [[ -z "$GLUEX_TOP" && "$(basename $HDPM_PATH)" == .hdpm ]]; then
    echo "Setting GLUEX_TOP..."
    export GLUEX_TOP=$(dirname $HDPM_PATH)
    echo GLUEX_TOP=$GLUEX_TOP
    if [[ -z "$HALLD_MY" ]]; then
        echo "Setting HALLD_MY..."
        export HALLD_MY=$GLUEX_TOP/plugins
        echo HALLD_MY=$HALLD_MY
    fi
fi
