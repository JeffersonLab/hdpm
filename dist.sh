# Build and ship hdpm distribution to JLab CUE
target=/group/halld/www/halldweb/html/dist/
target2=/group/halld/Software/
name=hdpm
ver=$(hdpm version | awk '{print $3}')

# Build for 64-bit Linux
GOOS=linux GOARCH=amd64 go build

# Pack and ship to targets
dir=${name}-${ver}
mkdir -p ${dir}/bin
mv hdpm ${dir}/bin
cp -p README.md ${dir}
tarfile=${dir}.linux.tar.gz
tar czf ${tarfile} ${dir}
chgrp halld ${tarfile}
mv ${tarfile} ${target}

if [[ ! -f "${target2}/${name}/${ver}" ]]; then
	mkdir -p ${target2}/${name}
	mv ${dir} ${ver}
	mv ${ver} ${target2}/${name}
	chgrp -R halld ${target2}/${name}
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
chgrp halld ${tarfile}
rm -rf ${dir}
mv ${tarfile} ${target}
