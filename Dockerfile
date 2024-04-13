FROM golang:1.22

# set working directory
WORKDIR /jane2

# copy source code
COPY . .

RUN chmod +x /jane2

RUN apt update && apt-get install -y ffmpeg

# dependencies
RUN go mod download

# build
RUN CGO_ENABLED=0 go build -o jane2.0

# run
 CMD ["./jane2.0"]