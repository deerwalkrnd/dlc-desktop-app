#!/bin/sh
# Signs a Windows executable in place with osslsigncode, so Windows SmartScreen
# shows the publisher instead of "Unknown publisher".
#
#   WINDOWS_CERT_PFX       path to the code-signing certificate (.pfx)
#   WINDOWS_CERT_PASSWORD  its password
#
# Used by `make release SIGN=scripts/sign-windows.sh`.
set -eu
exe=$1
: "${WINDOWS_CERT_PFX:?set WINDOWS_CERT_PFX}"
: "${WINDOWS_CERT_PASSWORD:?set WINDOWS_CERT_PASSWORD}"
osslsigncode sign \
	-pkcs12 "$WINDOWS_CERT_PFX" -pass "$WINDOWS_CERT_PASSWORD" \
	-n "DLC Desktop App" -i "https://github.com/deerwalkrnd/dlc-desktop-app" \
	-h sha256 -ts http://timestamp.digicert.com \
	-in "$exe" -out "$exe.signed"
mv "$exe.signed" "$exe"
echo "signed $exe"
