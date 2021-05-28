#!/usr/bin/env python

import datetime
import json

import configargparse
import paho.mqtt.client as mqtt


def congratulations():
    colours = ["green", "red", "magenta", "yellow", "cyan"]
    lines = [f"{{{colour}}}congratulations!{{clear}}" for colour in colours]
    return "\n".join(lines)


def nosana():
    return "no life"


def on_connect(client, userdata, flags, rc):
    print(f"Connected with result code {rc}")

    client.subscribe("/gowon/input")


def on_message(client, userdata, msg):
    try:
        mj_in = json.loads(msg.payload.decode())
    except JSONDecodeError:
        print("Error parsing message json")
        return

    msg_in, dest = mj_in["msg"], mj_in["dest"]

    matches = [
        ("congratulations", congratulations),
        ("no sana", nosana),
    ]

    for r, f in matches:
        if r in msg_in.lower():
            out = f()
            msg_out = {"module": "module2", "msg": out, "dest": dest}
            client.publish("/gowon/output", json.dumps(msg_out))

            return


def main():
    print(datetime.datetime.now())
    print("gowon-module2 started")

    p = configargparse.ArgParser()
    p.add(
        "-H", "--broker-host", env_var="GOWON_BROKER_HOST", default="localhost"
    )
    p.add(
        "-P",
        "--broker-port",
        env_var="GOWON_BROKER_PORT",
        type=int,
        default=1883,
    )
    opts = p.parse_args()

    client = mqtt.Client("gowon-module2")

    client.on_connect = on_connect
    client.on_message = on_message

    client.connect(opts.broker_host, opts.broker_port)

    client.loop_forever()


if __name__ == "__main__":
    main()
