#!/bin/bash
#
#

set -ev

sudo apt-get update
echo y | sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

sudo apt-key fingerprint 0EBFCD88

sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"

#Install Docker and Docker Compose
sudo apt-get update
echo y | sudo apt-get install docker-ce docker-ce-cli containerd.io
sudo apt-cache madison docker-ce
sudo docker run hello-world

#sudo service docker start
sudo systemctl start docker
sudo systemctl enable docker
sudo gpasswd -a ${USER} docker
#Activate the changes to groups
newgrp docker

echo y | sudo apt install docker-compose

#Install Node.js Runtime (v10) and NPM
#curl -sL https://deb.nodesource.com/setup_10.x | sudo -E bash -
sudo apt install nodejs

echo
echo "Log out and log back in so that your group membership is re-evaluated"
echo

#Setup docker images and repositories
mkdir -p ~/go/src/github.com/hyperledger
cd ~/go/src/github.com/hyperledger
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.3.1 1.4.9
