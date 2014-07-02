Assets
======

Magically bundle (some) assets into your Golang binaries.

Uses a combination of boring function calls and a (horrible) pre-processor to
generate a magical assets file for your project.

**When does it work?** When you use the very small and limited setup that this
was tested in and designed for.
**When does it _NOT_ work?** When you don't live on Planet Sysrq.

Definitely don't use this when you're working with very large files - really,
you shouldn't want to embed those in the first place. This means trying to embed
a copy of the world's greatest film (_Hackers (1995)_) in your app is probably
not a good idea.

