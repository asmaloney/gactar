#!/bin/sh

dir=${0%/*}
if [ -d "$dir" ]; then
  cd "$dir" || exit
fi

cd .. || exit

# create the virtual env dir
python3 -m venv env

# activate it
source ./env/bin/activate

echo "Using python3 from here:" `which python3`
python3 --version

# update pip
pip install --upgrade pip

# should we use requirements.txt?
pip install autopep8 pyactr pylint requests

cd ./env
python3 ../scripts/getImplementations.py
if [ $? -eq 0 ]; then
    echo "SUCCESS"
    echo "Your environment is set up."
    echo "To load your virtual enviroment, run: source ./env/bin/activate"
else
    echo "INSTALLATION ERROR"
    echo "There was a problem setting up your environment. Please check the errors."
    echo "If you can't figure out the problem, search for solutions in the github issues at:"
    echo "  https://github.com/asmaloney/gactar/issues"
fi
