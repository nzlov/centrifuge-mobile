#!/bin/bash
gomobile bind -target=android -tags="mobile" github.com/nzlov/centrifuge-mobile
gomobile bind -target=ios -tags="mobile" github.com/nzlov/centrifuge-mobile
cp centrifuge.aar examples/android/CentrifugoAndroid/app/libs/centrifuge.aar
cp centrifuge.aar examples/android/CentrifugoAndroid/centrifuge/centrifuge.aar
cp -a Centrifuge.framework examples/ios-oc/CentrifugoIOS/
cp -a Centrifuge.framework examples/ios-swift/CentrifugoIOS/