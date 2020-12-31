mkdir -p docs/source/_themes

if [ ! -d "docs/source/_themes/sphinx_rtd_theme_source" ]
then
    git clone --depth 1 --branch 0.5.0 https://github.com/readthedocs/sphinx_rtd_theme.git docs/source/_themes/sphinx_rtd_theme_source
    cp -r docs/source/_themes/sphinx_rtd_theme_source/sphinx_rtd_theme docs/source/_themes/sphinx_rtd_theme
fi

cd docs && make html
