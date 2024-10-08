import argparse
import sys

import multiaddr
import trio

from libp2p import new_host
from libp2p.network.stream.net_stream_interface import INetStream
from libp2p.peer.peerinfo import info_from_p2p_addr
from libp2p.typing import TProtocol

PROTOCOL_ID = TProtocol("/chat/1.0.0")
MAX_READ_LEN = 2**32 - 1


async def read_data(stream: INetStream) -> None:
    while True:
        read_bytes = await stream.read(MAX_READ_LEN)
        if read_bytes:
            read_string = read_bytes.decode()
            if read_string != "\n":
                # Green console color: \x1b[32m
                # Reset console color: \x1b[0m
                print("\x1b[32m%s\x1b[0m" % read_string, end="")


async def write_data(stream: INetStream) -> None:
    async_f = trio.wrap_file(sys.stdin)
    while True:
        line = await async_f.readline()
        await stream.write(line.encode())


async def run(port: int, destination: str, relay_ip: str, relay_peerid: str) -> None:
    # Listening address
    listen_addr = multiaddr.Multiaddr(f"/ip4/0.0.0.0/tcp/{port}")

    # Create the relay node's multiaddress
    relay_maddr_str = f"/ip4/{relay_ip}/tcp/4001/p2p/{relay_peerid}"
    relay_maddr = multiaddr.Multiaddr(relay_maddr_str)
    relay_info = info_from_p2p_addr(relay_maddr)

    # Create a new libp2p host
    host = new_host()

    async with host.run(listen_addrs=[listen_addr]), trio.open_nursery() as nursery:
        # Connect to the relay node
        await host.connect(relay_info)
        print(f"Connected to relay node at {relay_maddr_str}")

        # Add the relay address to our own addresses
        host.get_network().add_addrs(
            host.get_id(),
            [relay_maddr],
            ttl=10000,
        )

        if not destination:  # It's the server (listener)

            async def stream_handler(stream: INetStream) -> None:
                nursery.start_soon(read_data, stream)
                nursery.start_soon(write_data, stream)

            host.set_stream_handler(PROTOCOL_ID, stream_handler)

            # Show the multiaddress with relay
            # Our address will be the relay address encapsulated with our peer ID
            relay_listen_maddr = relay_maddr.encapsulate(
                multiaddr.Multiaddr(f"/p2p/{host.get_id().pretty()}")
            )
            print(
                f"Run this command on another console:\n\n"
                f"python chat.py -p {int(port) + 1} "
                f"-d '{relay_listen_maddr}' "
                f"--relay-ip {relay_ip} --relay-peerid {relay_peerid}\n"
            )
            print("Waiting for incoming connection...")

        else:  # It's the client (dialer)
            maddr = multiaddr.Multiaddr(destination)
            info = info_from_p2p_addr(maddr)

            # Connect to the relay node (already connected above)
            # Use the relay to connect to the destination peer
            await host.connect(info)

            # Start a stream with the destination peer
            stream = await host.new_stream(info.peer_id, [PROTOCOL_ID])

            nursery.start_soon(read_data, stream)
            nursery.start_soon(write_data, stream)
            print(f"Connected to peer {info.peer_id.pretty()}")

        await trio.sleep_forever()


def main() -> None:
    description = """
    This program demonstrates a simple p2p chat application using libp2p with a relay node.
    To use it, first run 'python chat.py -p <PORT> --relay-ip <RELAY_IP> --relay-peerid <RELAY_PEERID>',
    where <PORT> is the port number, <RELAY_IP> is the relay node's IP address, and <RELAY_PEERID> is its Peer ID.
    Then, run another host with 'python chat.py -p <ANOTHER_PORT> -d <DESTINATION> --relay-ip <RELAY_IP> --relay-peerid <RELAY_PEERID>',
    where <DESTINATION> is the multiaddress of the previous listener host.
    """
    parser = argparse.ArgumentParser(description=description)
    parser.add_argument(
        "-p", "--port", default=8000, type=int, help="Source port number"
    )
    parser.add_argument(
        "-d",
        "--destination",
        type=str,
        help="Destination multiaddr string",
    )
    parser.add_argument(
        "--relay-ip",
        type=str,
        required=True,
        help="IP address of the relay node",
    )
    parser.add_argument(
        "--relay-peerid",
        type=str,
        required=True,
        help="Peer ID of the relay node",
    )
    args = parser.parse_args()

    if not args.port:
        raise RuntimeError("Was not able to determine a local port")

    try:
        trio.run(
            run,
            args.port,
            args.destination,
            args.relay_ip,
            args.relay_peerid,
        )
    except KeyboardInterrupt:
        pass


if __name__ == "__main__":
    main()

