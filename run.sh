#!/bin/bash

# コンテナの作成
docker-compose up -d --build

# 起動したコンテナにログイン
docker exec -it sqltest-mysql-1 bash -p

# MySQLを起動
#mysql -u root -p -h 127.0.0.1