#!/usr/bin/env python

from optparse import OptionParser
import os
import sys
import time
import importlib


def run_plugin(pluginname, conffile, notify, poll):
    try:
        mod = importlib.import_module(pluginname)
    except ModuleNotFoundError:
        sys.stderr.write(f"{pluginname} is not a valid module. Exiting.\n")
        sys.exit(3)

    m = mod.Plugin(open(conffile, "r"), notify)
    try:
        if poll > 0:
            while 1:
                time.sleep(poll)
                m.run()
        else:
            m.run()
    except KeyboardInterrupt:
        sys.stderr.write("Ctrl+c exitting.")
        sys.exit(0)

p = OptionParser()
p.add_option('-l', '--list-modules', action='store_true', dest="LISTMODULES",
             help="List the supported modules")
p.add_option('-P', '--plugin-dir', action='append', dest='PLUGINDIRS',
             help="Directories from which to load plugins")
p.add_option('-p', '--plugin', action='store', dest='PLUGIN',
             help="Plugin to execute")
p.add_option('-n', '--notify', action='store_true', dest="NOTIFY",
             help="Send notifications for results")
p.add_option('-c', '--config', action='store', dest="CONFFILE",
             help="Path to configuration file (defaults to ~/.mystatus.conf")
p.add_option('-d', '--delay', action='store', dest="POLLING_DELAY",
             type=int,
             help="If >0, this is the number of seconds to delay between polling")
p.add_option('-a', '--all', action='store_true', dest="ALL",
             help="If set, return outputs for all plugins")

opts, args = p.parse_args()

pluginbase = os.path.join(os.path.dirname(__file__), 'plugins')
try:
    opts.PLUGINDIRS.append(pluginbase)
except AttributeError:
    opts.PLUGINDIRS = [pluginbase]

if opts.LISTMODULES:
    for plug in opts.PLUGINDIRS:
        for p in os.listdir(plug):
            if p == "__pycache__":
                continue
            print(p[:-3])
    sys.exit(0)

if opts.PLUGIN is None and opts.ALL is None:
    sys.stderr.write("No plugin specified. Exiting.\n")
    sys.exit(3)

if opts.CONFFILE is None:
    conffile = os.path.join(os.path.expanduser("~"), ".mystatus.conf")
else:
    conffile = opts.CONFFILE

for plug in opts.PLUGINDIRS:
    if plug not in sys.path:
        sys.path.append(plug)

if opts.PLUGIN:
    try:
        poll = int(opts.POLLING_DELAY)
    except TypeError:
        poll = 0
    run_plugin(opts.PLUGIN, conffile, opts.NOTIFY, poll)

if opts.ALL:

    plugins = []
    for plug in opts.PLUGINDIRS:
        for p in os.listdir(plug):
            if p == "__pycache__":
                continue
            plugins.append(p[:-3])

    for p in plugins:
        sys.stderr.write(f"----------- {p} --------\n\n")
        run_plugin(p, conffile, opts.NOTIFY, 0)
        sys.stderr.write("\n")
