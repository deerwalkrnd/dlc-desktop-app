#!/usr/bin/env bash
# Creates a demo DLC library for development and testing.
#
#   scripts/make-demo-library.sh [folder]     (default: DLC)
#
# With ffmpeg installed, each lesson is a short, playable test video with a
# tone. Without ffmpeg, empty placeholder files are created (they show up in
# the app but will not play).
#
# The set deliberately covers the tricky cases: nested folders, an upper-case
# extension, hyphens inside names, a missing space before ".mp4", a lower-case
# syllabus, a duplicate, and files that must be skipped.
set -euo pipefail

root="${1:-DLC}"
mkdir -p "$root"

files=(
	"Arjun Khatiwada - 1.1 - Factorization of Polynomials - Polynomials in Maths - Mathematics - OLD - 8 .mp4"
	"Arjun Khatiwada - 1.2 - Factorization of Polynomials - Common Factors - Mathematics - OLD - 8 .mp4"
	"Arjun Khatiwada - 1.10 - Factorization of Polynomials - Revision - Mathematics - OLD - 8 .mp4"
	"Arjun Khatiwada - 2.1 - Sets - Introduction to Sets - Mathematics - NEW - 8.mp4"
	"Class 8/Science/Sita Rai - 1.1 - Matter - States of Matter - Science - NEW - 8.MP4"
	"Class 8/Science/Sita Rai - 1.2 - Matter - Changes of State - Science - NEW - 8 .mp4"
	"Class 9/Ram-Prasad Sharma - 3.1 - Light - Mid-term Review - Science - NEW - 9 .mp4"
	"Class 9/Ram-Prasad Sharma - 3.2 - Light - Reflection - Science - NEW - 9 .mp4"
	"Class 9/Gita Thapa - 1.1 - Vyakaran - Sandhi - Nepali - old - 9 .mp4"
	"Class 10/Hari Bhandari - 4.1 - Trigonometry - Heights and Distances - Mathematics - NEW - 10 .mp4"
	"backup/Arjun Khatiwada - 1.2 - Factorization of Polynomials - Common Factors - Mathematics - OLD - 8 .mp4"
	"random recording.mp4"
	"Sita Rai - 1.3 - Matter - Missing class - Science - NEW .mp4"
)

have_ffmpeg=false
command -v ffmpeg >/dev/null 2>&1 && have_ffmpeg=true
$have_ffmpeg || echo "ffmpeg not found: creating empty placeholder files"

i=0
for f in "${files[@]}"; do
	path="$root/$f"
	mkdir -p "$(dirname "$path")"
	[ -e "$path" ] && continue
	if $have_ffmpeg; then
		i=$((i + 1))
		ffmpeg -loglevel error -y \
			-f lavfi -i "testsrc2=size=640x360:rate=25:duration=$((20 + i))" \
			-f lavfi -i "sine=frequency=$((300 + i * 40)):duration=$((20 + i))" \
			-c:v libx264 -preset ultrafast -pix_fmt yuv420p -c:a aac -shortest \
			-movflags +faststart -f mp4 "$path"
	else
		: >"$path"
	fi
	echo "created $f"
done
