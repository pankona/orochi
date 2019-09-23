
APPNAME    = orochi
EXECUTABLE = $(CURDIR)/cmd/orochi/$(APPNAME)

all:
	cd $(CURDIR)/cmd/$(APPNAME); go build

launch:
	nohup $(EXECUTABLE) --port 3000 >> log.txt 2>&1 &
	nohup $(EXECUTABLE) --port 3001 >> log.txt 2>&1 &
	nohup $(EXECUTABLE) --port 3002 >> log.txt 2>&1 &

stop:
	pkill -9 -f $(APPNAME)
