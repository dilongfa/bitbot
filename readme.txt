Local dev setup
===============
- Install Docker Toolbox (https://www.docker.com/toolbox)
- Add the docker host to /etc/hosts (mac only)
$ echo $(docker-machine ip default) localhost | sudo tee -a /etc/hosts
- Compile and start
$ docker-compose build
$ docker-compose up

- run dB
$ cd citusdb
$ vagrant up

- create database
$ psql -h localhost -p 5433 -d postgres -U vagrant -c "create database bitbot;"

- create table
$ psql -h localhost -p 5433 -d bitbot -U vagrant -f sql/201510072202-create-orderbook.sql

Deployment
==========

- Install Ansible
$ git clone git://github.com/ansible/ansible.git --recursive
$ cd ./ansible
$ source ./hacking/env-setup
$ git pull --rebase
$ git submodule update --init --recursive
$ sudo make install

- run ansible-playbook
$ ansible-playbook ansible/deploy.yaml -i ansible/gce_hosts
