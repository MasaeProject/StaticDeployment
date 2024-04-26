name="StaticDeployment"
version="1.0.0"
rm -rf bin
mkdir -p bin
cp README.md bin/
cp ico/icon.png bin/
cp ${name}.desktop bin/
export CGO_ENABLED=0

export GOOS=windows

echo "Compiling Windows x86"
export GOARCH=386
go generate
mkdir -p "bin/${name}_windows-x86"
go build -o "bin/${name}_windows-x86/${name}.exe"
cp README.md "bin/${name}_windows-x86/"
rm -f *.syso

echo "Compiling Windows x64"
export GOARCH=amd64
go generate
mkdir -p "bin/${name}_windows-x64"
go build -o "bin/${name}_windows-x64/${name}.exe"
cp README.md "bin/${name}_windows-x64/"
rm -f *.syso

echo "Compiling Windows ARM32"
export GOARCH=arm
go generate
mkdir -p "bin/${name}_windows-arm32"
go build -o "bin/${name}_windows-arm32/${name}.exe"
cp README.md "bin/${name}_windows-arm32/"
rm -f *.syso

echo "Compiling Windows ARM64"
export GOARCH=arm64
go generate
mkdir -p "bin/${name}_windows-arm64"
go build -o "bin/${name}_windows-arm64/${name}.exe"
cp README.md "bin/${name}_windows-arm64/"
rm -f *.syso

export GOOS=darwin

mkdir -p "bin/${name}.app"
mkdir -p "bin/${name}.app/Contents"
cp Info.plist "bin/${name}.app/Contents"
mkdir -p "bin/${name}.app/Contents/Resources"
cp ico/icon.icns "bin/${name}.app/Contents/Resources"
mkdir -p "bin/${name}.app/Contents/MacOS"
mkdir -p "bin/${name}_macos-x64"
cp -R "bin/${name}.app" "bin/${name}_macos-x64/${name}.app"
mkdir -p "bin/${name}_macos-arm64"
cp -R "bin/${name}.app" "bin/${name}_macos-arm64/${name}.app"
rm -rf "bin/${name}.app"

echo "Compiling macOS x64"
export GOARCH=amd64
go build -o "bin/${name}_macos-x64/${name}"
cp README.md "bin/${name}_macos-x64/"
cp "bin/${name}_macos-x64/${name}" "bin/${name}_macos-x64/${name}.app/Contents/MacOS/"

echo "Compiling macOS ARM64"
export GOARCH=arm64
go build -o "bin/${name}_macos-arm64/${name}"
cp README.md "bin/${name}_macos-arm64/"
cp "bin/${name}_macos-arm64/${name}" "bin/${name}_macos-arm64/${name}.app/Contents/MacOS/"

export GOOS=linux

echo "Compiling Linux x86"
export GOARCH=386
mkdir -p "bin/${name}_linux-x86"
go build -o "bin/${name}_linux-x86/${name}"
cp README.md "bin/${name}_linux-x86/"
cp ico/icon.png "bin/${name}_linux-x86/"
cp "${name}.desktop" "bin/${name}_linux-x86/"

echo "Compiling Linux x64"
export GOARCH=amd64
mkdir -p "bin/${name}_linux-x64"
go build -o "bin/${name}_linux-x64/${name}"
cp README.md "bin/${name}_linux-x64/"
cp ico/icon.png "bin/${name}_linux-x64/"
cp "${name}.desktop" "bin/${name}_linux-x64/"

echo "Compiling Linux ARM32"
export GOARCH=arm
mkdir -p "bin/${name}_linux-arm32"
go build -o "bin/${name}_linux-arm32/${name}"
cp README.md "bin/${name}_linux-arm32/"
cp ico/icon.png "bin/${name}_linux-arm32/"
cp "${name}.desktop" "bin/${name}_linux-arm32/"

echo "Compiling Linux ARM64"
export GOARCH=arm64
mkdir -p "bin/${name}_linux-arm64"
go build -o "bin/${name}_linux-arm64/${name}"
cp README.md "bin/${name}_linux-arm64/"
cp ico/icon.png "bin/${name}_linux-arm64/"
cp "${name}.desktop" "bin/${name}_linux-arm64/"

cd bin
for D in */ ; do
    7z a -tzip -mx=9 "${D%/}.zip" "$D"
    rm -rf "$D"
done
rm -f *.md
rm -f *.png
rm -f *.desktop
openssl sha256 ${name}_*.zip >"${name}_${version}.sha256.txt"
cd ..

unset version
unset CGO_ENABLED
unset GOOS
unset GOARCH

echo "Compiling Local"
mkdir "${GOPATH}/bin"
go generate
go build -o "${GOPATH}/bin/${name}"
rm -f *.syso
go clean
echo "${GOPATH}/bin/${name}"
unset name
