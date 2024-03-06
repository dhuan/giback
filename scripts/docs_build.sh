sh scripts/docs_update_releases.sh refs/remotes/origin/master
cd docs
make html
cp ./build/html/about.html ./build/html/index.html
