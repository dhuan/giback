Installation
============

Downloading the binaries directly
---------------------------------

The easiest way to download is to get a binary directly from the `releases <https://github.com/dhuan/giback/releases>`_ page.
Or also from your terminal:

.. code:: sh

    $ wget https://github.com/dhuan/giback/releases/download/v0.1.0-beta-3/giback_v0.1.0-beta-3_linux-386.zip
    $ unzip giback_v0.1.0-beta-3_linux-386.zip
    $ chmod u+x giback
    $ ./giback --help

From source
-----------

.. code:: sh

    $ git clone
    $ cd giback
    $ make build
    $ ./giback --help

Giback's execution logic is independent of which directory you're currently
located. You can make it easier to execute Giback without needing to be at the
source code folder:

.. code:: sh

   $ ln -s $(pwd)/bin/giback ~/bin/giback
   $ giback --help



    
