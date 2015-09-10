envdir
======

Not sure how to deploy your 12-factor app during early development?

`envdir` makes setting environment variables easy by reading files in a directory.
The filenames are the variable names.
The contents of the files are the variable values.

```
mkdir env
echo 'localhost:8080' > env/BIND_HOST
echo '10.1.1.3:28015' > env/DB_HOST
envdir ./env/ sh -c 'echo $BIND_HOST'  # prints localhost:8080
```


## Usage

```
Usage:
  envdir --version
  envdir --help
  envdir [-i] <directory> <command> [<arguments>...]

Arguments:
  <directory>  The directory of files representing environment variables.
  <command>    The command to run.
  <arguments>  The arguments of the command to run.

Options:
  -i, --ignore-environment  Start with an empty environment.
  --version  Show version.
  --help     Show help.

Interface:
  Each filename in <directory> is the name of an environment variable.
  The contents of the file is the value of the environment variable.
  The last newline of each file is ignored.
  If the file is empty (containing only 0 bytes or 1 newline),
    that environment variable is unset.

  envdir exits 111 if:
   * The directory's files can't be read
   * A filename contains "="
   * A file contains the null character
   * The command can't be run
```


## Why?

 * djb's envdir is removed from the Arch repositories because of burdensome licensing
 * this envdir ignores the last newline in a file, so you can:
    * Make environment files with `echo` instead of needing `echo -n`
    * Edit environment files with `vim`, which [adds `<EOL>` to the end of every line](https://stackoverflow.com/a/16224292)
 * Go may be overkill, but at least it is very easily distributable
 * It's a simple and practical Go exercise


## Compiling

Run `make`


## Notes

Basic program behavior is not likely to change in the future, but edge cases and error messages might.
