# Contributing to Captain's Log 

The contributors are listed in AUTHORS (add yourself). This project uses the MPL v2 license, see LICENSE.

# Our Process
Before you send a pull request, please familiarize yourself with the [C4.1 Collective Code Construction Contract](http://rfc.zeromq.org/spec:22). A quick summary (but please, do read the process document):
* A Pull Request should be described in the form of a problem statement.
* The code included with the pull request should be a proposed solution to that problem.
* The submitted code should adhere to our style guidelines (described below).
* The submitted code should include tests.
* The submitted code should not break any existing tests.

"A Problem" should be one single clear problem. Large complex problems should be broken down into a series of smaller problems when ever possible.

# Style Guide
* Your code must be formatted with [Gofmt](https://blog.golang.org/go-fmt-your-code)
* Your code should pass [golint](https://github.com/golang/lint). If for some reason it cannot, please provide an explanation.
* Your code should pass [go vet](https://golang.org/cmd/vet/)


