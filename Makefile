permission:
	chmod +x start_services.sh; \
	chmod +x stop_services.sh; \
	chmod +x clean_services.sh;

start:
	./start_services.sh

stop:
	./stop_services.sh

clean:
	./clean_services.sh
