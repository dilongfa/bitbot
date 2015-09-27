sudo apt-get update;
sudo apt-get install -y curl git;

curl -sSL https://get.docker.com/ | sh;

echo -e "Host bitbucket.org\n\tStrictHostKeyChecking no\n" >> ~/.ssh/config;
git clone git@bitbucket.org:bhirbec/bitbot.git;

sudo docker build -t recorder-img bitbot/src;
sudo docker run -d --name exchanger-db -p 9200:9200 -v /var/data/elasticsearch:/usr/share/elasticsearch/data elasticsearch:1.7.1;
sudo docker run --name recorder -d --link exchanger-db recorder-img /go/bin/record  -p 300;
