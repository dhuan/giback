set -xe

git config --global user.email "giback_bot@noemail.com"
git config --global user.name "Giback Bot"

cp -r ./docs/build/html ./public

git checkout .

git remote set-url origin https://giback:$GH_KEY_DOCS@github.com/dhuan/giback.git

git fetch

git checkout gh-pages

git add ./public

git commit -m "Update docs"

git push origin gh-pages
