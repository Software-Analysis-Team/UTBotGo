<!--[
# {PROJECT_NAME}
|-->
# UTBotGo
<!--]-->

Program for generating unit tests for Go files.

## How to install

<!--$ command += make install -->
<!--[
```bash
{command}
```
|-->
```bash
make install
```
<!--]-->

<!--[
This command installs program `{PROGRAM_NAME}` on your system, so you will be
able to run the shell command `{PROGRAM_NAME}` from any directory.
|-->
This command installs program `utbotgo` on your system, so you will be
able to run the shell command `utbotgo` from any directory.
<!--]-->

> **NOTE**: You may need to run this command together with `sudo`.

If you want just build this program without installation, you can type

<!--[
```bash
make build
```
|-->
```bash
make build
```
<!--]-->

## How to use

<!--[
```bash
{PROGRAM_NAME} <option> <file_or_dir>
```
|-->
```bash
utbotgo <option> <file_or_dir>
```
<!--]-->

### Option

* `--gofuzz` - use package [gofuzz](https://github.com/google/gofuzz) for
fuzzing

* `--go-fuzz` - use package [go-fuzz](https://github.com/dvyukov/go-fuzz) for
fuzzing

### File or dir

It is name of Go file or directory with Go files, for  which you want to
generate testing files (files `*_test.go`).

### Output

Output will contain amount of generated files.

For example, if you have Go file `hello.go` and run
<!--$ command += {PROGRAM_NAME} --gofuzz hello.go -->
<!--[
```bash
{command}
```
|-->
```bash
utbotgo --gofuzz hello.go
```
<!--]-->
testing file `hello_test.go` will be generated, and output will be `1` (or `0`
if error will be thrown).

## How to uninstall

<!--$ command += make uninstall -->
<!--[
```bash
{command}
```
|-->
```bash
make uninstall
```
<!--]-->

<!--[
It removes command `{PROGRAM_NAME}` from your system.
|-->
It removes command `utbotgo` from your system.
<!--]-->

> **Note**: You may need to run this command together with `sudo`.

## Tests

<!--[
```
{$ {command}|
}
```
|-->
```
$ make install
$ utbotgo --gofuzz hello.go
$ make uninstall
```
<!--]-->
