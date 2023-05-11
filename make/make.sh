#!/usr/bin/env bash

set -eox pipefail

function test-local() {
    set +x

    docker-compose -f make/docker-compose-dev.yml down --remove-orphans && \
    docker-compose -f make/docker-compose-dev.yml up --build -d

    export CONFIG_FILE=$(pwd)/config/config_dev.toml

    test
}

function test() {
    set +x
    set -e # exit if a command fails

    PKG_LIST=$(go list ./... | grep -v /vendor/ | grep -v /docs/swagger)
    COVERAGE_DIR="${COVERAGE_DIR:-.coverage}"

    echo "Run go mod tidy"
    go mod tidy -v

    echo "Run formating"
    go fmt $PKG_LIST

    echo "Run golangci-lint ..."
    golangci-lint run ./...

    # Remove the coverage files directory
    if [ -d "$COVERAGE_DIR" ]; then rm -Rf "$COVERAGE_DIR"; fi

    # run race by default

    RACEFLAG="-race"

    # script race on CI
    SKIPRACEFLAG=$1
    if [ -n "$SKIPRACEFLAG" ]; then
    RACEFLAG=""
    fi
    echo "RACEFLAG=$RACEFLAG"

    echo "Running tests and code coverage ..."

    # Create the coverage files directory
    mkdir -p "$COVERAGE_DIR";

    # Create a coverage file for each package
    # test minim coverage
    MINCOVERAGE=71

    # stop tests at first test fail
    TFAILMARKER="FAIL:"
    REGEXCOVERAGE="^coverage:"

    for package in $PKG_LIST; do
        go test $RACEFLAG -covermode=count -coverprofile "${COVERAGE_DIR}/${package##*/}.cov" "$package" -v -count=1 -p=1 | { IFS=''; while read -r line; do
            echo "$line"

            if [ -z "$line" ]; then
                continue
            fi

            if [ -z "${line##*$TFAILMARKER*}" ] ; then
                exit 10
            fi

            if [[ "${line}" =~ $REGEXCOVERAGE ]] ; then
                pcoverage=$(echo "$line"| grep "coverage" | sed -E "s/.*coverage: ([0-9]*\.[0-9]+)\% of statements/\1/g")

                if [ $(echo ${pcoverage%%.*}) -lt $MINCOVERAGE ] ; then
                    echo ""
                    echo "🚨 Test coverage of $package is $pcoverage%"
                    echo "FAIL: min coverage is $MINCOVERAGE%"
                    echo ""
                    exit 11
                else
                    echo ""
                    echo "🟢 Test coverage of $package is $pcoverage%"
                    echo ""
                fi
            fi
        done }
    done

    # Merge the coverage profile files
    echo 'mode: count' > "${COVERAGE_DIR}"/coverage.cov
    for fcov in "${COVERAGE_DIR}"/*.cov
    do
        if [ $fcov != "${COVERAGE_DIR}/coverage.cov" ]; then
            tail -q -n +2 $fcov >> "${COVERAGE_DIR}"/coverage.cov
        fi
    done


    # global code coverage
    pcoverage=$(go tool cover -func="${COVERAGE_DIR}"/coverage.cov | grep 'total:' | sed -E "s/^total:.*\(statements\)[[:space:]]*([0-9]*\.[0-9]+)\%.*/\1/g")
    echo "coverage: $pcoverage% of project"


    if [ $(echo ${pcoverage%%.*}) -lt $MINCOVERAGE ] ; then
        echo ""
        echo "🚨 Test coverage of project is $pcoverage%"
        echo "FAIL: min coverage is $MINCOVERAGE%"
        echo ""
        exit 12
    else
        echo ""
        echo "🟢 Test coverage of project is $pcoverage%"
        echo ""
    fi
}