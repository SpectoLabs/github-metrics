deps:
	go get github.com/stretchr/testify/assert

test: deps
	hoverctl start && \
	hoverctl mode simulate && \
	hoverctl import simulation.json && \
	go test -v && \
	hoverctl stop
