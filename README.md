# alertmanager-history
Build a alertmanager history server for search and analysis, all alerts will save to elasticsearch and ready for search


[![CircleCI](https://circleci.com/gh/zhangmingkai4315/alertmanager-history.svg?style=svg)](https://circleci.com/gh/zhangmingkai4315/alertmanager-history)


1. Frist rename config-example.yml -> config.yml , change the configration parms, then make sure the sevices in your config file is ready to communicate: 

 - alertmanager 
 - elasticsearch

2. Build or download from release page and start history server:

```
  go build . && ./alertmanager-history -c config.yml
```


