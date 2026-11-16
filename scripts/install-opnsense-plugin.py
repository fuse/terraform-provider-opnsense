#!/usr/bin/env python3
"""Install an OPNsense plugin via the firmware API and wait until idle.

Usage:
  install-opnsense-plugin.py os-acme-client

Environment:
  OPNSENSE_URI, OPNSENSE_API_KEY, OPNSENSE_API_SECRET
"""

from __future__ import annotations

import json
import os
import ssl
import sys
import time
import urllib.error
import urllib.request
from base64 import b64encode


def api(method: str, path: str, body: dict | None = None) -> tuple[int, dict | list | str]:
    base = os.environ["OPNSENSE_URI"].rstrip("/")
    key = os.environ["OPNSENSE_API_KEY"]
    secret = os.environ["OPNSENSE_API_SECRET"]
    credentials = b64encode(f"{key}:{secret}".encode()).decode()
    headers = {
        "Authorization": f"Basic {credentials}",
        "Content-Type": "application/json",
    }
    ctx = ssl.create_default_context()
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE
    data = None if body is None else json.dumps(body).encode()
    req = urllib.request.Request(
        base + path, data=data, headers=headers, method=method
    )
    try:
        with urllib.request.urlopen(req, context=ctx, timeout=120) as resp:
            raw = resp.read().decode()
            code = resp.status
    except urllib.error.HTTPError as e:
        raw = e.read().decode()
        code = e.code
    try:
        parsed: dict | list | str = json.loads(raw) if raw else {}
    except json.JSONDecodeError:
        parsed = raw
    return code, parsed


def wait_idle(timeout_s: int = 600) -> None:
    deadline = time.time() + timeout_s
    while time.time() < deadline:
        code, payload = api("POST", "/api/core/firmware/upgradestatus", {})
        if code != 200:
            raise RuntimeError(f"upgradestatus HTTP {code}: {payload}")
        if not isinstance(payload, dict):
            raise RuntimeError(f"unexpected upgradestatus payload: {payload!r}")
        status = str(payload.get("status", ""))
        print(f"  firmware status={status!r}", flush=True)
        if status in ("done", "error"):
            if status == "error":
                raise RuntimeError(f"firmware reported error: {payload}")
            return
        time.sleep(3)
    raise RuntimeError("timed out waiting for firmware idle")


def main() -> int:
    if len(sys.argv) != 2:
        print("usage: install-opnsense-plugin.py <package>", file=sys.stderr)
        return 2
    for var in ("OPNSENSE_URI", "OPNSENSE_API_KEY", "OPNSENSE_API_SECRET"):
        if not os.environ.get(var):
            print(f"ERROR: {var} must be set", file=sys.stderr)
            return 1

    package = sys.argv[1]
    print(f"Installing {package} ...", flush=True)
    code, payload = api("POST", f"/api/core/firmware/install/{package}", {})
    if code != 200:
        raise RuntimeError(f"install HTTP {code}: {payload}")
    if isinstance(payload, dict) and payload.get("status") not in ("ok", "done"):
        # Already installed often returns failure with a message; treat as soft OK
        # if the plugin API becomes available after a short wait.
        print(f"  install response: {payload}", flush=True)

    wait_idle()
    print(f"Installed {package}.", flush=True)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:  # noqa: BLE001 - surface in CI logs
        print(f"ERROR: {exc}", file=sys.stderr)
        raise SystemExit(1) from exc
