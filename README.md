# simple-feed-gen

Tool for generating basic atom feeds from `.gmi` files. And nothing more.

## Install

### From source

clone & `make install`

### Nix

I packaged it in my NUR repo: [pnpkgs](https://github.com/pniedzwiedzinski/pnpkgs)

```
nix-shell -p nur.repos.pn.sfg
```

## Usage

```
sfg gemini://yoursite.abc/blog dir/with/gmi/files
```
