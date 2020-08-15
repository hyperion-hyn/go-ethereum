## make lib for mobile

### android
```bash
make android_lib
```

### ios
```bash
gomobile bind -v -o ~/Downloads/Mobile.framework -target=ios github.com/hyperion-hyn/re-encrypt-server/mobile
```

See more at https://github.com/golang/go/wiki/Mobile