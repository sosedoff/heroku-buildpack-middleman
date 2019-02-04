# Middleman BuildPack

This is a buildpack for [Middleman](http://middlemanapp.com) that will
create your static site.

It uses a simple HTTP server written in Go (see server.go) to serve all static
pages and assets instead of standard `middleman server` command. There's also support
for basic authentication with user/password and a handler for missing pages.

## Usage

Add the buildpack to your Heroku project:

```
heroku buildpacks:add https://github.com/sosedoff/heroku-buildpack-middleman
```

### Root Directory

By default application will try to serve files from `build` directory. You can change
that with `STATIC_DIR` environment variable:

```
heroku config:set STATIC_DIR=my-assets
```

### Authentication

To require basic authentication set the following environment variables (both are required):

```
heroku config:set AUTH_USER=admin AUTH_PASSWORD=password
```

### 404 Pages

If you want to customize the default 404 page (which is plain text), add the following var:

```
heroku config:set NOT_FOUND_PATH=/404.html
```

Make sure you have `404.html` file in your static assets directory.