#!/bin/bash

# 固定リポジトリの script/start.sh を test モードで実行
curl -fsSL https://raw.githubusercontent.com/kdinstall/system-base3/master/script/start.sh | bash -s -- -test
