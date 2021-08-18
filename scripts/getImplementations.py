import io
import os
import pathlib
import requests
import shutil
import subprocess
import sys
import zipfile

# I am not a python person! I'm using python rather than shell script for portability
# and so we don't require more tools to be installed.
#
# If you know how to do this better, please submit an issue:
#   https://github.com/asmaloney/gactar/issues
# or a pull request with fixes:
#   https://github.com/asmaloney/gactar/pulls


def removeFile(fileName):
    if os.path.isfile(fileName):
        try:
            os.remove(fileName)
        except OSError as err:
            print(err)
            sys.exit()


def removeDir(dirName):
    if os.path.isdir(dirName):
        try:
            shutil.rmtree(dirName)
        except OSError as err:
            print(err)
            sys.exit()


def unpackFile(fileName):
    if os.path.isfile(fileName):
        shutil.unpack_archive(fileName)


def download_url(url):
    local_filename = url.split('/')[-1]
    with requests.get(url, stream=True) as r:
        with open(local_filename, 'wb') as f:
            shutil.copyfileobj(r.raw, f)

    return local_filename


def downloadCCM():
    print('Downloading and installing CCM-PyACTR...')

    # Because CCMSuite isn't a proper package and we can't use pip, we need to copy files to the right place.
    # I'm using a fork (CCM-PyACTR) to avoid pulling all the tmp and .pyc files in the original repo.

    url = 'https://github.com/asmaloney/CCM-PyACTR/archive/refs/heads/master.zip'
    unpackedDir = 'CCM-PyACTR-master'
    targetDir = 'lib/python3.9/site-packages/ccm'

    # remove old files if they exists for some reason
    removeDir(unpackedDir)
    removeDir(targetDir)

    # get the CCM-PyACTR files
    zipFile = download_url(url)
    unpackFile(zipFile)

    shutil.move(unpackedDir + '/ccm', targetDir)

    # clean up
    removeFile(zipFile)
    removeDir(unpackedDir)


def downloadVanilla():
    print('Downloading and installing Vanilla ACT-R...')

    url = 'https://github.com/asmaloney/ACT-R/archive/refs/tags/v7.21.6.tar.gz'
    unpackedDir = 'ACT-R-7.21.6'
    targetDir = 'actr'

    # remove old files if they exists for some reason
    removeDir(unpackedDir)
    removeDir(targetDir)

    # get the ACT-R files
    zipFile = download_url(url)
    unpackFile(zipFile)

    os.rename(unpackedDir, targetDir)

    # clean up
    removeFile(zipFile)


def downloadSBCL():
    import platform
    sys = platform.system()
    platform = platform.machine()

    if sys == 'Darwin':
        print('Downloading and installing SBCL...')
        sys = 'darwin'

        if platform == 'x86_64':
            platform = 'x86-64'
        # arm64 should just work...

        baseDir = pathlib.Path().resolve()
        dirName = 'sbcl-1.2.11-' + platform + '-' + sys
        unpackedDir = dirName
        compressedFile = dirName + '-binary.tar.bz2'
        url = 'https://prdownloads.sourceforge.net/sbcl/' + compressedFile

        # remove old file if it exists for some reason
        removeDir(unpackedDir)

        # download sbcl
        compressedFile = download_url(url)
        unpackFile(compressedFile)

        # run the install
        os.environ['INSTALL_ROOT'] = str(baseDir)
        os.chdir(dirName)
        result = subprocess.run(['./install.sh'], env=os.environ)
        result.check_returncode()
        os.chdir(baseDir)

        # compile the actr files
        os.unsetenv('INSTALL_ROOT')
        os.environ['SBCL_HOME'] = str(baseDir / 'lib/sbcl')
        result = subprocess.run(['./bin/sbcl', '--script',
                                 'actr/load-single-threaded-act-r.lisp'], env=os.environ)
        result.check_returncode()

        # clean up
        removeFile(compressedFile)
        removeDir(unpackedDir)
    else:
        print(
            "ERROR: I don't know how to install the sbcl compiler for your platform (" + sys + ' - ' + platform + ')')
        print("Please see the gactar README for how to download and setup sbcl.")


if __name__ == "__main__":
    downloadCCM()
    downloadVanilla()
    downloadSBCL()
