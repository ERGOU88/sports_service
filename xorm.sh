#!/bin/bash
# 数据表->数据结构
# tips: ./xorm reverse mysql $username:$password@tcp\($host\)/$dbname?charset=utf8 配置文件目录 输出目录
cd ../pkg/mod/github.com/go-xorm/cmd/xorm@v0.0.0-20190426080617-f87981e709a1
xorm reverse mysql root:bluetrans888@tcp\(192.168.5.12:3306\)/sports_service?charset=utf8 ../../../../../../server/models/modelTpl ../../../../../../server/models/
