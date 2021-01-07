.PHONY: cross help clean mk dist cpu_torch cpu_pprof mem_torch mem_pprof trace doc

cmd :
	mkdir -p cmd
	go build -o cmd/example/douglas-peuker example/douglas-peuker/main.go
	go build -o cmd/example/drift example/drift/main.go
	go build -o cmd/example/kalman-filter example/kalman-filter/main.go
	go build -o cmd/example/random-buffer-track example/random-buffer-track/main.go
	go build -o cmd/example/random-drift-track example/random-drift-track/main.go


douglas-peuker :
	cmd/example/douglas-peuker

drift :
	cmd/example/drift

kalman-filter :
	cmd/example/kalman-filter

random-buffer-track :
	cmd/example/random-buffer-track

random-drift-track :
	cmd/example/random-drift-track


clean :
	@rm -rf cmd
	@rm -rf *.html
