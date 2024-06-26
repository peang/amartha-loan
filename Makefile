cleanup:
	go mod tidy && go mod vendor

download-file:
	@echo "Downloading file from Google Drive"
	curl -L -o datasets2.parquet "https://drive.usercontent.google.com/download?id=1QLBGFOoKw_3-iM58q4unWfwHmPqfnrYr&export=download&authuser=0&confirm=t&uuid=0e0a5b66-1b23-4d93-a9f8-4814b82a26ae&at=APZUnTWQH6Py1QF0cizaF7hzYyIJ%3A1719379161030"


