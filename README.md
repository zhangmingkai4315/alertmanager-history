# alertmanager-history
Build a alertmanager history server for search and analysis, all alerts will save to elasticsearch and ready for search


#### Start Server 

[![CircleCI](https://circleci.com/gh/zhangmingkai4315/alertmanager-history.svg?style=svg)](https://circleci.com/gh/zhangmingkai4315/alertmanager-history)


Frist rename config-example.yml -> config.yml , change the configration parms, then make sure your sevices is ready: 

 - alertmanager 
 - elasticsearch

Build yourself, or download from release page and start history server:

```
  go build . && ./alertmanager-history -c config.yml
```


