---
sidebar_position: 2
title: Adding videos
---

# Adding videos

The app reads everything it shows (class, subject, chapter, lesson and
teacher) from the **file name** of each video. There is nothing else to
configure.

## The naming format

```
TEACHER - LESSON NO - CHAPTER TITLE - LESSON TITLE - SUBJECT - SYLLABUS - CLASS.mp4
```

| Part          | Example                        | Notes                                                              |
| ------------- | ------------------------------ | ------------------------------------------------------------------ |
| Teacher       | `Arjun Khatiwada`              | Shown under each lesson and on the Teachers page.                  |
| Lesson number | `1.2`                          | The part before the dot is the **chapter number**.                 |
| Chapter title | `Factorization of Polynomials` | Lessons with the same chapter number and title are grouped.        |
| Lesson title  | `Polynomials in Maths`         |                                                                    |
| Subject       | `Mathematics`                  | Spelled the same way (including capitals) on every file.           |
| Syllabus      | `OLD` or `NEW`                 | Write `OLD` in capitals for the old syllabus. Anything else counts as new. |
| Class         | `8`                            | A whole number.                                                    |

Example:

```
Arjun Khatiwada - 1.2 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8.mp4
```

This shows up as **Class 8 → Mathematics (old syllabus) → Chapter 1:
Factorization of Polynomials → Lesson 1.2: Polynomials in Maths**.

:::caution The third part is the chapter, the fourth is the lesson
Earlier versions of this guide listed these two the other way round. The app
has always read the **chapter title third** and the **lesson title fourth**.
:::

## Rules and tips

- **Separate the parts with a hyphen that has a space on each side** (` - `).
  A hyphen *inside* a name is fine when it has no spaces around it, for
  example `Ram-Prasad Sharma` or `Mid-term Review`.
- **Lesson order** follows the lesson number like a version number:
  `1.1, 1.2, … 1.9, 1.10, 2.1`. You don't need to write `1.01`.
- **Folders don't matter.** Organise the `DLC` folder however you like, e.g.
  `DLC\Class 8\Science\…`. Only the file name is read.
- **The extension** can be `.mp4` or `.MP4`. A space before it (`8 .mp4`) is
  fine.
- Videos should be H.264 MP4, which is what almost every phone, camera and
  editing program produces. Other formats might not play in the browser.

## Checking your files

The **Library** page in the app lists:

- **Skipped files**: videos whose names don't follow the format. They don't
  appear in the app. The reason is shown next to each one.
- **Warnings**: videos that do appear, but whose names look wrong, for
  example `old` in lower case or the same lesson in two places.

To check a folder *before* copying it to the school computers, run the
`dlc-check` tool that comes with the app. Open a Command Prompt in the app
folder and run:

```
dlc-check.exe "D:\Lessons\DLC"
```

Example output:

```
Skipped (these files will NOT appear in the app):
  random recording.mp4
      expected 7 parts separated by " - ", found 1

11 lessons in 3 classes, 0 warnings, 1 skipped
```

Add `-v` to also list every lesson that was found. `dlc-check` exits with
status 1 when any file is skipped, so you can use it in scripts.

## Compatibility with older names

Every file name that worked with version 1.x still works and shows up in the
same place. This is enforced by automated tests that compare the new parser
with the old one. Version 2 also accepts names the old version rejected,
such as names with hyphens inside them, `.MP4` files and videos in
sub-folders.
