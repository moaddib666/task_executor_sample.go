#!/usr/bin/env python3
import sys
import json

def parse_header() -> dict:
    return json.load(sys.stdin)

def result(data: dict, exc: Exception = None):
    print(
        json.dumps(
            {"status": int(exc == None), "reason": str(exc) if exc else "ok", "payload": data}
        )
    )

def main():
    header = parse_header()
    print("stdin parsed", json.dumps(header, indent=3))
    print("hello")
    return result({"foo":"bar"})

if __name__ == "__main__":
    main()
    print("Not captured")
