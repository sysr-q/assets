Assets
======

Magically bundle (some) assets into your Golang binaries.

Uses a combination of boring function calls and a (~~horrible~~ genius)
pre-processor to generate a magical binary with pre-bundled assets.

**When does it work?** When you use the very small and limited setup that this
was tested in and designed for.
**When does it _NOT_ work?** When you don't live on Planet Sysrq.

Definitely don't use this when you're working with very large files - really,
you shouldn't want to embed those in the first place. This means trying to embed
a copy of the world's greatest film (_Hackers (1995)_) in your app is probably
~~not a good idea~~ the greatest example ever made. (see: `example/`)

The gist is this:

1. `go get github.com/sysr-q/assets`
2. `go install github.com/sysr-q/assets/goappc`
3. `import "github.com/sysr-q/assets"`
4. `foo := string(assets.MustRead("templates/blah.html"))` (for example)
5. `goappc .` (runs magical preprocessor, bundles assets, runs _go build_)
6. Deploy!

**NOTE!** `goappc` is still not complete - currently it will only process a
_single_ file. Not really what you want, right? I'll remedy this sometime soon.
