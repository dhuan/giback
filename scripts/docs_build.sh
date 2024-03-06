sh scripts/docs_update_releases.sh v0.1.1
cd docs
make html
cp ./build/html/about.html ./build/html/index.html
