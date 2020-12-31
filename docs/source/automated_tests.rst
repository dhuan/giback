Automated Tests
===============

.. Important::

    Before attempting to run the automated tests in your computer, make sure to
    have either `Docker <https://docker.com/>`_ or
    `Podman <https://podman.io/>`_ installed.

Tests are executed in a containerized environment. As the test script starts,
an image containing a git server and the Giback executable is created, so that
it's not necessary for you to have a Git server setup running in your host
machine.

Both Podman and Docker are supported. By default, Podman will be used, when the
the test command is issued:

.. code:: sh

    $ make tests

If you wanted instead to use Docker:

.. code:: sh

    $ GIBACK_TESTS_USE_DOCKER=true make tests
