.PHONY: generate_homm2

run_test: build_test
	./test

build_test:
	go build

generate_homm2: homm2map.sbf
	python generator.py homm2map.sbf homm2 homm2/homm2.go
