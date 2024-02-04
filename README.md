# Go API REST Example

Based on [Francesco Ciulla Video](https://www.youtube.com/watch?v=aLVJY-1dKz8)

I added/changed a few things

- [x] A server struct contains db and router. IMO to make testing and mocking easier.
- [x] unit test with sqlmock.
- [x] Dockerfile multistage.
- [x] gitignore.
- [x] [modd](https://github.com/cortesi/modd) to run the tests when a file is modified.


Run test

    go test

Watching the files and running the tests automatically with [modd](https://github.com/cortesi/modd)

    modd
