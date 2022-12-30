### What's new in 0.2.2-s20221230
***

1. remote commands must be sent by `POST` requests instead of using url params
2. add command `version`. this command will return current D2Lib version.
3. add command `reload`. you can use it to reload configs or templates with arguments: `config` and `template`