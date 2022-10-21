cd ../..
gitId=$(git show HEAD --pretty=format:"%H" --no-patch)
cd -
go get github.com/BalazsNyiro/iris/iris@${gitId}
go mod vendor

go run iris_windows.go
