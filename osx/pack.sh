# pack for distribution
tag=osx
top=../pkgs/osx
target=`pwd`/.pkgs
name=sim-recon
id_deps=`cat .id-deps-$tag`
mkdir -p .pkgs; cd .pkgs
cwd=`pwd`
cp -pr ../$top $tag
mkdir $name-$tag; cd $name-$tag
cp -pr ../../../settings .; mv ../$tag/* .; rm -rf ../$tag
cp -p ../../.id-deps-$tag .
mv env-setup/osx.sh env-setup/master.sh
mv env-setup/osx.csh env-setup/master.csh
commit=$(echo $(grep -i sim-recon sim-recon/master/*/success.hdpm) | sed -E 's/sim-recon-//g')
mkdir $cwd/$name-$tag-tmp
mv hdds $cwd/$name-$tag-tmp; mv sim-recon $cwd/$name-$tag-tmp
cp -p .id-deps-$tag $cwd/$name-$tag-tmp/hdds/master/; cp -p .id-deps-$tag $cwd/$name-$tag-tmp/sim-recon/master/
cd $cwd
mv $name-$tag $name-deps-$tag; mv $name-$tag-tmp $name-$tag
if ! test -f $target/$name-deps-$id_deps-$tag.tar.gz; then
    tar czf $name-deps-$id_deps-$tag.tar.gz $name-deps-$tag
    mv $name-deps-$id_deps-$tag.tar.gz $target
fi
if ! test -f $target/$name-$commit-$id_deps-$tag.tar.gz; then
    tar czf $name-$commit-$id_deps-$tag.tar.gz $name-$tag
    mv $name-$commit-$id_deps-$tag.tar.gz $target
fi
rm -rf $name-$tag; rm -rf $name-deps-$tag
cd ../
