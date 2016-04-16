Install
=======

Vagrant Install
---------------

Install Virtualbox and Vagrant.

Create the Vagrant machine:
$ vagrant up

Connect to the VM:
$ vagrant ssh

Create MySQL database:
$ mysql -u bitbot -p bitbot < db/init.sql
password

Start the services
------------------

You first need to compile the Go code and build the JavaScript application:
$ cd /vagrant
$ make

Now you can start the "record" service which fetches Bitcoin orderbooks from several exchangers:
$ bin/record

Start the "web" service that power the UI and provides the API: (you will another SSH session for now):

$ bin/web

Open your browser and point it at localhost:8080

Note: you will need to run make and restart the services each time you make a change to the code.


Deploy the code on GCE
======================

Create a VM instance and disk
-----------------------------
$ ./ansible/setup.sh

Deployment
----------

- Install Ansible
$ git clone git://github.com/ansible/ansible.git --recursive
$ cd ./ansible
$ source ./hacking/env-setup
$ git pull --rebase
$ git submodule update --init --recursive
$ sudo make install

- run ansible-playbook
$ ansible-playbook ansible/deploy.yaml -i ansible/gce_hosts
