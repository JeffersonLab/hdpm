# Download dependencies and move to halldweb
# $tarfile is for docker builds (see Dockerfile-deps)
tarfile=gxsrc.tar.gz
target=/group/halld/www/halldweb/html/dist/hdpm
export GLUEX_TOP=$PWD/gx
mkdir $GLUEX_TOP
# Fetch all deps
#hdpm fetch -d cernlib amptools evio jana rcdb geant4
# Fetch all deps except ROOT
hdpm fetch amptools ccdb cernlib evio geant4 jana rcdb xerces-c
tar czf $tarfile gx
chgrp halld $tarfile 
mv $tarfile $target/
rm -rf $GLUEX_TOP
