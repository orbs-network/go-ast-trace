rm -rf ./test/output
cp -R ./test/input ./test/output

go run *.go locks ./test/output/*.go

echo "helloworld.go:"
go run ./test/output/helloworld.go
echo ""

echo "mutex.go:"
go run ./test/output/mutex.go
echo ""

echo "channel.go:"
go run ./test/output/channel.go
echo ""