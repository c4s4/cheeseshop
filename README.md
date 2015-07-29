CheeseShop
==========

- Project: <https://github.com/c4s4/cheeseshop>
- Downloads: <https://github.com/c4s4/cheeseshop/releases>

CheeseShop is a Python package repository.

Installation
------------

Download binary archive at <https://github.com/c4s4/cheeseshop/releases>, unzip it and copy the binary executable for your platform (named *cheeseshop-system-platform* in the *bin* directory) somewhere in yout *PATH* and rename it *cheeseshop*. This executable doesn't need any dependency or virtual machine to run.

There are binaries for following platforms:

- Linux 386, amd64 and arm.
- FreeBSD 386, amd64 and arm.
- NetBSD 386, amd64 and arm.
- OpenBSD 386 and amd64.
- Darwin (MacOSX) 386 and amd64.
- Windows 386 and amd64.

Usage
-----

To run CheeseShop, type on command line:

    $ cheeseshop

It will look for a configuration file at following locations:

- *~/.cheeseshop.yml*
- */etc/cheeseshop.yml*

You may also pass the path to the configuration file on the command line:

    $ cheeseshop /path/to/cheeseshop.yml

This configuration file should look like this:

    # The port CheeseShop is listening
    port: 8000
    # The URL path
    path: simple
    # The root directory for packages
    root: repo
    # Redirection when not found
    shop: http://pypi.python.org/simple
    # List of users and their MD5 hashed password
    # To get MD5 sum for password foo, type 'echo -n foo | md5sum'
    # To disable auth when uploading packages, set auth to ~
    auth:
        spam: acbd18db4cc2f85cedef654fccc4a4d8
        eggs: 37b51d194a7513e45b56f6524f2d51f2

To compute MD5 sum for a given password, in order to fill the authentication file, you may type following command :

    $ echo -n foo | md5sum
    acbd18db4cc2f85cedef654fccc4a4d8  -
    $ echo -n bar | md5sum
    37b51d194a7513e45b56f6524f2d51f2  -

*Enjoy!*
