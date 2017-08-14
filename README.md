# URL SHORTENER API
===========================

# Software Requirements
1. Go SDK +1.6 (https://golang.org)
2. GO ORM (https://github.com/jinzhu/gorm)
3. Viper (https://github.com/spf13/viper)
4. Gorilla Mux (https://github.com/gorilla/mux)
5. Glide (https://github.com/Masterminds/glide)

# Hardware Requirements
1. Memory  : +- 50MB
2. Storage : +- 10MB
3. CPU     : 1 Core

# Go Installation
> download Go SDK from https://golang.org/dl/

> extract using command : `tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz`

> edit ``$HOME/.profile` or `./bashrc` using text editor.

> add `export PATH=$PATH:/usr/local/go/bin`

> create folder workspace, inside folder workspace create below structure folder.

    `> src
     > pkg
     > bin`

> Set GOPATH value to path of our workspace folder.

> export GOPATH=$PATH_TO_YOUR_FOLDER/workspace

> test if go has successfully configured in system by run `go env`.

# Project installation (dependency)
> run ./glide install

# Build application
> $PATH/url-shortener/go build


# Run application
> $PATH/url-shortener/ ./url-shortener

# API endpoint
> Create / Shorten URL (POST)

> http://[YOUR_HOST]:[YOUR_PORT]/shortener/create

> Get shorten URL (GET)

> http://[YOUR_HOST]:[YOUR_PORT]/shortener/{id}
