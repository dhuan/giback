Giback: Backup your files to git
================================

Giback is a command-line utility for easily backing up your files to git
repositories. After setting up a configuration file listing your desired files,
and the respective target repository, you can then at any time command giback
to backup your files.


.. code-block:: sh

    $ giback all
    Running unit 'my_backup'.
    Pulling git changes.
    Identifying files...
    /home/john_doe/Documents/personal_notes.txt
    /home/john_doe/Documents/work.org
    /home/john_doe/photos/trip.jpg
    Files copied.
    Committing: Backing up with Giback!
    Pushing...
    Done!

Get it now
----------

If you have Go installed, you can quickly install Giback:

.. code:: shell

    $ go get github.com/dhuan/giback

The ``giback`` executable utility should then be available for you. If not,
check if the ``GOBIN`` or ``GOPATH`` environment variables are set properly.

.. Note::

    For other installation methods, `check the installation manual. </installation>`_

Resources
---------

- Homepage/Manual: https://dhuan.github.io/giback
- Downloads: https://github.com/dhuan/giback/releases
- Source code/Git URL: git://github.com/dhuan/giback.git

