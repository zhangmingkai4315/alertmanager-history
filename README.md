# alertmanager-history
Build a alertmanager history server for search and analysis, all alerts will save to elasticsearch


#### Start Server 


Rename config-example.yml -> config.yml , change the configration parms

Make sure the sevices is ready 

 - alertmanager 
 - elasticsearch


Start History Server
```
  go build . && ./alertmanager-history -c config.yml
```

