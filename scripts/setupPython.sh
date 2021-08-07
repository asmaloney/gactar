#!/bin/sh

dir=${0%/*}
if [ -d "$dir" ]; then
  cd "$dir" || exit
fi

cd .. || exit

# create the virtual env dir
python3 -m venv pyenv

# activate it
source ./pyenv/bin/activate

echo "Using python3 from here:" `which python3`
python3 --version

# update pip
pip install --upgrade pip

# because CCMSuite isn't a proper package, we need to copy files to the right place
cd ./pyenv
git clone https://github.com/CarletonCognitiveModelingLab/CCMSuite3
cp -rpf CCMSuite3/ccm lib/python3.9/site-packages/

echo "Your environment is set up."
echo "To load your virtual enviroment, run: source ./pyenv/bin/activate"