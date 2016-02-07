# ship osx dist. to halldweb
target=/group/halld/www/halldweb/html/dist/
mkdir -p .pkgs; cd .pkgs
image=quay.io/nathansparks/sim-recon-osx
docker pull $image
docker rmi $(docker images -f "dangling=true" -q)
id=$(docker run -d $image)
docker export -o osx.tar $id; docker rm $id
mkdir -p osx; tar xf osx.tar -C osx; rm -f osx.tar
mv osx/home/sim-recon-* .; rm -rf osx
for file in $(ls); do
    if ! test -f $target/$file; then
        chgrp halld $file
        mv $file $target
    else
        rm $file
    fi
done
cd ../
