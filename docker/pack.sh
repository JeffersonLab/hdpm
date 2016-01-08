# pack for distribution
tag=$1
target=/group/halld/www/halldweb/html/dist/
name=sim-recon
id_deps=`cat .id-deps-$tag`
mkdir -p .pkgs; cd .pkgs
cwd=`pwd`
id=`docker run -d sim-recon:$tag`
docker export -o $tag.tar $id; docker rm $id
mkdir $tag; tar xf $tag.tar -C $tag; chmod -R u+w $tag; rm -f $tag.tar
mkdir $name-$tag; cd $name-$tag
if test $tag == c6; then
    mv ../$tag/opt .
fi
mv ../$tag/home/hdpm/settings .; mv ../$tag/home/hdpm/pkgs/* .; rm -rf ../$tag
cp -p ../../.id-deps-$tag .; cp -p ../../.log-sim-recon-$tag sim-recon/
commit=$(echo $(grep -i sim-recon sim-recon/*/success.hdpm) | sed -r 's/sim-recon-//g')
mkdir $cwd/$name-$tag-tmp
mv hdds $cwd/$name-$tag-tmp; mv sim-recon $cwd/$name-$tag-tmp
cp -p .id-deps-$tag $cwd/$name-$tag-tmp/hdds/; cp -p .id-deps-$tag $cwd/$name-$tag-tmp/sim-recon/
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
