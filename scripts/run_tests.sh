set -x
set -e

# Remove previously generated SSH keys, if they exist.
rm -f ./test/tmp/id_rsa*

# Generate new SSH keys that will be used by Giback.
ssh-keygen -t rsa -f ./test/tmp/id_rsa -q -N ""

# Generate a new SSH key to assert that invalid keys are dealt with
ssh-keygen -t rsa -f ./test/tmp/id_rsa_invalid -q -N ""

CONTAINER_PROGRAM=podman
CONTAINER_BUILD_FLAGS="--cgroup-manager=cgroupfs"

if [ "$GIBACK_TESTS_USE_DOCKER" ]
then
    CONTAINER_PROGRAM=docker
fi

if [ "$GIBACK_TESTS_USE_DOCKER" ]
then
    CONTAINER_BUILD_FLAGS=""
fi

# Build a simple image containing a git server where backup repositories will exist that Giback will push to.
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

# If a container is already running of that image we just created, let's remove it.
CONTAINER_ID=$("$CONTAINER_PROGRAM" ps -aqf "name=giback_test")

if [ ! -z "$CONTAINER_ID" ]
then
    echo "Removing container $CONTAINER_ID."

    "$CONTAINER_PROGRAM" stop --ignore $CONTAINER_ID
    "$CONTAINER_PROGRAM" rm -fi $CONTAINER_ID
fi

# Run the git server container we just created.
"$CONTAINER_PROGRAM" run --rm --name giback_test -v $(pwd):/giback:z giback_test
