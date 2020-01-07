<div align="center">
<img src="https://raw.github.com/wiki/mpppk/everest/images/social_preview.png" title="logo">
</div>

everest is the static file server with no dependencies.

The most unique feature of everest is the `rebuild-with` command to import the specified files into the binary.  
By execute the command once, everest can serve imported files in a single binary.

<div align="center">
<img src="https://raw.github.com/wiki/mpppk/everest/images/terminal.png" title="logo">
</div>

## Installation

Download from [GitHub Releases](https://github.com/mpppk/everest/releases).

### From source

```bash
$ go get -d github.com/mpppk/everest
$ make install
```

## Usage
### serve static files

```bash
$ everest path/to/files
Files are served on http://localhost:3000
```

### Import files to binary

*Note: under the hood, `rebuild-with` command use `go build`, so you need to install Go compiler first.*

```bash
$ everest rebuild-with path/to/files
$ everest # if you execute everest with no arguments, imported files are served.
Embedded files are served on http://localhost:3000
```

## Credits

<div>Logo icon made by <a href="https://www.flaticon.com/authors/google" title="Google">Google</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a></div>
