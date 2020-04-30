# beego-bbs
Quick Start:
bee run -gendoc=true -downdoc=true

Go to:
http://127.0.0.1:8000/swagger/ 


conf文件夹：系统配置信息和数据库配置信息。本程序数据库使用postgresql，需要替换数据库可以到main.go中做替换。

系统主要包含三个模块：用户管理(user)、发帖(topic)、回复(reply)，user模块使用JWT做认证。
