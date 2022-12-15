cd ../..
gitId=$(git show HEAD --pretty=format:"%H" --no-patch)
echo "git id: $gitId"
cd -
echo "go get..."
go get github.com/BalazsNyiro/iris/iris@${gitId}
echo "go tidy - remove unused dependencies"
go tidy
echo "go mod vendor..."
go mod vendor

echo "go run..."
go run iris_windows.go
