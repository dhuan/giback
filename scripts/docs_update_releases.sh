ARG_REV="${1}"

if [ -z "${ARG_REV}" ]
then
    printf "Git revision missing!\n"

    exit 1
fi

if ! git rev-parse --quiet --verify "${ARG_REV}" > /dev/null
then
    printf "Invalid git revision: %s\n" "${ARG_REV}"

    exit 1
fi

TMP_CHANGELOG=$(mktemp -d)"/CHANGELOG.md"

git show "${ARG_REV}":CHANGELOG.md > $TMP_CHANGELOG

pandoc -s "${TMP_CHANGELOG}" -o docs/source/releases.rst 

