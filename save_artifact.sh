#! /bin/sh
mkdir dist
for f in ./tmp/*;do
  filename=`basename "$f"`
  tar czvf dist/$filename.tar.gz $f config-example.yml
done 