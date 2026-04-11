---
title: "Spring Break, Rivulet, and An Ode to XMonad"
date: 2026-04-11
description: This post was originally made to the Haskell discourse.
---

### This post was originally posted to the [Haskell Discourse](https://discourse.haskell.org/t/ann-rivulet-for-wayland-prototype/13921)

Recently, I switched to Asahi Linux on my Mac. I've been a Linux user on my desktop PC for years, and for the most part I've always used Fedora. So I thought it would be neat to experiment with the Fedora Asahi Remix! However, it uses Wayland exclusively and the Asahi team has (rightfully) said that they don't want to work on X11 compatibility, and, frankly, good for them! The X Window System might be the most useful thing I used for years and years that I still hated every single time I had to deal with it in any capacity beyond being an end-user. So Wayland is good. But Wayland compositors work quite differently than the X server does, and I think it's been really hard for the Haskell community to make the switch, too, in recent years. It's no secret that we've always struggled a bit with GUI in Haskell, and Wayland bindings/compositing are no different lol.

The Problem of The Compositor in Haskell is that there's barely any good bindings to [`wlroots`](https://gitlab.freedesktop.org/wlroots/wlroots/) or [`smithay`](https://github.com/Smithay/smithay/)... the official [`hsroots`](https://github.com/swaywm/hsroots) bindings have been abandoned for ~8 years, Rust FFI in Haskell is sort of a dead end, so `smithay` is out of the picture. [`hayland`](https://hackage.haskell.org/package/hayland), [`sudbury`](https://github.com/abooij/sudbury) and [`wayland-wire`](https://github.com/sivertb/wayland-wire) have all been dormant for years. The sole active effort is [`wlhs`](https://github.com/bradrn/wlhs) (from 2024) that started as an outgrowth of XMonad, which, while very cool, development has been paused since the primary contributor was injured (I think, correct me if I'm wrong) and the bindings are still self-described as "highly incomplete." No standalone haskell bindings exist for `xdg-shell`, `xdg-decoration`, `wlr-foreign-toplevel-management` or the `wayland-protocols` collection. The only bright spot here is everything that's been automated by `haskell-gi`, like `gi-gtk4-layer-shell`, but it's deeply painful to use because you're basically writing a bunch of imperative C in Haskell.

