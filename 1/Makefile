main: main.go ../bidobo/stub.go ../bidobo/bidobo.s ../bidobo/bidobo.go
	go build -o main

../bidobo/stub.go ../bidobo/bidobo.s: ../bidobo/asm.go
	$(MAKE) -C ../bidobo
