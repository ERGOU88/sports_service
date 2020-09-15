#!/bin/bash
# 数据表->数据结构

cd ../pkg/mod/github.com/go-xorm/cmd/xorm@v0.0.0-20190426080617-f87981e709a1
xorm reverse mysql root:a3202381@tcp\(127.0.0.1:3306\)/fpv2?charset=utf8 ../../../../../../server/models/modelTpl ../../../../../../server/models/
