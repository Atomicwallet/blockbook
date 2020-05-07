# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure('2') do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://vagrantcloud.com/search.
  config.vm.box = 'atomlab/bionic64'
  config.vm.box_version = "0.0.2"
  config.vm.synced_folder '.', '/src', type: "rsync", rsync__auto: "true"
  config.ssh.private_key_path = '~/.ssh/id_rsa_develop'
  config.vm.provider 'hetznercloud' do |hcloud, override|

    ## Tokens in Vagrantfiles should be avoided. You can use:
    #  - priority1: env var "HETZNERCLCOUD_TOKEN"
    #  - priority2: content in ~/.config/hcloud/cli.toml
    #  - priority3: hcloud.token in this file (last resort as it tends to leak to git)
    # hcloud.token = "YOUR_TOKEN or set via HETZNERCLCOUD_TOKEN environment variable"

    ## SSH Key must exist and must match override.ssh.private_key_path
    hcloud.ssh_keys = ['develop']

    ## Select an explicit context in the config file
    ## If not speicified, the default one is used
    # hcloud.active_context = 'non_default_context'

    ## if you have a multi provider vagrantfile
    #  you can also override config.ssh.private_key_path
    #  provider specific.
    #override.ssh.private_key_path = './keys/testkey2'
  end
  config.vm.define 'blockbook-testing' do |blockbook|
    blockbook.trigger.after :up do |t|
      t.info = "rsync auto"
      t.run = {inline: "/usr/bin/screen -S blockbook-rsync-auto -dm vagrant rsync-auto blockbook-testing"}
      # If you want it running in the background switch these
      #t.run = {inline: "bash -c 'vagrant rsync-auto bork &'"}
    end
    blockbook.trigger.after :halt do |t|
      t.info = "rsync auto stop"
      t.run = {inline: "/usr/bin/screen -S blockbook-rsync-auto -X quit"}
      # If you want it running in the background switch these
      #t.run = {inline: "bash -c 'vagrant rsync-auto bork &'"}
    end
    config.vm.provider 'hetznercloud' do |hcloud|
      hcloud.name = 'blockbook-testing'
      hcloud.server_type = 'cx21'
      hcloud.image = 'ubuntu-18.04'
    end
    blockbook.vm.provision 'shell', inline: <<-SHELL
        set -e
        apt-get update && \
        apt-get install -y build-essential git wget pkg-config lxc-dev libzmq3-dev \
                           libgflags-dev libsnappy-dev zlib1g-dev libbz2-dev \
                           liblz4-dev graphviz && \
        apt-get clean

        export GOLANG_VERSION=go1.12.4.linux-amd64
        export ROCKSDB_VERSION=v5.18.3
        export GOPATH=/go
        export PATH=$PATH:$GOPATH/bin
        export CGO_CFLAGS="-I/opt/rocksdb/include"
        export CGO_LDFLAGS="-L/opt/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4"

        mkdir /build
        mkdir /out

        # install and configure go
        cd /opt && wget https://storage.googleapis.com/golang/$GOLANG_VERSION.tar.gz && \
        tar xf $GOLANG_VERSION.tar.gz
        ln -s /opt/go/bin/go /usr/bin/go
        mkdir -p $GOPATH
        echo -n "GO version: " && go version
        echo -n "GOPATH: " && echo $GOPATH

        # install rocksdb
        cd /opt && git clone -b $ROCKSDB_VERSION --depth 1 https://github.com/facebook/rocksdb.git
        cd /opt/rocksdb && CFLAGS=-fPIC CXXFLAGS=-fPIC make -j 4 release
        strip /opt/rocksdb/ldb /opt/rocksdb/sst_dump && \
        cp /opt/rocksdb/ldb /opt/rocksdb/sst_dump /build

        # install build tools
        go get github.com/golang/dep/cmd/dep
        go get github.com/gobuffalo/packr/...

        # download pre-loaded depencencies
        cleanup() { rm -rf $GOPATH/src/blockbook; } && \
        trap cleanup EXIT && \
        cd $GOPATH/src && \
        git clone https://github.com/trezor/blockbook.git && \
        cd blockbook && \
        dep ensure -vendor-only && \
        cp -r vendor /build/vendor
    SHELL
  end
end
