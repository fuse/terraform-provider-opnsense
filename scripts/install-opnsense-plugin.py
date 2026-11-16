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
        "Accept": "application/json",
    }
    ctx = ssl.create_default_context()
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE
    # POST/PUT: always send JSON (OPNsense rejects body-less POSTs with
    # "Invalid JSON syntax"). GET: no body.
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


def wait_idle(timeout_s: int = 600, *, require_busy_first: bool = False) -> None:
    """Poll firmware upgradestatus until status is done.

    When require_busy_first is True, keep polling until a non-done status is
    observed (or a short grace window elapses), so a stale "done" from before
    install does not make us return before the new job starts.
    """
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
            # Stale done from before install started; keep waiting briefly.
            time.sleep(1)
            continue

        time.sleep(3)

    raise RuntimeError("timed out waiting for firmware idle")


def plugin_api_ready(method: str, path: str, timeout_s: int = 120) -> None:
    deadline = time.time() + timeout_s
    while time.time() < deadline:
        code, payload = api(method, path, {} if method.upper() != "GET" else None)
        if code == 200:
            print(f"  probe {method} {path} OK", flush=True)
            return
        print(f"  probe {method} {path} HTTP {code}: {payload!r}", flush=True)
        time.sleep(3)
    raise RuntimeError(f"plugin API not ready: {method} {path}")


def main() -> int:
    if len(sys.argv) != 2:
        print("usage: install-opnsense-plugin.py <package>", file=sys.stderr)
        return 2
    for var in ("OPNSENSE_URI", "OPNSENSE_API_KEY", "OPNSENSE_API_SECRET"):
        if not os.environ.get(var):
            print(f"ERROR: {var} must be set", file=sys.stderr)
            return 1

    package = sys.argv[1]
    # Prefer POST search endpoints: GET without a body often returns
    # "Invalid JSON syntax" on OPNsense plugin controllers.
    probes = {
        "os-acme-client": ("POST", "/api/acmeclient/accounts/search"),
    }

    # Already present?
    probe = probes.get(package)
    if probe:
        method, path = probe
        code, _ = api(method, path, {})
        if code == 200:
            print(f"{package} already available ({method} {path}).", flush=True)
            return 0

    print(f"Installing {package} ...", flush=True)
    code, payload = api("POST", f"/api/core/firmware/install/{package}", {})
    if code != 200:
        raise RuntimeError(f"install HTTP {code}: {payload}")
    if isinstance(payload, dict) and payload.get("status") not in ("ok", "done"):
        # Already installed often returns failure with a message; still wait and probe.
        print(f"  install response: {payload}", flush=True)

    # Give the install job a moment to flip upgradestatus away from a stale "done".
    time.sleep(2)
    wait_idle(require_busy_first=True)

    if probe:
        method, path = probe
        plugin_api_ready(method, path)

    print(f"Installed {package}.", flush=True)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:  # noqa: BLE001 - surface in CI logs
        print(f"ERROR: {exc}", file=sys.stderr)
        raise SystemExit(1) from exc
