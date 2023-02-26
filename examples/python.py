#!/usr/bin/env python3
import sys
import json
import socket


def parse_header() -> dict:
    return json.load(sys.stdin)


def result(data: dict, exc: Exception = None):
    print(
        json.dumps(
            {
                "status": int(exc == None),
                "reason": str(exc) if exc else "ok",
                "payload": data,
            }
        )
    )


def main():
    header = parse_header()
    print("stdin parsed", json.dumps(header, indent=3))
    cmd = header["data"]["cmd"]
    if cmd == "hostname":
        return result({"hostname": socket.gethostname()})
    print("unknown command", cmd)
    return result({"foo": "bar"})


INPUT_SCHEMA = {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "properties": {
        "user": {"type": "string"},
        "cmd": {"type": "string"},
        "requestId": {"type": "string"},
    },
    "required": ["user", "cmd", "requestId"],
}
OUTPUT_SCHEMA = {
    "$schema": "http://json-schema.org/draft-07/schema#",
    "type": "object",
    "properties": {
        "status": {"type": "integer"},
        "reason": {"type": "string"},
        "payload": {},
    },
    "required": ["status", "reason", "payload"],
}


def schema():
    _schema = {
        "name": "python",
        "description": "Python example",
        "inputs": INPUT_SCHEMA,
        "outputs": OUTPUT_SCHEMA,
    }
    print(json.dumps(_schema, indent=3))


if __name__ == "__main__":
    if len(sys.argv) > 1 and sys.argv[1] == "schema":
        schema()
        sys.exit(0)
    main()
    print("Not captured")
