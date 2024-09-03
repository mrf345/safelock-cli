<h2></h2>
<h1>
safelock-cli
<a href='https://github.com/mrf345/safelock-cli/actions/workflows/ci.yml'>
  <img src='https://github.com/mrf345/safelock-cli/actions/workflows/ci.yml/badge.svg'>
</a>
<a href="https://pkg.go.dev/github.com/mrf345/safelock-cli/safelock">
  <img src="https://pkg.go.dev/badge/github.com/mrf345/safelock-cli/.svg" alt="Go Reference">
</a>
</h1>

Fast files encryption (AES-GCM) package and command-line tool built for speed with Go and [Archiver](https://github.com/mholt/archiver) ⚡

### Install

For command-line

```shell
go install https://github.com/mrf345/safelock-cli@latest
```

For packages

```shell
go get https://github.com/mrf345/safelock-cli@latest
```

Or using one of the latest release binaries [here](https://github.com/mrf345/safelock-cli/releases)

### Examples

Encrypt a path with default options

```shell
safelock-cli encrypt path_to_encrypt encrypted_file_path
```
And to decrypt

```shell
safelock-cli decrypt encrypted_file_path decrypted_files_path
```
> [!TIP]
> If you want it to run silently with no interaction use `--quiet` and pipe the password

```shell
echo "password123456" | safelock-cli encrypt path_to_encrypt encrypted_file_path --quiet
```

You can find interactive examples of using it as a package to [encrypt](https://pkg.go.dev/github.com/mrf345/safelock-cli/safelock#example-Safelock.Encrypt) and [decrypt](https://pkg.go.dev/github.com/mrf345/safelock-cli/safelock#example-Safelock.Decrypt).

### Performance

With the default settings it should be about **19.1** times faster than `gpgtar`

> [!NOTE]
> You can reproduce the results by running [bench_and_plot.py](benchmark/bench_and_plot.py) (based on [matplotlib](https://github.com/matplotlib/matplotlib) and [hyperfine](https://github.com/sharkdp/hyperfine))

<p align="center">
  <a href="https://raw.githubusercontent.com/mrf345/safelock-cli/master/benchmark/encryption-time.webp" target="_blank">
    <img src="benchmark/encryption-time.webp" align="center" alt="encryption time" />
  </a>
  <a href="https://raw.githubusercontent.com/mrf345/safelock-cli/master/benchmark/decryption-time.webp" target="_blank">
    <img src="benchmark/decryption-time.webp" align="center" alt="encryption time" />
  </a>
</p>

And you could gain a slight file size reduction

```shell
> du -hs test/
1.2G test/

> ls -lh --block-size=MB test.sla test.gpg
-rw-r--r-- 1 mrf3 mrf3 1.2G Sep  3 17:55 test.gpg
-rw-r--r-- 1 mrf3 mrf3 959M Sep  3 17:29 test.sla
```
