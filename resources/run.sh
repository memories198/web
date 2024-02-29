rm -rf /home/dbManager/*
cp -a redis /home/dbManager/
cp -a mysql /home/dbManager/
cp docker-compose.yaml /home/dbManager/
chown -R  dbManager:dbManager /home/dbManager/
su - dbManager