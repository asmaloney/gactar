"""The install_vanilla module will download & install ACT-R (Lisp) and the SBCL Lisp compiler if possible."""
import os
import platform
import shutil
import subprocess
import sys

import pathlib
import requests

# I am not a python person! I'm using python rather than shell script for portability
# and so we don't require more tools to be installed.
#
# If you know how to do this better, please submit an issue:
#   https://github.com/asmaloney/gactar/issues
# or a pull request with fixes:
#   https://github.com/asmaloney/gactar/pulls


def remove_file(file_name: str):
    """Remove a file from the file system."""
    if os.path.isfile(file_name):
        try:
            os.remove(file_name)
        except OSError as err:
            print(err)
            sys.exit()


def remove_dir(dir_name: str):
    """Remove a directory from the file system."""
    if os.path.isdir(dir_name):
        try:
            shutil.rmtree(dir_name)
        except OSError as err:
            print(err)
            sys.exit()


def unpack_file(file_name: str):
    """Unpack a file."""
    if os.path.isfile(file_name):
        shutil.unpack_archive(file_name)


def download_url(url: str) -> str:
    """Download a URL."""
    local_filename = url.split('/')[-1]
    with requests.get(url, stream=True) as response:
        with open(local_filename, 'wb') as local_file:
            shutil.copyfileobj(response.raw, local_file)

    return local_filename


def download_vanilla():
    """Download the lisp ACT-R files and install in the correct location."""
    print('Downloading and installing Vanilla ACT-R...')

    version = 'v7.27.0'

    url = f'https://github.com/asmaloney/ACT-R/releases/download/{version}/actr-super-slim-{version}.zip'
    unpacked_dir = f'actr-super-slim-{version}'
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


def download_sbcl() -> str:
    """Download the SBCL (Lisp) compiler and install in the correct location."""
    # See: http://www.sbcl.org/platform-table.html
    system = platform.system()
    arch = platform.machine()

    if system == 'Darwin':
        print('Downloading and installing SBCL...')
        system = 'darwin'

        version = '2.1.2'  # latest arm64 version

        if arch == 'x86_64':
            version = '1.2.11'  # latest x86_64 version
            arch = 'x86-64'     # URL needs hyphen

        dir_name = f'sbcl-{version}-{arch}-{system}'
    else:
        raise Exception(
            f'ERROR: I don\'t know how to install the sbcl compiler for your platform({system} - {arch})\n'
            'Please see the gactar README for how to download and setup sbcl.'
        )

    # remove old file if it exists for some reason
    remove_dir(dir_name)

    # download sbcl
    url = f'https://prdownloads.sourceforge.net/sbcl/{dir_name}-binary.tar.bz2'

    compressed_file = download_url(url)
    unpack_file(compressed_file)

    return dir_name


def install_vanilla(dir_name: str):
    print(f'Installing from {dir_name}')

    base_dir = pathlib.Path().resolve()

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
    remove_dir(dir_name)


if __name__ == "__main__":
    try:
        # download_vanilla()
        dir_name = download_sbcl()
        install_vanilla(dir_name)
    except BaseException as err:
        print(err)
        sys.exit()
