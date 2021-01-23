set -xe

git config --global user.email "giback_bot@noemail.com"
git config --global user.name "Giback Bot"

git checkout .

git remote set-url origin https://giback:$GH_KEY_DOCS@github.com/dhuan/giback.git

git fetch --tags --prune --unshallow

git checkout gh-pages

mv docs sphinx_docs

cp -r ./sphinx_docs/build/html ./docs

if [ ! -f './docs/.nojekyll' ]
then
    touch ./docs/.nojekyll
fi

LATEST_VERSION=$(git describe --tags --abbrev=0 origin/master)

FILES_TO_REPLACE=$(grep -rl '%GIBACK_VERSION%' docs | grep '\.html$')

for FILE in "$FILES_TO_REPLACE"
do
    sed -i "s/%GIBACK_VERSION%/$LATEST_VERSION/g" $FILE
done

git add ./docs

git commit -m "Update docs" || true

git push origin gh-pages
