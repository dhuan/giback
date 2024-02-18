# Giback - Backup your files to Git

[![Build Status](https://github.com/dhuan/giback/actions/workflows/go.yml/badge.svg)](https://github.com/dhuan/giback/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dhuan/giback)](https://goreportcard.com/report/github.com/dhuan/giback)

Giback is a command-line utility for easily backing up your files to git repositories. After setting up a configuration file listing your desired files, and the respective target repository, you can then at any time command giback to backup your files.

```yml
$ giback all
Running unit 'my_backup'.
Pulling git changes.
Identifying files...
/home/my_user/Documents/personal_notes.txt
/home/my_user/Documents/work.org
/home/my_user/photos/trip.jpg
Files copied.
Committing: Backing up with Giback!
Pushing...
Done!
```

## Resources

- Manual: https://dhuan.github.io/giback
- Downloads: https://github.com/dhuan/giback/releases

## Installation from source

The only requirement to compile is to have Go installed.

```sh
git clone git@github.com:dhuan/giback.git

cd giback

make build
```

Once compilation has finished successfully, a `giback` executable should then be available at the current dir.

## Known Issues

- Only has been tested in Linux. Probably will not work in Windows.
- Works only with files. Backing up directories still not supported.
