#!/bin/sh

dir=${0%/*}
if [ -d "$dir" ]; then
  cd "$dir" || exit
fi

cd .. || exit

# Windows
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
    # check if python3 command exists, make symbolic link to python.exe if not
    if ! [ -x "$(command -v python3)" ]; then
        cd ./scripts
        touch symlink.bat
        echo -e "@ECHO OFF" > symlink.bat
        echo -e "FOR /F \"tokens=*\" %%g IN ('WHERE python') do (SET PY_PATH=%%g)" >> symlink.bat
        echo -e "\nSET PY3_PATH=%PY_PATH:~0,-4%3.exe" >> symlink.bat
        echo -e "\nmklink %PY3_PATH% %PY_PATH%" >> symlink.bat
        cmd.exe /C symlink.bat
        rm symlink.bat
        if ! [ -x "$(command -v python3)" ]; then
            echo "ERROR setting symbolic link."
            exit 1
        else
            echo "Symbolic link set successfully."
        fi
        cd ..
    fi
fi

# create the virtual env dir
python3 -m venv env

# activate it
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
    source ./env/Scripts/activate
else
    . ./env/bin/activate
fi

echo "Using python3 from here:" $(which python3)
python3 --version

# update pip
pip install --upgrade pip

# install required packages
pip install -r scripts/requirements.txt

# run our own script to download and install non-pip-compatible things
cd ./env
python3 ../scripts/install_vanilla.py
if [ $? -eq 0 ]; then
    echo "SUCCESS"
    echo "Your environment is set up."
    echo "To get help on gactar's command line options, run: ./gactar help"
else
    echo "INSTALLATION ERROR"
    echo "There was a problem setting up your environment. Please check the errors."
    echo "If you can't figure out the problem, search for solutions in the github issues at:"
    echo "  https://github.com/asmaloney/gactar/issues"
fi
