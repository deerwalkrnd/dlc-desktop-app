---
sidebar_position: 1
title: Getting started
slug: /intro
---

# Getting started

The DLC Desktop App plays Deerwalk Learning Center lesson videos **without an
internet connection**. It runs on one computer, finds the lesson videos in a
folder and shows them, grouped by class, subject and chapter, in the web
browser.

## Requirements

- Windows 10 or later (64-bit). Windows 7 and 8 are not supported.
- A modern web browser: Chrome, Edge or Firefox.
- Disk space for your videos. The app itself is under 10 MB.

## Install

1. Download the latest `dlc-desktop-app-…-windows-amd64.zip` from the
   [releases page](https://github.com/deerwalkrnd/dlc-desktop-app/releases).
2. Extract the zip anywhere, for example `C:\DLC Desktop App`. You get:

   ```
   dlc-desktop-app-…-windows-amd64
   ├── dlc-desktop-app.exe   ← the app
   ├── dlc-check.exe         ← checks video file names (optional)
   ├── DLC                   ← put your videos here
   └── README.txt
   ```

   The web interface is built into `dlc-desktop-app.exe`. There is no separate
   `web` folder any more.

3. Copy your lesson videos into the `DLC` folder. Sub-folders are fine. Each
   file name has to follow the [naming format](./adding-videos.md).

## Run

Double-click `dlc-desktop-app.exe`. A console window opens and your browser
shows the app at [http://localhost:3000](http://localhost:3000).

- **Keep the console window open** while you use the app. Closing it stops
  the app.
- If the browser does not open by itself, type `http://localhost:3000` into
  the address bar.
- Starting the app again while it is already running just opens the browser
  again.

## Adding more videos later

Copy the new files into the `DLC` folder, open the **Library** page in the app
and press **Rescan library**. You don't need to restart anything.

:::tip Upgrading from version 1.x
Older versions kept a database file called `dlc.sqlite` and needed it deleted
whenever videos changed. That file is no longer used, so you can delete it.
Replace the old `dlc-desktop-app.exe` and `web` folder with the new
`dlc-desktop-app.exe`, and keep your `DLC` folder where it is.
:::
