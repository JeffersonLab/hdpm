# Build and ship hdpm distribution to JLab CUE
name=hdpm
ver=$(hdpm version | awk '{print $3}')

target=/group/halld/www/halldweb/html/dist/${name}
mkdir -p ${target}

target2=/group/halld/Software/${name}
mkdir -p ${target2}

# Build for 64-bit Linux
GOOS=linux GOARCH=amd64 go build -ldflags '-s'

# Pack and ship to halldweb
dir=${name}-${ver}
mkdir -p ${dir}/bin
mv hdpm ${dir}/bin
cp -p README.md setup.*sh ${dir}
tarfile=${dir}.linux.tar.gz
tar czf ${tarfile} ${dir}
mv ${tarfile} ${target}
chgrp -R halld ${target}

# If it's a release, move it to halld software folder
if [[ ! -d "${target2}/${ver}" && "${ver}" != "dev" ]]; then
	mv ${dir} ${ver}
	mv ${ver} ${target2}
	chgrp -R halld ${target2}
else
	rm -rf ${dir}
fi

# Build for 64-bit macOS
GOOS=darwin GOARCH=amd64 go build

# Pack and ship to halldweb
dir=${name}-${ver}
mkdir -p ${dir}/bin
mv hdpm ${dir}/bin
cp -p README.md setup.*sh ${dir}
tarfile=${dir}.macOS.tar.gz
tar czf ${tarfile} ${dir}
rm -rf ${dir}
mv ${tarfile} ${target}
chgrp -R halld ${target}
