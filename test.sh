docker build -t tests -f test.Dockerfile .
docker run --name tests --rm tests