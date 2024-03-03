Usage and Configuration
=======================

Giback needs a configuration file describing what files you want to backup and
to which repository they should be pushed to. Below is an example of a
configuration file that defines two Backup Units, ``notes`` and ``photos``:

.. code-block:: yaml

    units:
    - id: notes
      repository: ssh://git@github.com:my_user/my_personal_notes.git
      author_name: John Doe
      author_email: johndoe@example.com
      commit_message: 'Backing up with Giback.'
      files:
      - "/home/john_doe/Documents/*.txt"
      - "/home/john_doe/Documents/*.org"
      exclude:
      - "/home/john_doe/Documents/unimportant_notes.txt"
    - id: photos
      repository: ssh://git@github.com:my_user/photos.git
      author_name: John Doe
      author_email: johndoe@example.com
      commit_message: 'Backing up with Giback.'
      files:
      - "/home/john_doe/Documents/*.jpg"
      - "/home/john_doe/Documents/*.png"

Given these two Backup Units defined above, we can just command Giback at any
time to backup everything:

.. code-block:: sh

    $ giback all

In case we wanted to backup a single Backup Unit, say ``notes``, instead of
all:

.. code-block:: sh

    $ giback notes

Backup Units' Configuration Parameters
--------------------------------------

id
^^

An ID that identifies a Backup Unit.

repository
^^^^^^^^^^

The Git repository to which your backup files will be pushed.

author_name
^^^^^^^^^^^

The author's name to which commits will be associated.

author_email
^^^^^^^^^^^^

The author's email.

commit_message
^^^^^^^^^^^^^^

A commit message.

files
^^^^^

A list of files that will pushed to your backup repository. Absolute
paths should be used. Glob Patterns are supported when you want to backup
multiple files in the same directory or that have the same extension.

exclude
^^^^^^^

Like `files`, this takes a list of files, only that those files
will be ignored and not pushed to your backup repository.

ssh_key
^^^^^^^

Absolute path to a SSH key. If this paramter is not defined, your user's
default SSH key will be used (like git normally does).

Command Line Options
--------------------

-w
    Path to a workspace. Use it if you don't want your workspace to be
    "/home/user/.giback". All repositories will be kept there.

-c
    Path to a configuration file. Use it if you want a configuration file
    that's not located at the default path "/home/user/.giback.yml".

-v
    By default giback prints useful log messages describing its operations as
    it runs. It doesn't however show the output of the git command it runs.
    With this flag all git output will be visible. This is specifically useful
    for debugging.
