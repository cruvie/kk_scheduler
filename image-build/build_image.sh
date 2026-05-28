clear

cd ..

#(
#echo "building ui"
#cd ui && npm run generate
#)
#(
#echo "server vendor"
#go mod vendor
#)

#docker pull --platform linux/amd64 golang:1.26.0
#docker pull --platform linux/arm64 golang:1.26.0
#docker pull --platform linux/amd64 alpine:latest
#docker pull --platform linux/arm64 alpine:latest

export GOVERSION=1.26.0

export TAG=0.2.3


# local build test
#(
#docker build \
#--build-arg GOVERSION=${GOVERSION} \
#-t kk_scheduler-local:${TAG} \
#-f ./image-build/Dockerfile .
#)

# docker login

#(
#docker buildx build --platform linux/amd64,linux/arm64  \
#            --build-arg GOVERSION=${GOVERSION} \
#            -t cruvie/kk_scheduler:${TAG} \
#            -t cruvie/kk_scheduler:latest \
#            -f ./image-build/Dockerfile . \
#            --load
#)
(
docker push cruvie/kk_scheduler:${TAG}
docker push cruvie/kk_scheduler:latest
)

