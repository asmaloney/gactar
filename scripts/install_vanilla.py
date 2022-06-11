"""The install_vanilla module will download & install ACT-R (Lisp) and the Clozure Common Lisp compiler."""
import os
import platform
import shutil
import sys

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


def download_ccl():
    """Download the Clozure Common Lisp compiler."""
    # See: https://github.com/Clozure/ccl
    system = platform.system().lower()

    if system != 'darwin' and system != 'linux' and system != 'windows':
        raise Exception(
            f'ERROR: I don\'t know how to install the Clozure Common Lisp compiler for your platform ({system})\n'
            '\tPlease see the gactar README for how to download and setup ccl.'
        )

    version = '1.12.1'
    extension = 'tar.gz'

    if system == 'windows':
        extension = 'zip'

    print(
        f'Downloading and installing Clozure Common Lisp v{version} for {system}...')

    dir_name = f'ccl-{version}-{system}x86'

    # remove old file if it exists
    remove_dir(dir_name)

    # download ccl
    url = f'https://github.com/Clozure/ccl/releases/download/v{version}/{dir_name}.{extension}'

    compressed_file = download_url(url)
    unpack_file(compressed_file)


if __name__ == "__main__":
    try:
        download_vanilla()
        download_ccl()
    except BaseException as err:
        print(err)
        sys.exit()
