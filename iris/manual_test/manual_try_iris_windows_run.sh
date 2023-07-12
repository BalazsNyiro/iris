# manual local test 

cp ../iris*.go .
rm *test*.go
sed --in-place "s/package iris/package main/g" iris*.go
go run manual_test_iris_windows.go *.go

