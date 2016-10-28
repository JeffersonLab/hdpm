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
if test $tag == c6; then
    mv ../$tag/opt .
    cd opt/rh/devtoolset-3/root; mv usr ../; cd ../
    rm -rf root; mkdir root; mv usr root
    cd root/usr/bin; rm -f ld; rm -f ../tmp
    ln -s ld.bfd ld; cd $cwd/$name-$tag
fi
mv ../$tag/home/gx/* .
if [[ $tag != u14 && $tag != u16 ]]; then
    mv ../$tag/usr/lib*/libblas.a cernlib/2005/lib/
    mv ../$tag/usr/lib*/liblapack.a cernlib/2005/lib/liblapack3.a
fi
#else
#    mv ../$tag/usr/lib/*/libblas.a cernlib/2005/lib/
#    mv ../$tag/usr/lib/*/liblapack.a cernlib/2005/lib/liblapack3.a
#fi
rm -rf ../$tag
cp -p ../../.id-deps-$tag .; cp -p ../../.log-sim-recon-$tag sim-recon/master/
commit=$(echo $(grep -i sim-recon sim-recon/master/*/success.hdpm) | sed -r 's/sim-recon-//g')
mkdir $cwd/$name-$tag-tmp
mv hdds $cwd/$name-$tag-tmp; mv sim-recon $cwd/$name-$tag-tmp
cp -p .id-deps-$tag $cwd/$name-$tag-tmp/hdds/master/; cp -p .id-deps-$tag $cwd/$name-$tag-tmp/sim-recon/master/
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
