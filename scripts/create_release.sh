GIBACK_VERSION=$(echo $GITHUB_REF | cut -d '/' -f 3)

export GH_TOKEN=${GH_KEY}

echo $GIBACK_VERSION

gh release create "$GIBACK_VERSION" -t "$GIBACK_VERSION"

RELEASE_FILES=$(ls ./release_downloads/*.zip)

for RELEASE_FILE in $RELEASE_FILES
do
    gh release upload "$GIBACK_VERSION" "$RELEASE_FILE"
done
