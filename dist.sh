# Build and ship hdpm distribution to JLab CUE
name=hdpm
ver=$(hdpm version | awk '{print $3}')

target=/group/halld/www/halldweb/html/dist/${name}
mkdir -p ${target}

target2=/group/halld/Software/${name}
mkdir -p ${target2}

# Build for 64-bit Linux
GOOS=linux GOARCH=amd64 go build

# Pack and ship to targets
dir=${name}-${ver}
mkdir -p ${dir}/bin
mv hdpm ${dir}/bin
cp -p README.md ${dir}
tarfile=${dir}.linux.tar.gz
tar czf ${tarfile} ${dir}
mv ${tarfile} ${target}
chgrp -R halld ${target}

if [[ ! -f "${target2}/${ver}" ]]; then
	mv ${dir} ${ver}
	mv ${ver} ${target2}
	chgrp -R halld ${target2}
else
	rm -rf ${dir}
fi

# Build for 64-bit macOS
GOOS=darwin GOARCH=amd64 go build

# Pack and ship to target
dir=${name}-${ver}
mkdir -p ${dir}/bin
mv hdpm ${dir}/bin
cp -p README.md ${dir}
tarfile=${dir}.macOS.tar.gz
tar czf ${tarfile} ${dir}
rm -rf ${dir}
mv ${tarfile} ${target}
chgrp -R halld ${target}
