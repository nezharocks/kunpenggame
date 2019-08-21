#env GOOS=linux GOARCH=arm
#go build -v -o ./bin/client src/client/main.go
git archive -o dist/henryleu.zip HEAD
#zip -r dist/henryleu.zip . -x "./.git/*" "./dist/*" "./src/*" "./temp/*" "./.gitignore" "./.DS_Store"  "./dist.sh"
