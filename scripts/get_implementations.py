"""The get_implementations module will download & install multiple ACT-R frameworks."""
import os
import platform
import shutil
import subprocess
import sys
import sysconfig

import pathlib
import requests

# I am not a python person! I'm using python rather than shell script for portability
# and so we don't require more tools to be installed.
#
# If you know how to do this better, please submit an issue:
#   https://github.com/asmaloney/gactar/issues
# or a pull request with fixes:
#   https://github.com/asmaloney/gactar/pulls


def remove_file(file_name):
    """Remove a file from the file system."""
    if os.path.isfile(file_name):
        try:
            os.remove(file_name)
        except OSError as err:
            print(err)
            sys.exit()


def remove_dir(dir_name):
    """Remove a directory from the file system."""
    if os.path.isdir(dir_name):
        try:
            shutil.rmtree(dir_name)
        except OSError as err:
            print(err)
            sys.exit()


def unpack_file(file_name):
    """Unpack a file."""
    if os.path.isfile(file_name):
        shutil.unpack_archive(file_name)


def download_url(url):
    """Download a URL."""
    local_filename = url.split('/')[-1]
    with requests.get(url, stream=True) as response:
        with open(local_filename, 'wb') as local_file:
            shutil.copyfileobj(response.raw, local_file)

    return local_filename


def download_ccm():
    """Download The CCMSuite and install the files we need in the correct location."""
    print('Downloading and installing CCM-PyACTR...')

    # Because CCMSuite isn't a proper package and we can't use pip, we need
    # to copy files to the right place.
    # I'm using a fork (CCM-PyACTR) to avoid pulling all the tmp and .pyc
    # files in the original repo.

    url = 'https://github.com/asmaloney/CCM-PyACTR/archive/refs/heads/master.zip'
    unpacked_dir = 'CCM-PyACTR-master'
    target_dir = sysconfig.get_paths()["purelib"] + '/ccm'

    # remove old files if they exists for some reason
    remove_dir(unpacked_dir)
    remove_dir(target_dir)

    # get the CCM-PyACTR files
    zip_file = download_url(url)
    unpack_file(zip_file)

    shutil.move(unpacked_dir + '/ccm', target_dir)

    # clean up
    remove_file(zip_file)
    remove_dir(unpacked_dir)


def download_vanilla():
    """Download the lisp ACT-R files and install in the correct location."""
    print('Downloading and installing Vanilla ACT-R...')

    url = 'https://github.com/asmaloney/ACT-R/releases/download/v7.21.6/actr-super-slim-v7.21.6.zip'
    unpacked_dir = 'actr-super-slim-v7.21.6'
    target_dir = 'actr'

    # remove old files if they exists for some reason
    remove_dir(unpacked_dir)
    remove_dir(target_dir)

    # create dir and change to it
    os.mkdir(target_dir)
    os.chdir(target_dir)

    # get the ACT-R files
    zip_file = download_url(url)
    unpack_file(zip_file)

    # clean up
    remove_file(zip_file)
    os.chdir('..')


def download_sbcl():
    """Download the SBCL (Lisp) compiler and install in the correct location."""
    # See: http://www.sbcl.org/platform-table.html
    system = platform.system()
    arch = platform.machine()

    if system == 'Darwin':
        print('Downloading and installing SBCL...')
        system = 'darwin'

        version = "2.1.2"  # latest arm64 version

        if arch == 'x86_64':
            version = "1.2.11"  # latest x86_64 version
            arch = 'x86-64'     # URL needs hyphen

        base_dir = pathlib.Path().resolve()

        dir_name = 'sbcl-' + version + '-' + arch + '-' + system
        unpacked_dir = dir_name
        compressed_file = dir_name + '-binary.tar.bz2'
        url = 'https://prdownloads.sourceforge.net/sbcl/' + compressed_file

        # remove old file if it exists for some reason
        remove_dir(unpacked_dir)

        # download sbcl
        compressed_file = download_url(url)
        unpack_file(compressed_file)

        # run the install
        os.environ['INSTALL_ROOT'] = str(base_dir)
        os.chdir(dir_name)
        subprocess.run(['./install.sh'], check=True, env=os.environ)
        os.chdir(base_dir)

        # compile the actr files
        os.unsetenv('INSTALL_ROOT')
        os.environ['SBCL_HOME'] = str(base_dir / 'lib/sbcl')
        subprocess.run(['./bin/sbcl', '--script',
                        'actr/load-single-threaded-act-r.lisp'], check=True, env=os.environ)

        # clean up
        remove_file(compressed_file)
        remove_dir(unpacked_dir)
    else:
        print(
            "ERROR: I don't know how to install the sbcl compiler for your platform ("
            + system + ' - ' + arch + ')')
        print("Please see the gactar README for how to download and setup sbcl.")


if __name__ == "__main__":
    download_ccm()
    download_vanilla()
    download_sbcl()
