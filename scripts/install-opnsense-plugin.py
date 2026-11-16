#!/usr/bin/env python3
"""Install an OPNsense plugin and wait until its API is ready.

Prefers QEMU guest-agent `pkg install` when QEMU_GA_SOCKET is available
(CI), and falls back to the firmware API.

Usage:
  install-opnsense-plugin.py os-acme-client

Environment:
  OPNSENSE_URI, OPNSENSE_API_KEY, OPNSENSE_API_SECRET
  QEMU_GA_SOCKET (optional, default /tmp/qemu-virtserialport.sock)
"""

from __future__ import annotations

import base64
import json
import os
import select
import socket
import ssl
import sys
import time
import urllib.error
import urllib.request
from base64 import b64encode

QEMU_GA_SOCKET = os.environ.get("QEMU_GA_SOCKET", "/tmp/qemu-virtserialport.sock")

PROBES = {
    # Prefer POST search: GET without a body often returns "Invalid JSON syntax".
    "os-acme-client": ("POST", "/api/acmeclient/accounts/search"),
}


def api(method: str, path: str, body: dict | None = None) -> tuple[int, dict | list | str]:
    base = os.environ["OPNSENSE_URI"].rstrip("/")
    key = os.environ["OPNSENSE_API_KEY"]
    secret = os.environ["OPNSENSE_API_SECRET"]
    credentials = b64encode(f"{key}:{secret}".encode()).decode()
    headers = {
        "Authorization": f"Basic {credentials}",
        "Accept": "application/json",
    }
    ctx = ssl.create_default_context()
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE
    data = None
    if method.upper() != "GET":
        payload = {} if body is None else body
        data = json.dumps(payload).encode()
        headers["Content-Type"] = "application/json"
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


def plugin_ready(package: str) -> bool:
    probe = PROBES.get(package)
    if not probe:
        return False
    method, path = probe
    code, _ = api(method, path, {} if method.upper() != "GET" else None)
    return code == 200


def wait_plugin_ready(package: str, timeout_s: int = 180) -> None:
    probe = PROBES.get(package)
    if not probe:
        return
    method, path = probe
    deadline = time.time() + timeout_s
    while time.time() < deadline:
        if plugin_ready(package):
            print(f"  probe {method} {path} OK", flush=True)
            return
        code, payload = api(method, path, {} if method.upper() != "GET" else None)
        print(f"  probe {method} {path} HTTP {code}: {payload!r}", flush=True)
        time.sleep(3)
    raise RuntimeError(f"plugin API not ready: {method} {path}")


def qemu_ga(command: dict, timeout: float = 30) -> dict:
    with socket.socket(socket.AF_UNIX, socket.SOCK_STREAM) as client:
        client.settimeout(timeout)
        client.connect(QEMU_GA_SOCKET)
        client.sendall((json.dumps(command) + "\n").encode())
        response = b""
        while True:
            ready, _, _ = select.select([client], [], [], timeout)
            if not ready:
                raise TimeoutError("qemu-ga response timeout")
            chunk = client.recv(8192)
            if not chunk:
                break
            response += chunk
            if b"\n" in chunk:
                break
    return json.loads(response.decode())


def install_via_qemu_ga(package: str) -> None:
    print(f"Installing {package} via qemu-ga pkg ...", flush=True)
    started = qemu_ga(
        {
            "execute": "guest-exec",
            "arguments": {
                "path": "/usr/sbin/pkg",
                "arg": ["install", "-y", package],
                "capture-output": True,
            },
        }
    )
    pid = started.get("return", {}).get("pid")
    if not pid:
        raise RuntimeError(f"guest-exec failed: {started}")

    deadline = time.time() + 600
    while time.time() < deadline:
        status = qemu_ga({"execute": "guest-exec-status", "arguments": {"pid": pid}})
        ret = status.get("return") or {}
        if not ret.get("exited"):
            time.sleep(2)
            continue
        exitcode = ret.get("exitcode")
        out = base64.b64decode(ret.get("out-data") or "").decode(errors="replace")
        err = base64.b64decode(ret.get("err-data") or "").decode(errors="replace")
        if out.strip():
            print(out.strip(), flush=True)
        if err.strip():
            print(err.strip(), flush=True)
        if exitcode not in (0, None):
            # pkg may return non-zero if already installed; still probe API.
            print(f"  pkg exitcode={exitcode}", flush=True)
        return
    raise RuntimeError("timed out waiting for pkg install via qemu-ga")


def wait_idle(timeout_s: int = 600, *, require_busy_first: bool = False) -> None:
    deadline = time.time() + timeout_s
    saw_busy = not require_busy_first
    grace_deadline = time.time() + 15 if require_busy_first else 0.0

    while time.time() < deadline:
        code, payload = api("POST", "/api/core/firmware/upgradestatus", {})
        if code != 200:
            raise RuntimeError(f"upgradestatus HTTP {code}: {payload}")
        if not isinstance(payload, dict):
            raise RuntimeError(f"unexpected upgradestatus payload: {payload!r}")
        status = str(payload.get("status", ""))
        print(f"  firmware status={status!r}", flush=True)

        if status not in ("done", "error", ""):
            saw_busy = True

        if status == "error":
            raise RuntimeError(f"firmware reported error: {payload}")

        if status == "done":
            if saw_busy or time.time() >= grace_deadline:
                return
            time.sleep(1)
            continue

        time.sleep(3)

    raise RuntimeError("timed out waiting for firmware idle")


def install_via_firmware_api(package: str) -> None:
    print(f"Installing {package} via firmware API ...", flush=True)
    code, payload = api("POST", f"/api/core/firmware/install/{package}", {})
    if code != 200:
        raise RuntimeError(f"install HTTP {code}: {payload}")
    if isinstance(payload, dict) and payload.get("status") not in ("ok", "done"):
        print(f"  install response: {payload}", flush=True)

    time.sleep(2)
    wait_idle(require_busy_first=True)


def main() -> int:
    if len(sys.argv) != 2:
        print("usage: install-opnsense-plugin.py <package>", file=sys.stderr)
        return 2
    for var in ("OPNSENSE_URI", "OPNSENSE_API_KEY", "OPNSENSE_API_SECRET"):
        if not os.environ.get(var):
            print(f"ERROR: {var} must be set", file=sys.stderr)
            return 1

    package = sys.argv[1]
    if plugin_ready(package):
        probe = PROBES[package]
        print(f"{package} already available ({probe[0]} {probe[1]}).", flush=True)
        return 0

    qemu_available = os.path.exists(QEMU_GA_SOCKET)
    errors: list[str] = []

    if qemu_available:
        try:
            install_via_qemu_ga(package)
        except Exception as exc:  # noqa: BLE001
            errors.append(f"qemu-ga: {exc}")
            print(f"WARNING: qemu-ga install failed: {exc}", flush=True)

    if not plugin_ready(package):
        try:
            install_via_firmware_api(package)
        except Exception as exc:  # noqa: BLE001
            errors.append(f"firmware: {exc}")
            print(f"WARNING: firmware install failed: {exc}", flush=True)

    wait_plugin_ready(package)
    print(f"Installed {package}.", flush=True)
    if errors:
        print(f"  notes: {'; '.join(errors)}", flush=True)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:  # noqa: BLE001 - surface in CI logs
        print(f"ERROR: {exc}", file=sys.stderr)
        raise SystemExit(1) from exc
