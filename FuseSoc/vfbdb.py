#!/usr/bin/python3
"""
FuseSoc wrapper for vfbdb.
"""

import os
import sys
import yaml
import subprocess

if __name__ == "__main__":

    with open(sys.argv[1], 'r') as config_file:
        config = yaml.safe_load(config_file)

    files_root = config['files_root']
    os.environ['FBDPATH'] = files_root

    try:
        main = os.path.join(files_root, config['parameters']['main'])
    except:
        print("ERROR: Input .fbd main file ('main' parameter) musts be specified!")
        sys.exit(1)

    args = ['vfbdb', '-fusesoc', '-fusesoc-vlnv', config['vlnv']]
    global_args = []
    targets_args = []

    for param, val in config['parameters'].items():
        if param == 'main':
            continue

        if param == 'global':
            args_list = global_args
        else:
            targets_args.append(param)
            args_list = targets_args

        prev_v = None
        for v in val:
            if prev_v == '-path':
                args_list.append(os.path.join(files_root, v))
            else:
                args_list.append(v)
            prev_v = v

    args = args + global_args + targets_args
    args.append(main)

    ret = subprocess.run(args)
    if ret.returncode != 0:
        exit(ret.returncode)
