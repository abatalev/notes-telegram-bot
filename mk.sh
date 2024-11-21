#!/bin/sh

CDIR=$(pwd)

if [ ! -d build ]; then
    mkdir build
fi

if [ ! -f build/prj2hash ]; then
    cd build
    # git clone https://github.com/abatalev/prj2hash prj2hash.git
    git clone http://localhost:3000/andrey/prj2hash prj2hash.git
    cd prj2hash.git
    ./build.sh
    cp prj2hash ../
    cd ${CDIR}/build
    rm -f -R prj2hash.git
    echo "### done build tools"
fi

cd $CDIR
if [ ! -f "build/gototcov" ]; then
    cd build
    git clone https://github.com/jonaz/gototcov gototcov.git
    cd gototcov.git
    go get golang.org/x/tools/cover
    go build .
    cp gototcov.git ../gototcov
    cd ${CDIR}/build
    rm -f -R gototcov.git
    echo "### done build tools"
fi

cd $CDIR
if [ ! -f "build/golangci-lint" ]; then
  echo "Install golangci-lint"
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b build/

  build/golangci-lint --version
fi

# cd $CDIR
# if [ ! -d examples ]; then
#     mkdir examples
# fi

# if [ ! -d examples/distcalc ]; then
#     cd examples
#     git clone http://localhost:3000/andrey/distcalc.git distcalc
# fi

# cd $CDIR
# if [ ! -d examples/demo-camunda ]; then
#     cd examples
#     git clone https://github.com/abatalev/demo-camunda.git demo-camunda
# fi

# cd $CDIR
# if [ ! -d examples/demo-stand ]; then
#     cd examples
#     git clone https://github.com/abatalev/demo-stand.git demo-stand
# fi

# cd $CDIR
# if [ ! -d examples/gitdiff2fly ]; then
#     cd examples
#     git clone https://github.com/abatalev/gitdiff2fly.git gitdiff2fly
# fi

# cd $CDIR
# if [ ! -d examples/prj2hash ]; then
#     cd examples
#     git clone http://localhost:3000/andrey/prj2hash.git prj2hash
# fi

# cd $CDIR
# if [ ! -d examples/gitstats ]; then
#     cd examples
#     git clone http://localhost:3000/andrey/gitstats.git gitstats
# fi

cd $CDIR
echo "### -[*]-[ Mod ]------------"
go mod tidy

cd ${CDIR}
echo "### -[*]-[ Lint ]------------"
build/golangci-lint run ./...
if [ "$?" != "0" ]; then
    echo "### aborted"
    exit 1
fi

echo "### -[*]-[ Test ]------------"
go test -v -coverpkg=./... -coverprofile=coverage.out ./... > /dev/null
if [ "$?" != "0" ]; then
    echo "### aborted"
    exit 1
fi

echo "### total coverage"
./build/gototcov -f coverage.out -limit 47
if [ "$?" != "0" ]; then
    echo "### open browser"
    go tool cover -html=coverage.out
    echo "### aborted"
    exit 1
fi

echo "### -[*]-[ Mutating tests ]------------"
#cd $CDIR/internal
~/go/bin/gremlins unleash
if [ "$?" != "0" ]; then
    echo "### aborted"
    exit 1
fi

cd $CDIR

echo "### -[*]-[ Build ]------------"
go build .
if [ "$?" != "0" ]; then
    echo "### aborted"
    exit 1
fi

function build_app_git(){
    GIT_HASH=$1
    if [ -f "./build/prj2hash" ]; then
        P2H_HASH=$(./build/prj2hash)
    fi
    go build -ldflags "-X main.gitHash=${GIT_HASH} -X main.p2hHash=${P2H_HASH}" -o notes-telegram .
}

echo "### build application with version"
build_app_git $(git rev-list -1 HEAD)

echo "### -[*]-[ Show version ]------------"
./notes-telegram -version
echo "### -[*]-[ The End ]------------"