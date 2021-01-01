Installation
============

Downloading the binaries directly
---------------------------------

The easiest way to download is to get a binary directly from the `releases <https://github.com/dhuan/giback/releases>`_ page.
Or also from your terminal:

.. code:: sh

    $ wget https://github.com/dhuan/giback/releases/download/%GIBACK_VERSION%/giback_%GIBACK_VERSION%_linux-386.zip
    $ unzip giback_%GIBACK_VERSION%_linux-386.zip
    $ chmod u+x giback
    $ ./giback --help

From source
-----------

.. code:: sh

    $ git clone
    $ cd giback
    $ make build
    $ ./bin/giback --help

Giback's execution logic is independent of which directory you're currently
located. You can make it easier to execute Giback without needing to be at the
source code folder:

.. code:: sh

   $ ln -s ./bin/giback ~/bin/giback
   $ giback --help



    
