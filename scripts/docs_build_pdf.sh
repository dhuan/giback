set -xe

if [ ! -z "$GITHUB_ACTION" ]
then
    sudo apt install --yes latexmk python3-sphinx texlive-formats-extra
    echo "master_doc = 'index'" >> docs/source/conf.py
fi

cd docs && make latexpdf
