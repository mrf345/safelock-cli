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

Fast files encryption (AES-GCM) package and command-line tool built for speed with Go ⚡

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
safelock-cli decrypt path_to_encrypt encrypted_file_path
```
If you want it to run silently with no interaction

```shell
echo "password123456" | safelock-cli encrypt path_to_encrypt encrypted_file_path --quiet
```

<details>
  <summary>Simple example of using it within a package</summary>

  > Checkout [GoDocs](https://pkg.go.dev/github.com/mrf345/safelock-cli/safelock) for more examples and references

  ```go
  package main

  import "github.com/mrf345/safelock-cli/safelock"

  func main() {
    lock := safelock.New()
    inputPath := "/home/testing/important"
    outputPath := "/home/testing/encrypted.sla"
    password := "testing123456"

    // Encrypts `inputPath` with the default settings
    if err := lock.Encrypt(nil, inputPath, outputPath, password); err != nil {
      panic(err)
    }

    // Decrypts `outputPath` with the default settings
    if err := lock.Decrypt(nil, outputPath, "/home/testing", password); err != nil {
      panic(err)
    }
  }
  ```
</details>

### Performance

With the default settings it should be **twice** as fast as `gpgtar`

```shell
> du -hs testing/
1.2G testing/

> time gpgtar --encrypt --output testing.gpg -r user testing/
real	0m42.710s
user	0m41.148s
sys	0m7.943s

> time echo "testing123456" | safelock-cli encrypt testing/ testing.sla --quiet
real	0m22.902s
user	0m29.171s
sys	0m10.868s
```
You can get faster performance using the `--sha256` flag (considered less secure)

```shell
> time echo "testing123456" | safelock-cli encrypt testing/ testing.sla --quiet --sha256
real	0m18.843s
user	0m20.619s
sys	0m10.901s
```

And no major file size difference

```shell
> ls -lh --block-size=MB testing.gpg
-rw-rw-r-- 1 user user 1247MB Aug 10 12:15 testing.gpg

> ls -lh --block-size=MB testing.sla
-rw-rw-r-- 1 user user 1273MB Aug 10 11:30 testing.sla
```
