#!/bin/bash
# 数据表->数据结构
# tips: ./xorm reverse mysql $username:$password@tcp\($host\)/$dbname?charset=utf8 配置文件目录 输出目录
xorm reverse mysql root:123456@tcp\(127.0.0.1:3306\)/sports_service?charset=utf8 models/modelTpl models/
xorm reverse mysql root:123456@tcp\(127.0.0.1:3306\)/venue?charset=utf8 models/modelTpl models/
