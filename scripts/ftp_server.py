#!/usr/bin/env python

# Copyright (C) 2007-2016 Giampaolo Rodola' <g.rodola@gmail.com>.
# Use of this source code is governed by MIT license that can be
# found in the LICENSE file.

"""A basic FTP server which uses a DummyAuthorizer for managing 'virtual
users', setting a limit for incoming connections.
"""

import json
import os
import argparse

from pyftpdlib.authorizers import DummyAuthorizer
from pyftpdlib.handlers import FTPHandler
from pyftpdlib.servers import FTPServer


def main():
    # Parse cmdline
    parser = argparse.ArgumentParser(formatter_class=argparse.ArgumentDefaultsHelpFormatter)
    parser.add_argument('-c', '--conf', help='configuration file', default='config.json')
    parser.add_argument('-r', '--root', help='root dir', default='./')
    parser.add_argument('-p', '--port', help='listen port', type=int)
    flag = parser.parse_args()
    print flag 
    
    # Load configuration
    cfg = json.load(open(flag.conf, 'rb'))

    # Instantiate a dummy authorizer for managing 'virtual' users
    authorizer = DummyAuthorizer()

    # Define a new user having full r/w permissions and a read-only
    perm = 'elradfmwM'
    for acl in cfg.get('acls', []):
        print flag.root
        base_dir = os.path.join(flag.root, acl.get('directory', './').lstrip('/'))
        print base_dir
        authorizer.add_user(acl['username'], acl['password'], 
                 base_dir, perm=perm)

    # anonymous user
    if cfg.get('anonymous'):
        authorizer.add_anonymous(os.getcwd())

    # Instantiate FTP handler class
    handler = FTPHandler
    handler.authorizer = authorizer

    # Define a customized banner (string returned when client connects)
    handler.banner = cfg.get('banner', "pyftpdlib based ftpd ready.")

    # Specify a masquerade address and the range of ports to use for
    # passive connections.  Decomment in case you're behind a NAT.
    # handler.masquerade_address = '151.25.42.11'
    handler.passive_ports = range(60000, 65535)

    # Instantiate FTP server class and listen on 0.0.0.0:2121
    address = ('', flag.port or cfg.get('port', 2121))
    server = FTPServer(address, handler)

    # set a limit for connections
    server.max_cons = 256
    server.max_cons_per_ip = 5

    # start ftp server
    server.serve_forever()

if __name__ == '__main__':
    main()
