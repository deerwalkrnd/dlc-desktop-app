---
sidebar_position: 3
title: Using the app
---

# Using the app

## Finding a lesson

1. **Classes**: the home page shows every class that has at least one video.
2. **Subjects**: pick a class. If it has both old- and new-syllabus content,
   switch with the **New syllabus / Old syllabus** tabs.
3. **Chapters and lessons**: pick a subject. Lessons are grouped by chapter
   in lesson-number order. The sidebar switches between subjects of the same
   class.
4. Click a lesson to play it.

The **Teachers** page lists every teacher with their lesson count and
subjects.

## Searching

Type in the **Search lessons** box at the top of every page, or press `/` to
jump to it. Search looks at lesson titles, chapters, subjects, teachers,
class numbers (`class 8` or just `8`) and lesson numbers (`1.2`). Every word
you type has to match somewhere, so `class 9 light` finds the Light lessons
in Class 9. Results update as you type, and clicking one starts the lesson.

## Picking up where you left off

The app remembers how far you got in each lesson:

- Lessons you've finished show a **✓ Watched** mark, and each chapter shows
  how many of its lessons you've watched.
- Lessons you've started show a progress bar and **Resume at 4:32**.
  Opening one continues from there, and a **Start over** button appears for
  a few seconds.
- The home page lists recent unfinished lessons under **Continue watching**.
  Hover over one and press × to remove it.

A lesson counts as watched once you reach the end, or 95% of the way
through (or the last 15 seconds of a longer video).

Progress is saved **in this computer's browser**. It isn't shared between
computers or browsers, and clearing the browser's data resets it.

## The video player

The player opens on top of the lesson list and starts playing. Use the
buttons along the bottom, or the keyboard:

| Key                 | Action                         |
| ------------------- | ------------------------------ |
| `Space` or `K`      | Play / pause                   |
| `←` / `J`           | Back 10 seconds                |
| `→` / `L`           | Forward 10 seconds             |
| `↑` / `↓`           | Volume up / down               |
| `M`                 | Mute / unmute                  |
| `F`                 | Full screen                    |
| `Shift` + `N` / `P` | Next / previous lesson         |
| `/`                 | Search (when the player is closed) |
| `Esc`               | Close the player               |

You can also:

- drag the progress bar to jump anywhere in the video,
- change the **playback speed** (0.5× to 2×) with the selector at the bottom
  right,
- double-click the video for full screen.

The app remembers your volume and speed.

Below the video, **Previous** and **Next** show which lessons come before
and after. When a lesson ends, you can **Watch again** or go to the **Next
lesson**. If a video can't be played (because the file was moved, deleted or
damaged), the player says so and offers **Try again**.

## The Library page

The **Library** page shows:

- the folder the app reads videos from,
- how many lessons it found and when it last looked,
- skipped files and warnings (see [Adding videos](./adding-videos.md#checking-your-files)),
- the **Rescan library** button, which picks up new, renamed or deleted
  videos straight away.

## Sharing with other computers on the network

By default only the computer running the app can open it. To let other
computers in the same lab watch from it, start the app from a Command Prompt
with:

```
dlc-desktop-app.exe -addr 0.0.0.0:3000
```

Other computers can then open `http://<this-computer's-IP>:3000`. Windows
Firewall will ask for permission the first time. Only do this on a trusted
network: anyone who can reach the computer can watch the videos.
