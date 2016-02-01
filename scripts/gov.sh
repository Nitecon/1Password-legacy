#!/bin/bash
#####################################################################
## gov.sh is a simple application for dealing with go vendoring    ##
## once go has full support + tooling for vendoring then this will ##
## dissappear, until then as a utility for me to get updated       ##
## vendor packages we will using this to keep vendor libs updated  ##
## License: See LICENSE in root                                    ##
## Usage: ./bin/gov.sh --update                                    ##
## Please note you should be running this from the repository root ##
#####################################################################
GOSRC="$GOPATH/src/"
CURPATH=`pwd`

if [ -z "$GOPATH" ]; then
    echo "Need to set GOPATH"
    exit 1
fi

if ! [[ $GOPATH == */ ]]; then
	GOPATH=$GOPATH/
fi

if ! [[ $PWD == ${GOPATH}src/* ]]; then
	echo Must be in '$GOPATH/src/<projname>' to use vendoring
	exit 1
fi

if [[ $PWD == *"vendor"* ]]; then
	echo "Cannot do vendoring inside vendor"
	exit 1
fi

GOLEN=${#GOSRC}
PACKLEN=${#CURPATH}
MaxLEN=`expr $PACKLEN - $GOLEN`
CURPACKAGE=${CURPATH:GOLEN:MaxLEN}
DEPLIST=`go list -f '{{join .Deps "\n"}}' |  xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}'|grep -v $CURPACKAGE`
for PNAME in $DEPLIST; do
	IFS='/' read -ra i <<< "$PNAME"
	DEPPACKAGE="${i[0]}/${i[1]}/${i[2]}"
	if [ -d $GOPATH/src/$DEPPACKAGE ]; then
		#Package is available so lets see if it exists and if not pull it / sync
		if [ ! -d $CURPATH/vendor/$DEPPACKAGE ]; then
			mkdir -p $CURPATH/vendor/$DEPPACKAGE
			rsync -a --exclude=.git $GOPATH/src/$DEPPACKAGE/ $CURPATH/vendor/$DEPPACKAGE/
		fi
	else
		go get -u $DEPPACKAGE
		mkdir -p $CURPATH/vendor/$DEPPACKAGE
		rsync -a --exclude=.git $GOPATH/src/$DEPPACKAGE/ $CURPATH/vendor/$DEPPACKAGE/
	fi
done

if [ "$1" == "--update" ]; then
	if [ -d vendor ]; then
		for dist in `ls $CURPATH/vendor/`; do
			for owner in `ls $CURPATH/vendor/$dist/`; do
				if [ "$dist" == "gopkg.in" ]; then
					echo "Upgrading $dist/$owner"
					go get -u $dist/$owner
					rsync -a --del --exclude=.git $GOPATH/src/$dist/$owner/ $CURPATH/vendor/$dist/$owner/
					#cd $CURPATH/vendor/$dist/$owner && git pull
				else
					for pack in `ls $CURPATH/vendor/$dist/$owner/`; do
						echo "Upgrading $dist/$owner/$pack"
						go get -u $dist/$owner/$pack
						rsync -a --del --exclude=.git $GOPATH/src/$dist/$owner/$pack/ $CURPATH/vendor/$dist/$owner/$pack/
						#cd $CURPATH/vendor/$dist/$owner/$pack && git pull
					done
				fi
			done
		done
	else
		echo "This is not a vendored go app."
	fi
fi