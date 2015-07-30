CheeseShop
==========

- Project: <https://github.com/c4s4/cheeseshop>
- Downloads: <https://github.com/c4s4/cheeseshop/releases>

CheeseShop is a Python package repository. This is a local version of the well-known <http://pypi.python.org>. This is useful for enterprise users that need to share private Python libraries among developers.

To tell PIP where is your private CheeseShop, you must edit you *~/.pip/pip.conf* file:

    [global]
    index-url = http://my.shop.host/simple
    trusted-host = my.shop.host

Where *my.shop.host* is the hostname of the machine running CheeseShop. PIP will call your CheeseShop to get packages. If CheeseShop doesn't host this package it will redirect PIP to standard Pypi.

To tell *setup.py* where to upload your package, you must edit file *~/.pypirc*:

    [distutils]
    index-servers =
        pypi
        cheeseshop
    
    [pypi]
    username: ***
    password: ***
    
    [cheeseshop]
    username: spam
    password: foo
    repository: http://my.shop.host/simple/

*setup.py* will call your CheeseShop if you name it on command line:

    $ python setup.py sdist upload -r cheeseshop

Where `-r cheeseshop` is the option that indicates the connection you want to use. There must be a corresponding entry in your *~/.pypirc* configuration file. Don't forget to add *cheeseshop* in the *index-server* list at the beginning of the file.

CheeseShop is able to run on HTTP and/or HTTPS and performs basic authentication if necessary.

Installation
------------

Download binary archive at <https://github.com/c4s4/cheeseshop/releases>, unzip it and copy the binary executable for your platform (named *cheeseshop-system-platform* in the *bin* directory) somewhere in your *PATH* and rename it *cheeseshop*. This executable doesn't need any dependency or virtual machine to run.

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

    # The root directory for packages
    root:  /home/cheeseshop
    # Path to the server certificate
    cert:  /etc/ssl/certs/cheeseshop-cert.pem
    # Path to the server key
    key:   /etc/ssl/private/cheeseshop-key.pem
    # The HTTP port CheeseShop is listening
    http:  80
    # The HTTPS port CheeseShop is listening 
    https: 443
    # The URL path
    path:  simple
    # Redirection when not found
    shop:  http://pypi.python.org/simple
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

There is a sample configuration file in *etc* directory of the archive.

Of course, you must create an empty directory for the repository. Ensure that the user running CheeseShop has a right to write in this directory.

To disable HTTP or HTTPS, you must set port to *0*. If HTTPS is disabled, you don't have to set certificate and key paths. To disable basic authentication, you must set auth to `~` (which means none in YAML).

To generate a key, you can use openssl as follows:

    openssl genrsa -out cheeseshop-key.pem 2048

To generate au self signed certificate, you can type:

    openssl req -new -x509 -key cheeseshop-key.pem -out cheeseshop-cert.pem -days 3650

This command will ask you many fields, but the only that is necessary is the *FQDN* which is the name of the machine that is running CheeseShop.

You should copy the certificate in directory */etc/ssl/certs* and the key in */etc/ssl/private*.

Service
-------

To install CheeseShop as a System V service, edit sample init script in *etc/cheeseshop.init* file. You should edit *SCRIPT* variable to set the path to the *cheeseshop* command. Then copy this file as */etc/init.d/cheeseshop*.

You must also edit configuration file *etc/cheeseshop.yml* to set the repository location in the *root* variable. Copy this file in */etc/cheeseshop.yml* location.

You can then start the service with:

    $ sudo service cheeseshop start

And stop it with:

    $ sudo service cheeseshop stop

You can view the logs in */var/log/cheeseshop.log* file.

To start the service at boot, you should type:

    $ sudo update-rc.d cheeseshop defaults

And to disable start at boot:

    $ sudo update-rc.d -f cheeseshop remove

Build CheeseShop
----------------

To build CheeseShop, you must install [Goyaml](http://github.com/go-yaml/yaml) and [GOX](http://github.com/mitchellh/gox) with following commands:

    $ go get gopkg.in/yaml.v2
    $ go get github.com/mitchellh/gox
    $ gox -build-toolchain

Then you can use the make file to build the binary version for your platform:

    $ make build

To build binaries for all platforms, type:

    $ make compile

*Enjoy!*
