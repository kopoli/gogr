# gogr

A tool for running commands in multiple directories. There are a lot like it
and this one is mine.

## Installation

```
$ go get github.com/kopoli/gogr
```

## Description

The idea was to have a quick-and-dirty Go implementation of the
[gr tool](http://mixu.net/gr/).  The distinguishing feature is that this
requires no runtime dependencies and is able to run the commands concurrently.

It both supports the @tag idea and giving directories directly to the program. 

## Simple usage

```
# Add directories to tag @this

$ gogr tag add this . .. ../..

# run command in each directory

$ gogr @this ls -l

```

## Manual

### Concepts

- Running a single command in a group of directories.
- Tagging a group of directories under a single name.
- Setting a tag by discovering directories which contain a certain file.

### Running commands

Commands can be run by giving a list of directories and/or tags and the
command to run. Example:

```
$ gogr @projects ../extra/repository @src git status -sb
```

This will run `git status -sb` in directories that are in tags `@projects` and
`@src` and in directory `../extra/repository`.

Note: If a directory (or tag) is present multiple times the command will be
run at most one time in a directory.

By giving the `-j` flag the command is run in parallel in all directories:

```
$ gogr -j @src git remote update
```

See `gogr --help` for more information.

### Tagging

You can create and remove tags. Tags consist of directories which can be added
and removed. Two possible syntaxes exist.

#### Creating a tag

```
$ gogr tag add this .

# Alternatively

$ gogr +@this .
```

See `gogr tag add --help` for more information.

### Removing a tag

The `gogr tag rm` if given no arguments will remove the tag completely. If it
is given directories which are tagged with the tag, it will untag the directories.

```
# Remove the whole tag
$ gogr tag rm this


# Remove directories from a tag

$ gogr tag rm this .

# Alternatively

$ gogr -- -@this .
```

NOTE: The `-@this` syntax is ambiguous with command line flags, it must be
preceded by `--`.

### Listing tags

The created tags can be viewed with `gogr tag list`. The directories in a tag
can be viewed with `gogr tag list tagname`.


### Creating tags by discovery

Tags can be created by walking through the directory tree and tagging
directories which contain a given file. The following will tag all directories
containing a file or a directory called `.git` in the `~/src` directory and
subdirectories:

```
$ gogr discover src ~/src
```

The default maximum depth of directory tree is 5 levels. It can be changed
with the `-d` flag. The name of the file that is searched can be changed with
the `-f` -flag:

```
$ gogr discover -f README readmes ~/src/go ~/projects
```

For the file name, globbing is not supported; it needs a complete file
name. More information can be found via `gogr discover --help`.

## License

MIT license
