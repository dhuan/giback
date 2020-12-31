set -xe

git config --global user.email "giback_bot@noemail.com"
git config --global user.name "Giback Bot"

git checkout .

git remote set-url origin https://giback:$GH_KEY_DOCS@github.com/dhuan/giback.git

git fetch

git checkout gh-pages

cp -r ./docs/build/html ./public

FILES_TO_REPLACE=$(grep -rl '%GIBACK_VERSION%' docs | grep '\.html$')

for FILE in "$FILES_TO_REPLACE"
do
    sed -i 's/%GIBACK_VERSION%/1.0/g' $FILE
done

git add ./public

git commit -m "Update docs"

git push origin gh-pages
