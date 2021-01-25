set -x
set -e

CONTAINER_PROGRAM=podman
CONTAINER_BUILD_FLAGS="--cgroup-manager=cgroupfs"

run_tests() {
    "$CONTAINER_PROGRAM" exec giback_test bash /giback/test/giback_test_entrypoint.sh
}

if [ "$GIBACK_TESTS_USE_DOCKER" ]
then
    CONTAINER_PROGRAM=docker
fi

if [ "$GIBACK_TESTS_USE_DOCKER" ]
then
    CONTAINER_BUILD_FLAGS=""
fi

# If an image has been created already, this step will be skipped.
if [[ "$("$CONTAINER_PROGRAM" images -q giback_test 2> /dev/null)" == "" ]]
then
    echo "Git server image not found, generating one..."

    if [ "$GIBACK_TESTS_USE_DOCKER" ]
    then
        "$CONTAINER_PROGRAM" build -t giback_test -f dockerfile_test .
    else
        "$CONTAINER_PROGRAM" build "$CONTAINER_BUILD_FLAGS" -t giback_test -f dockerfile_test .
    fi
fi

CONTAINER_ID=$("$CONTAINER_PROGRAM" ps -aqf "name=giback_test")

# If Giback Test Container is already running, then just "run_tests()"
if [ ! -z "$CONTAINER_ID" ]
then
    run_tests

    exit 0
fi

# Otherwise, setup new ssh keys and start a Giback Test Container,
# and then "run_tests()"

# Remove previously generated SSH keys, if they exist.
rm -f ./test/tmp/id_rsa*

# Generate new SSH keys that will be used by Giback.
ssh-keygen -t rsa -f ./test/tmp/id_rsa -q -N ""

# Generate a new SSH key to assert that invalid keys are dealt with
ssh-keygen -t rsa -f ./test/tmp/id_rsa_invalid -q -N ""

"$CONTAINER_PROGRAM" run -d --rm --name giback_test -v $(pwd):/giback:z giback_test tail -f /dev/null

run_tests
