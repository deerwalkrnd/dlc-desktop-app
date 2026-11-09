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
| `Esc`               | Close the player               |

You can also:

- drag the progress bar to jump anywhere in the video,
- change the **playback speed** (0.5× to 2×) with the selector at the bottom
  right,
- double-click the video for full screen.

The app remembers your volume and speed.

When a lesson ends, you can **Watch again** or go to the **Next lesson**. If a
video can't be played (because the file was moved, deleted or damaged), the
player says so and offers **Try again**.

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
