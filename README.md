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

Configuration
-------------

The configuration file should look like this:

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

There is a sample configuration in file *etc/cheeseshop.yml* of the archive.

### root

The root directory is where live the Python packages. Under this root there is a directory for each package. Files for versions of this package are in these subdirectories. Thus, if our repository hosts packages *spam* (in versions *1.0.0* and *1.1.0*) and *eggs* (in versions *1.0.0* and *1.0.1*) we would have following directory structure :

    $ tree
    .
    ├── spam
    │   ├── spam-1.0.0.tar.gz
    │   └── spam-1.1.0.tar.gz
    └── eggs
        ├── eggs-1.0.0.tar.gz
        └── eggs-1.0.1.tar.gz

You must create this directory and ensure that user running the server has a right to write in this directory.

It is highly advised to backup this directory.

### key

This is the path to the CheeseShop private key. To generate such a key, you might type:

    $ openssl genrsa -out cheeseshop-key.pem 2048

This will generate a file *cheeseshop-key.pem* that you should copy in directory */etc/ssl/private*, which is the standard place.

This is only necessary when running HTTPS server. If you run only HTTP, you may set this value to *~*.

### cert

This is the path to the CheeseShop certificate. To generate a self signed certificate, you can type:

    $ openssl req -new -x509 -key cheeseshop-key.pem -out cheeseshop-cert.pem -days 3650

This command will ask you many fields, but the only that is necessary is the *FQDN* which is the hostname of the machine that is running CheeseShop. A file named *cheeseshop-cert.pem* will be generated; you should copy this file in directory */etc/ssl/certs*, which is the standard place.

Note that if you have a certificate generated by a Certification Authority, you might not have to add a *trusted-host* in your PIP configuration. But I have such certificate and was unable to test it.

### http

This the port number that HTTP server will listen for incoming connections. Set it to *0* to disable HTTP (and run only on HTTPS). Note that it is not a good idea to perform basic authentication on HTTP, as anybody that intercepts HTTP requests might know you username and password. Standard port for HTTP is *80* but the server must run as root to be able to listen on this port. If you don't run the server as root, you must listen on a port number greater than *1024*.

### https

This is the port number that HTTPS server is listening. Set it to *0* to disable HTTPS. If HTTPS is enabled, you must provide private key and certificate (in *key* and *cert* configuration fields). Standard port for HTTPS is *443* but the server must run as root to be able to listen on this port. If you don't run the server as root, you must listen on a port number greater than *1024*.

### path

This is the URL path that the server will listen. Default value is *simple*, thus to list all packages, you should open URL <http://my.shop.host/simple>. To list available version for package *spam*, you would open URL <http://my.shop.host/simple/spam>. To download version *1.2.3* of this package, you would open <http://my.shop.host/simple/spam/spam-1.2.3.tar.gz>. This value should not be changed.

### shop

This is the URL of the public package repository, aka <http://pypi.python.org/simple>. This should not be changed.

### auth

This is the basic authentication configuration. If you don't want authentication, set this value to *~*. This is a list of usernames and MD5 hash of their password. To get the MD5 hash of a given password, you can type following command:

    $ echo -n foo | md5sum
    acbd18db4cc2f85cedef654fccc4a4d8  -

Note that if you modify this configuration, you must restart server, because this configuration is loaded at startup.

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

Then you can use the makefile to build the binary version for your platform:

    $ make build

To build binaries for all platforms, type:

    $ make compile

*Enjoy!*
