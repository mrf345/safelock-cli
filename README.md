<h2></h2>
<h1>
safelock-cli
<a href='https://github.com/mrf345/safelock-cli/actions/workflows/ci.yml'>
  <img src='https://github.com/mrf345/safelock-cli/actions/workflows/ci.yml/badge.svg' alt='build status'>
</a>
<a href='https://github.com/golangci/golangci-lint/tree/master'>
  <img src='https://img.shields.io/badge/linter-golangci--lint-blue.svg?logo=go&logoColor=white' alt='linter badge'>
</a>
<a href="https://pkg.go.dev/github.com/mrf345/safelock-cli/safelock">
  <img src='https://img.shields.io/badge/reference-blue.svg?logo=go&logoColor=white' alt='Go Reference'>
</a>
</h1>

Fast files encryption package and command-line tool built for speed with Go and [Archiver](https://github.com/mholt/archiver) ⚡

Utilizing `argon2id` and `chacha20poly1305` for encryption, see [default options](#options).


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


### Options

 Following the default options remanded by [RFC9106](https://datatracker.ietf.org/doc/html/rfc9106#section-7.4) and [crypto/argon2](https://pkg.go.dev/golang.org/x/crypto/argon2#IDKey)

| Option                  | Value                                       |
|-------------------------|---------------------------------------------|
| Iterations              | 3                                           |
| Memory size             | 64 Megabytes                                |
| Salt length             | 16                                          |
| Key length              | 32                                          |
| Threads                 | Number of available cores `runtime.NumCPU()`|
| Minimum password length | 8                                           |


### Performance

> [!NOTE]
> You can reproduce the results by running [bench_and_plot.py](benchmark/bench_and_plot.py) (based on [Matplotlib](https://github.com/matplotlib/matplotlib) and [Hyperfine](https://github.com/sharkdp/hyperfine))

<p align="center">
  <a href="https://raw.githubusercontent.com/mrf345/safelock-cli/master/benchmark/encryption-time.webp" target="_blank">
    <img src="benchmark/encryption-time.webp" alt="encryption time" />
  </a>
  <a href="https://raw.githubusercontent.com/mrf345/safelock-cli/master/benchmark/decryption-time.webp" target="_blank">
    <img src="benchmark/decryption-time.webp" alt="decryption time" />
  </a>
  <a href="https://raw.githubusercontent.com/mrf345/safelock-cli/master/benchmark/file-size.webp" target="_blank">
    <img src="benchmark/file-size.webp" alt="file size" />
  </a>
</p>
