FROM golang:1.15-alpine as server-build
# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .


FROM iotivity-esp
ARG root_directory=/mf-builder
WORKDIR ${root_directory} 

ENV MF_EMBEDDED_SRC_PATH=${root_directory}/mf-embedded/src
ENV GEN_H_PATH=${root_directory}/mf-embedded/src/components/include/gen
ENV IOTIVITY_PATH=${root_directory}/iotivity-lite



# Clone mf-embedded

# Do it like this because the ADD obtains and stores
# the GitHub API's representation of your repo to a dummy location.
# The API call will return different results when the head changes, invalidating the docker cache.
# If not, Docker would Cache the git clone command
ADD https://api.github.com/repos/DavidMarquezF/mf-embedded/git/refs/heads/main version.json
RUN git clone -b main https://github.com/DavidMarquezF/mf-embedded.git 

# Builder code
COPY --from=server-build /build/main .
COPY run.sh /usr/local/bin/run.sh
#RUN go build -ldflags "-linkmode external -extldflags -static" -o builder .

CMD [ "run.sh" ]


