# pack for distribution
tag=$1
target=/group/halld/www/halldweb/html/dist/
name=sim-recon
id_deps=`cat .id-deps-$tag`
mkdir -p .pkgs; cd .pkgs
cwd=$(pwd)
id=`docker run -d sim-recon:$tag`
docker export -o $tag.tar $id; docker rm $id
mkdir $tag; tar xf $tag.tar -C $tag; chmod -R u+w $tag; rm -f $tag.tar
mkdir $name-$tag; cd $name-$tag
mv ../$tag/home/gx/.hdpm . && mv .hdpm/env .
rm -rf ../$tag/home/gx/.[!.]*
mv ../$tag/home/gx/* .
if [[ $tag != u16 ]]; then
    mv ../$tag/usr/lib*/libblas.a cernlib/2005/lib/
    mv ../$tag/usr/lib*/liblapack.a cernlib/2005/lib/liblapack3.a
fi
rm -rf ../$tag
cp -p ../../.id-deps-$tag .; cp -p ../../.log-sim-recon-$tag sim-recon/master/
commit=$(echo $(grep -i sim-recon sim-recon/master/*/success.hdpm) | sed -r 's/sim-recon-//g')
mkdir $cwd/$name-$tag-tmp
mv hdds sim-recon hdgeant4 gluex_root_analysis $cwd/$name-$tag-tmp
cd $cwd
mv $name-$tag $name-deps-$tag; mv $name-$tag-tmp $name-$tag
if ! test -f $target/$name-deps-$id_deps-$tag.tar.gz; then
    tar czf $name-deps-$id_deps-$tag.tar.gz $name-deps-$tag
    chgrp halld $name-deps-$id_deps-$tag.tar.gz
    mv $name-deps-$id_deps-$tag.tar.gz $target
fi
if ! test -f $target/$name-$commit-$id_deps-$tag.tar.gz; then
    tar czf $name-$commit-$id_deps-$tag.tar.gz $name-$tag
    chgrp halld $name-$commit-$id_deps-$tag.tar.gz
    mv $name-$commit-$id_deps-$tag.tar.gz $target
fi
rm -rf $name-$tag; rm -rf $name-deps-$tag
cd ../
