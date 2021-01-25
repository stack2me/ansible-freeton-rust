FROM ubuntu:18.04 as node_rust_build

RUN apt-get update && apt-get install -y \
    gpg \
    tar \
    cmake \
    build-essential \
    pkg-config \
    libssl-dev \
    libtool \
    m4 \
    automake \
    clang \
    git \
    curl

# RUN mkdir /root/tonlabs
RUN git clone https://github.com/tonlabs/ton-labs-node.git /root/ton-labs-node
WORKDIR /root/ton-labs-node
RUN git submodule init
RUN git submodule update
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

# install deps
ENV TZ=Europe/Moscow
ENV PATH="/root/.cargo/bin:${PATH}"
ENV RUST_BACKTRACE=1

RUN apt-get update && apt-get install -y  gnupg2 && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN cargo build --release

RUN git clone https://github.com/tonlabs/ton-labs-node-tools.git /root/ton-labs-node/ton-labs-node-tools

WORKDIR /root/ton-labs-node/ton-labs-node-tools
RUN git submodule init
RUN git submodule update
RUN cargo build --release


# build main image
FROM ubuntu:18.04 as node_rust

RUN apt-get update && apt-get install -y openssl && rm -rf /var/lib/apt/lists/*

RUN useradd -u 1090 -d /home/ton/ -m -g users -s /bin/bash ton

COPY --from=node_rust_build /root/ton-labs-node/ton-labs-node-tools/target/release/console /home/ton/console
COPY --from=node_rust_build /root/ton-labs-node/ton-labs-node-tools/target/release/keygen /home/ton/keygen
COPY --from=node_rust_build /root/ton-labs-node/target/release/ton_node /home/ton/ton_node_no_kafka
COPY entrypoint.sh /home/ton/entrypoint.sh
COPY gen_config.sh /home/ton/gen_config.sh

RUN mkdir /home/ton/configs /home/ton/logs /home/ton/node_db

RUN chown ton: /home/ton/*

USER ton
WORKDIR /home/ton

ENTRYPOINT ["sh", "entrypoint.sh"]