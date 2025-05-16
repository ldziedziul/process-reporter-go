#!/usr/bin/env bash

# init sudo
sudo echo

ansible-playbook -i hosts playbook-deploy.yml $@
