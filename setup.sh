# bash: Add hdpm binary to PATH
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