There have also been a couple attempts at a Haskell Wayland compositor, beginning with [`waymonad`](https://github.com/waymonad/waymonad). But now it's "badly bitrotted" according to the XMonad team, and it's been abandoned for years. `tiny-wlhs` had a lot of interest, and it's still early-stage ... and I really hope it continues, but, y'know. [Woburn](https://github.com/sivertb/woburn) attempted a more "pure" approach, but it never got much working and has been mostly dormant for years.

## So What?

I recently read Issac Freund's article [Separating the Wayland Compositor and Window Manager](https://isaacfreund.com/blog/river-window-management/) and found it fascinating. [`river`](https://codeberg.org/river/river), a project I've followed for some time, has now become what he calls a "non-monolithic" Wayland compositor, allowing you to build a WM/"layout engine" in a lazy/garbage-collected language instead of C or Zig or whatever. 

I think baby steps are needed in all directions. XMonad, while wonderful, seemingly will remain X11-only for the foreseeable future, and there's just overall not been much done here. A truly Haskell-pure Wayland compositor is far away, I think. River is a great solution to this!

So over spring break, I started on building a window manager/layout engine in Haskell called Rivulet. I think DSL-style configs are great, and really wanted to use something like it on my laptop. For some time, I've been chained to `sway` (which is great, it works fine) but its config is annoying (to me) at best. I used to flip-flop back and forth between XMonad and `bspwm`, and so all the design decisions I've made so far on my own WM are informed by both of them. `river` makes it pretty easy to implement your own WM in a language like Haskell, as long as you implement the River protocol bindings. The `river-window-management-v1` protocol is great and has been exactly what I've been looking for ... River handles the systems-level work, like DRM/KMS, inputs, buffer management, frame timing, Xwayland, etc ... and the WM/Rivulet handles the layouts, keybinds, focus, and decorations via the protocol. It's all atomic & seperated by the IPC barrier!

## Rivulet

Right now, Rivulet implements a lot of things but it's still not really usable as a daily driver. Here's an unfinished list of things that work:

1. The layout engine works and you can write your own layouts with the Layout typeclass.
2. The manage/render cycle is implemented to a somewhat-minimum specification so that River doesn't just immediately flip out and drop the connection; although there are some things that it might still get angry about because I haven't implemented them
3. The config DSL is pretty much in its "final state" syntax-wise, it supports autostart & keybinding but a lot of the actions aren't wired up yet
4. Workspaces work and you can switch between them
5. Layer-shell protocol half-works (wallpapers, status bars, etc)

I haven't tested a lot of things with my current implementation (especially multi-monitor support)... off the top of my head for what DOESN'T work: fullscreening is probably broken, colored borders are broken, a lot of the DSL actions (like changing layouts or floating windows) doesn't work, a lot of things River wants (like hints sent from the WM to windows about capabilities and the such) don't work, resizing windows isn't a thing yet, and so on. TL;DR: it's early, and I'm working on it.

I have, however, forced myself to switch to using Rivulet exclusively as my WM until I feel that it's in a very usable state. Hopefully this should help me identify anything wrong while I'm using it as my daily driver, although I'd encourage you NOT to use it as a daily driver until, like, v1.0.0. The code is also kind of a mess, there's massive monster lambdas everywhere and it's a little hard to read. That is something I intend to change in the future as well just by working on it naturally.

My favorite thing about this project (and the reason I decided to write it) is that it lets you write your config in pure Haskell. You can do absurd things if you want, like binding a key to run Haskell's garbage collector. Or you can use unicode character operators in your config. Or you can be really Haskell-idiomatic and do stuff like this in the `keybinds` block:

```haskell
mapM_ (uncurry (~>)) $ zipWith (\k a -> ([Super] # k, a))
  ['1'..'9']
  (map focusWorkspace ["I","II","III","IV","V","VI","VII","VIII","IX"])
```

which binds Super + 1..9 to focusing workspaces I...IX. I have also provided some syntactic sugar for this, in the form of the `#*` operator:

```haskell
Control           #* ['1'..'9'] ~> focusWorkspace
Control <+> Shift #* ['1'..'9'] ~> sendToWorkspace
```

Hopefully over the next year or so Rivulet gets to be as fully-featured as something like XMonad or `bspwm`. I'm currently working on implementing the BinarySpacePartitioning layout and trying to fill out Rivulet to use the full `river-window-management-v1` protocol spec.

You can check out the project on [my GitHub](https://github.com/jackiedorland/rivulet) and try it out! Be forewarned, the build instructions for the dependencies are a little fried... you may need a modern version of the Zig compiler and you probably need to build River from source (versions >0.4.x aren't really in many distro repos yet.)

Also, here's an example config if you want to try it out (be forewarned, Rivulet copies an extremely simple example config to `~/.config/rivulet/Config.hs` if you don't write one yourself):

```haskell
import Rivulet

main :: IO ()
main = rivulet $ do

    gaps 20

    layout Tall

    monitor "eDP-1" ["1", "2", "3", "4", "5"] --change eDP-1 to whatever your Wayland output is called
    
    autostart $ do
        start "foot"

    keybinds $ do
        Super              #  Return       ~> spawn "foot"
        Super              #* ['1'..'5']   ~> focusWorkspace
        Super <+> Ctrl     #* ['1'..'5']   ~> sendToWorkspace
        Ctrl  <+> Alt      #  Delete       ~> exitSession
        Super              #  'q'          ~> closeFocused
        -- and so on
```

Let me know of your thoughts or if anything sucks about my design... I want to make this a genuinely useful WM. Cheers!