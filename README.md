CheeseShop
==========

- Project: <https://github.com/c4s4/cheeseshop>
- Downloads: <https://github.com/c4s4/cheeseshop/releases>

CheeseShop is a Python package repository.

Installation
------------

Download binary archive at <https://github.com/c4s4/cheeseshop/releases>, unzip it and copy the binary executable for your platform (named *cheeseshop*) somewhere in yout *PATH* and rename it *cheeseshop*. This executable doesn't need any dependency or virtual machine to run.

There are binaries for following platforms:

- Linux 386, amd64 and arm.
- FreeBSD 386, amd64 and arm.
- NetBSD 386, amd64 and arm.
- OpenBSD 386 and amd64.
- Darwin 386 and amd64.
- Windows 386 and amd64.
- Plan9 386 and amd64.

Usage
-----

To run the server, type on command line :

    cheeseshop -port 8000 -path simple -root . -shop http://pypi.python.org

Options on command line :

- *-port* to set the port the server is litening (default to *8000*).
- *-path* to set the URL path (defaults to *simple*).
- *-root* to set the directory where packages are living (defaults to current directory).
- *-shop* to set the URL of the shop for packages that are not found.

The server outputs logs on the terminal. To get help on the console, type `cheeseshop -help`.

*Enjoy!*
