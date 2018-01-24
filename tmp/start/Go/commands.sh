# Build
GOPATH=/Users/luigi/Workspace/Go/branches/go/src GOROOT_BOOTSTRAP=/usr/local/Cellar/go/1.7.3/libexec ./make.bash
branches/go/bin/go test -i golang.org/x/tools/cmd/godoc

# Godoc
branches/go/bin/go build -v golang.org/x/tools/cmd/godoc && ./godoc -v -html cmd/net/test > page.html

# Godoc server
./godoc -http=:8080

# Clone
mkdir -p $GOPATH/src/golang.org/x/tools
git clone https://go.googlesource.com/tools $GOPATH/src/golang.org/x/tools
go build golang.org/x/tools/cmd/goimports

# Staged files
git diff --name-only --cached

# Revert
git checkout ../cmd/godoc/setup-godoc-app.bash
git reset --hard HEAD^

# Show unmerged
git diff --name-status --diff-filter=U

# Commit
git branch
git status
git checkout -b branch
git branch --set-upstream-to origin/master
git commit
git push origin HEAD:refs/for/master

# Contribution guidelines
git codereview sync
git checkout master
git sync
git change master
git mail


# Go commands
gofmt -w file.go

go run

go build

GOOS=windows go build

go install  --> $GOPATH/bin/file

go get github.com/golang/example/hello  --> $GOPATH/bin/hello

# List all packages
go list ...

# List all packages and their dependencies
go list -f "{{.ImportPath}} {{.Imports}}" ./...

# List only packages installed by user
go list ./... | ggrep -P '^([a-z\.\d]+/[a-z\.\d]+/[a-z\.\d]+)$'

go list -f '{{ .Name }}: {{ .Doc }} {{ join .Imports "\n" }} 'fmt

go doc fmt Printf

godoc -html:6060


errcheck

go vet

go test

# Generate gRPC code
protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld

# Test concurrent requests
ab -c 500 -n 500 -l -v 4 http://localhost:8080/


cmd/godoc: show comment before code block

The existing implementation shows a comment after its code block.
This causes confusion regarding how the comment relates to a code block.
Comments have been moved before their relative code blocks.

                                                   Fixes #16728
# MySQL server:
brew services start mysql
mysql -uroot
mysql -h host -u root -p cms


# Curl Post
curl -id "first_name=sausheong&last_name=chang" 127.0.0.1:8080/body
