# copy tarfiles into container and publish
docker-machine start default
yes | docker-machine regenerate-certs default
eval $(docker-machine env --shell=bash default)
mkdir -p .docker; cd .docker
c=$(cat ../.commit)
id_deps=$(cat ../.id-deps-osx)
file_deps=sim-recon-deps-$id_deps-osx.tar.gz
file=sim-recon-$c-$id_deps-osx.tar.gz
cp -p ../.pkgs/$file_deps .; cp -p ../.pkgs/$file .
echo "FROM alpine" >> Dockerfile
echo "COPY $file_deps /home/$file_deps" >> Dockerfile
echo "COPY $file /home/$file" >> Dockerfile
echo "CMD sh" >> Dockerfile
image=quay.io/nathansparks/sim-recon-osx
docker build -t $image .
rm $file_deps $file Dockerfile
docker push $image
docker rmi $(docker images -f "dangling=true" -q)
docker-machine kill default
cd ../
