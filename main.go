// ⚡ Fast files encryption command-line tool and package.
//
// See [safelock-cli/safelock] for package references and examples, And the GitHub [repo] for updates.
//
// # Install
//
// For command-line
//
//	go install https://github.com/mrf345/safelock-cli@latest
//
// For packages
//
//	go get https://github.com/mrf345/safelock-cli@latest
//
// # Examples
//
// Encrypt a path with the default options
//
//	safelock-cli encrypt path_to_encrypt encrypted_file_path
//
// And to decrypt
//
//	safelock-cli decrypt encrypted_file_path decrypted_files_path
//
// If you want it to run silently with no interaction
//
//	echo "password123456" | safelock-cli encrypt path_to_encrypt encrypted_file_path --quiet
//
// [safelock-cli/safelock]: https://pkg.go.dev/github.com/mrf345/safelock-cli/safelock
// [repo]: https://github.com/mrf345/safelock-cli
package main

import "github.com/mrf345/safelock-cli/cmd"

func main() {
	cmd.Execute()
}
