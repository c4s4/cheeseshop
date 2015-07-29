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

To run the server, type on command line :

    cheeseshop -port 8000 -path simple -root . -shop http://pypi.python.org

Options on command line :

- *-port* to set the port the server is litening (default to *8000*).
- *-path* to set the URL path (defaults to *simple*).
- *-root* to set the directory where packages are living (defaults to current directory).
- *-shop* to set the URL of the shop for packages that are not found.
- *-auth* to set the path to the authentication file

The server outputs logs on the terminal. To get help on the console, type `cheeseshop -help`.

The authentication file is made of lines with the username and the MD5 sum of the password separated with a space, such as (for user *foo* with password *bar*):

    foo 37b51d194a7513e45b56f6524f2d51f2

To compute MD5 sum for a given password, in order to fill the authentication file, you may type following command :

    $ echo -n bar | md5sum
    37b51d194a7513e45b56f6524f2d51f2  -

If no *-auth* option is set on command line, you won't have to authenticate to upload a package to *CheeseShop*.

*Enjoy!*
