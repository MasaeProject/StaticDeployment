#!/bin/bash
./buildclean.sh
wait
cd Minify
./build.sh &
wait
cd ..
cd ZhCodeConv
./build.sh &
wait
cd ..
cd Join
./build.sh &
wait
cd ..
cd PluginDemo
./build.sh &
wait
cd ..
./build.sh
wait
read -p "Press any key to continue..."
cp -r Join/bin/* bin/
rm -rf Join/bin
cp -r Minify/bin/* bin/
rm -rf Minify/bin
cp -r PluginDemo/bin/* bin/
rm -rf PluginDemo/bin
cp -r ZhCodeConv/bin/* bin/
rm -rf ZhCodeConv/bin
cd bin
for D in */; do
    D=${D%/}
    7z a -tzip -mx=9 "${D}.zip" "$D"
    rm -rf "$D"
done
cd ..
echo $GOPATH
