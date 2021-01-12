
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


# build main image
FROM ubuntu:18.04 as node_rust

RUN apt-get update && apt-get install -y openssl && rm -rf /var/lib/apt/lists/*

RUN useradd -u 1090 -d /home/ton/ -m -g users -s /bin/bash ton

COPY --from=node_rust_build /root/ton-labs-node/target/release/ton_node /home/ton/ton_node
COPY entrypoint.sh /home/ton/entrypoint.sh
RUN mkdir /home/ton/configs /home/ton/logs /home/ton/node_db


RUN chown ton: /home/ton/*

USER ton
WORKDIR /home/ton

ENTRYPOINT ["sh", "entrypoint.sh"]

