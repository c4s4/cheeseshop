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

To get help, type following command:

    $ cheeseshop -help
    Usage of build/cheeseshop:
      -auth="": Path to the authentication file
      -path="simple": The URL path
      -port=8000: The port CheeseShop is listening
      -root=".": The root directory for packages
      -shop="http://pypi.python.org/simple": Redirection when not found

The authentication file is made of lines with the username and the MD5 sum of the password separated with a space, such as (for user *foo* with password *bar*):

    foo 37b51d194a7513e45b56f6524f2d51f2

To compute MD5 sum for a given password, in order to fill the authentication file, you may type following command :

    $ echo -n bar | md5sum
    37b51d194a7513e45b56f6524f2d51f2  -

If no *-auth* option is set on command line, you won't have to authenticate to upload a package to *CheeseShop*.

*Enjoy!*
