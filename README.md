# Giback - Backup your files to Git

## Usage

Giback relies on a configuration file at your home folder in order to know which files to backup and to what repositories:

```yml
# /home/user/giback.yml
units:
- id: my_backup
  repository: ssh://git@github.com:some_user/some_repository.git
  author_name: Someone
  author_email: someone@example.com
  commit_message: "Backing up!"
  ssh_key: /path/to/a/ssh_key/id_rsa
  files:
  - "/path/to/some/folder/*.txt"
  - "/path/to/another/folder/*.jpg"
  exclude:
  - "/path/to/some/folder/unwanted_file.txt"
- id: another_backup
  repository: ssh://git@github.com:some_user/some_repository.git
  author_name: Someone
  author_email: someone@example.com
  commit_message: "Backing up!"
  ssh_key: /path/to/a/ssh_key/id_rsa
  files:
  - "/path/to/some/folder/*.txt"
  - "/path/to/another/folder/*.jpg"
  exclude:
  - "/path/to/some/folder/unwanted_file.txt"
```

In addition, a "workspace" folder should exist at `/home/user/.giback`. A workspace is a folder Giback uses to clone and maintain the repositories your files will be backed up to.

With the these requirements addressed, you could then command Giback to backup all units:

```sh
giback all
```

If you wanted to backup a specific unit:

```sh
giback my_backup
```

Refer to the [documentation distributed with the source](./giback.txt) to understand more how to use Giback.

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
