set -ex

GIBACK_VERSION=$(echo $GITHUB_REF | cut -d '/' -f 3)

TARGETS=(
    "linux,386"
    "linux,amd64"
    "linux,arm"
    "linux,arm64"
    "darwin,amd64"
)

rm -rf ./release_downloads

mkdir ./release_downloads

make docs_build_pdf

for TARGET in "${TARGETS[@]}"
do
    GOOS=$(echo $TARGET | cut -d "," -f 1)
    GOARCH=$(echo $TARGET | cut -d "," -f 2)

    TARGET_NAME="${GOOS}-${GOARCH}"

    printf "Generating build for ${TARGET_NAME}\n"

    TARGET_PATH="./release_downloads/$TARGET_NAME"

    mkdir $TARGET_PATH

    cp ./README.md $TARGET_PATH/.
    cp ./LICENSE $TARGET_PATH/.
    cp ./docs/build/latex/*.pdf $TARGET_PATH/.

    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X main.version=$GIBACK_VERSION" -o $TARGET_PATH/giback
done

TARGET_FOLDERS=$(ls ./release_downloads)

for TARGET_FOLDER in ${TARGET_FOLDERS[@]}
do
    zip "./release_downloads/giback_${GIBACK_VERSION}_${TARGET_FOLDER}.zip" -j ./release_downloads/"${TARGET_FOLDER}"/*
done
