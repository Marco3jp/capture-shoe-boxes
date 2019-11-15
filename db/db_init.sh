#!/bin/bash
mysql -uroot -p$ROOT_PASSWORD -e "DROP DATABASE IF EXISTS is_exist_researcher; CREATE DATABASE is_exist_researcher;"
mysql -uroot -p$ROOT_PASSWORD -e "DROP USER IF EXISTS capture_shoe_boxes@localhost; DROP USER IF EXISTS diff_shoe_boxes@localhost; DROP USER IF EXISTS is_exist_researcher_api@localhost;"
mysql -uroot -p$ROOT_PASSWORD is_exist_researcher <./tables.sql
mysql -uroot -p$ROOT_PASSWORD <./users.sql
