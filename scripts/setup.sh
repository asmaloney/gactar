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
pip install wget pyactr

# Because CCMSuite isn't a proper package, we need to copy files to the right place.
# I'm using a fork (CCM-PyACTR) to avoid pulling all the tmp and .pyc files in the original repo.
cd ./env
git clone https://github.com/asmaloney/CCM-PyACTR.git
cp -rpf CCM-PyACTR/ccm lib/python3.9/site-packages/

echo "Your environment is set up."
echo "To load your virtual enviroment, run: source ./env/bin/activate"