

launch:
	nohup ./orochi --port 3000 >> log.txt 2>&1 &
	nohup ./orochi --port 3001 >> log.txt 2>&1 &
	nohup ./orochi --port 3002 >> log.txt 2>&1 &

stop:
	pkill -9 -f orochi
