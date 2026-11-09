# Security Policy

## Scope
Fake DNS is meant for testing and mocking only. It answers every query with static data and must not be
used as a production DNS server or be reachable from untrusted networks. Use ``allowed_networks`` in the
config file to restrict access (see [README](README.md#access-restriction)).

## Reporting a vulnerability
Please do **not** open a public issue for security problems.
Report them privately via [GitHub Security Advisories](https://github.com/sguter90/fake-dns/security/advisories/new)
instead. You will get a response as soon as possible.

## Verifying releases
Every [release](https://github.com/sguter90/fake-dns/releases) contains a ``checksums.txt`` with the SHA256
hashes of all binaries and a detached GPG signature ``checksums.txt.asc``.
The signature is created automatically by the release workflow (``.github/workflows/build.yml``).

### 1. Import the release signing key
The public key is stored in [``release-key.asc``](release-key.asc):
```bash
curl -fsSL https://raw.githubusercontent.com/sguter90/fake-dns/main/release-key.asc | gpg --import
```

Make sure the fingerprint matches:
```
983C 0E12 97CE 6C94 4CED  3250 8D82 2938 3326 F7EF
```
```bash
gpg --fingerprint release@flying-lama.com
```

### 2. Verify the signature of the checksums
Download ``checksums.txt`` and ``checksums.txt.asc`` from the same release as your binary:
```bash
gpg --verify checksums.txt.asc checksums.txt
```
The output must contain ``Good signature from "fake-dns Release Signing <release@flying-lama.com>"``.
A warning that the key is not certified with a trusted signature is expected, as long as the fingerprint
matches the one above.

### 3. Verify the binary
```bash
sha256sum -c checksums.txt --ignore-missing
```
On macOS use ``shasum -a 256 -c checksums.txt --ignore-missing``.
The downloaded binary must be reported as ``OK``.

## Release signing key
Public key: [``release-key.asc``](release-key.asc)

| | |
|---|---|
| User ID | ``fake-dns Release Signing <release@flying-lama.com>`` |
| Fingerprint | ``983C 0E12 97CE 6C94 4CED  3250 8D82 2938 3326 F7EF`` |
| Algorithm | ed25519 (sign only) |
| Created | 2026-11-09 |
| Expires | 2028-11-08 |

If the key is ever compromised or replaced, it will be revoked and this file will be updated.
Releases signed with a revoked key should not be trusted.
