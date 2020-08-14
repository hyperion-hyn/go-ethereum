## make cipher lib for mobile

### android
```bash
gomobile bind -o ~/Downloads/cipher.aar -target=android/arm,android/arm64 github.com/hyperion-hyn/re-encrypt-server/mobile
```

### ios
```bash
gomobile bind -v -o ~/Downloads/Mobile.framework -target=ios github.com/hyperion-hyn/re-encrypt-server/mobile
```

See more at https://github.com/golang/go/wiki/Mobile