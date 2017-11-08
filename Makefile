
OUTPUT_DIR=./_output

build: 
	go build -o ${OUTPUT_DIR}/kclient ./

product: 
	env GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/kclient.linux ./

clean:
	rm -rf ${OUTPUT_DIR}/*
