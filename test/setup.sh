# Remove previously generated SSH keys, if they exist.
rm -f ./test/tmp/id_rsa*

# Generate new SSH keys that will be used by Giback.
ssh-keygen -t rsa -f ./test/tmp/id_rsa -q -N ""

# Build a simple image containing a git server where backup repositories will exist that Giback will push to.
# If an image has been created already, this step will be skipped.
if ! podman image exists giback_test_git_server
then
    echo "Git server image not found, generating one..."

    podman build --cgroup-manager=cgroupfs -t giback_test_git_server -f test/dockerfile_git
fi

# If a container is already running of that image we just created, let's remove it.
CONTAINER_ID=$(podman ps -aqf "name=giback_test_git_server")

if [ ! -z "$CONTAINER_ID" ]
then
    echo "Removing container $CONTAINER_ID."

    podman stop --ignore $CONTAINER_ID
    podman rm -fi $CONTAINER_ID
fi

# Run the git server container we just created.
podman run -d --rm -ti --name giback_test_git_server -p 2222:22 -v $(pwd)/test:/giback_test:z giback_test_git_server
