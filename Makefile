deploy:
	rsync -avhzL --delete \
				--no-perms --no-owner --no-group \
				--exclude .git \
				--exclude-from=".gitignore" \
				. ubuntu@13.125.136.45:/home/ubuntu/heta

reset:
	rm -rf .data
	LOCAL_CLIENT_ID=1 go run cli/main.go init
	LOCAL_CLIENT_ID=2 go run cli/main.go init
	LOCAL_CLIENT_ID=3 go run cli/main.go init
	LOCAL_CLIENT_ID=4 go run cli/main.go init
	LOCAL_CLIENT_ID=5 go run cli/main.go init

test-bootnode:
	go run main.go

test-producer-1:
	LOCAL_CLIENT_ID=1 HTTP_SERVICE_PORT=9001 go run main.go --producer --port 9101 \
	--bootnodes "enode://0e8eeac11d6d7eecad89a3b9606a8c494e986a4d7250557c5d5b658c7b5d79997048cca5e0f559798ea5455626a4814a01fc4b5df3120e8f9d76edc2e0be4b00@127.0.0.1:30301"

test-producer-2:
	LOCAL_CLIENT_ID=2 HTTP_SERVICE_PORT=9002 go run main.go --producer --port 9102 \
	--bootnodes "enode://0e8eeac11d6d7eecad89a3b9606a8c494e986a4d7250557c5d5b658c7b5d79997048cca5e0f559798ea5455626a4814a01fc4b5df3120e8f9d76edc2e0be4b00@127.0.0.1:30301"

test-producer-3:
	LOCAL_CLIENT_ID=3 HTTP_SERVICE_PORT=9003 go run main.go --producer --port 9103 \
	--bootnodes "enode://0e8eeac11d6d7eecad89a3b9606a8c494e986a4d7250557c5d5b658c7b5d79997048cca5e0f559798ea5455626a4814a01fc4b5df3120e8f9d76edc2e0be4b00@127.0.0.1:30301"

test-producer-4:
	LOCAL_CLIENT_ID=4 HTTP_SERVICE_PORT=9004 go run main.go --producer --port 9104 \
	--bootnodes "enode://0e8eeac11d6d7eecad89a3b9606a8c494e986a4d7250557c5d5b658c7b5d79997048cca5e0f559798ea5455626a4814a01fc4b5df3120e8f9d76edc2e0be4b00@127.0.0.1:30301"

test-client:
	LOCAL_CLIENT_ID=5 HTTP_SERVICE_PORT=9005 go run main.go --port 9105 \
	--bootnodes "enode://0e8eeac11d6d7eecad89a3b9606a8c494e986a4d7250557c5d5b658c7b5d79997048cca5e0f559798ea5455626a4814a01fc4b5df3120e8f9d76edc2e0be4b00@127.0.0.1:30301"
