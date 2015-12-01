envrender - push the environment into config files
==================================================

[![Circle
CI](https://circleci.com/gh/teepark/envrender/tree/master.svg?style=svg)](https://circleci.com/gh/teepark/envrender/tree/master)

Some of us like to keep our application's [configuration in the
environment](http://12factor.net/config), but software we may rely on doesn't
always play by those rules. If you are already managing configuration
directories for use with [`envdir`](http://cr.yp.to/daemontools/envdir.html),
it would be nice to be able to stay in that paradigm even with software that
demands configuration files.

`envrender` runs files through a templating engine using the environment as
context, and then execs another program.


Examples
--------

To overwrite nginx's config file from a template before starting nginx:

    $ envrender ./nginx.conf.tmpl:/etc/nginx.conf nginx

To convert a configuration directory set up for envdir to a redis.conf file
before starting redis:

    $ envdir ./redis-conf.env \
      envrender ./redis-conf.tmpl:redis.conf \
      redis-server ./redis.conf


Arguments
---------

Any leading arguments to `envrender` containing a colon (`:`) specify a
template to render and its destination. Multiple such arguments can be provided
and they will all be run up until the first argument that has no colon.

Everything before the colon is the source, and everything after is the
destination (`source:destination`). They can be paths referencing files, or a
dash (`-`) for stdin (source) or stdout (destination). One caution about using
stdin: envrender will read it in its entirety before rendering the result, so
the stream is unavailable for the program being exec'd.

If either the source or destination is empty (the argument ends or starts with
a colon), the non-empty path will be used as both. So an argument of
`nginx.conf:` would expect `nginx.conf` to be a template, and would overwrite
it with the rendered result. Similarly, `:-` or `-:` will read a template from
stdin and render the result to stdout.

One final thing: you can use a double dash (`--`) argument as a separator
between renders and the program to exec. This is a useful way to be more
explicit generally, but also is necessary to run an executable program with a
colon in its path.


Templates
---------

The templating language is go's
[text/template](https://godoc.org/text/template). Environment variables are
used as the context.

There is one slight modification applied to the templates -- they support
`{{-` and `-}}` as tag delimiters. These variants eat any whitespace
immediately outside the tags. This behavior can be disabled with the flag
`-w=false` (before all other args).
