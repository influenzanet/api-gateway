# reCaptcha configuration

The participant-api can be configured to make use of Google's reCaptcha. This short document is focusing on how to configure the participant-API for this.

Configuration can be `instance` specific, where each instance can use or not use reCaptcha. Also each instance can use a different reCaptcha secret key.
If the instance specific configuration is not found / not present, or the client does not send the instance ID, the fallback / default config would be used.

## Default / fallback setting

By default, if no other settings found, these values will be used:
```
USE_RECAPTCHA=true
RECAPTCHA_SECRET=<secret obtained through reCaptcha console>
```

If `USE_RECAPTCHA` is not defined, and otherwise for the specific request not instance config is found, reCaptcha validation will not take place. (Default value USE_RECAPTCHA false)

## Per instance setting

For every instance, that should use a different reCaptcha value then the default, values can be configured as:
```
USE_RECAPTCHA_FOR_<INSTANCE-ID>=true
RECAPTCHA_SECRET_FOR_<INSTANCE-ID>=<secret obtained through reCaptcha console>
```

`<INSTANCE-ID>` is the capitalized instance id string.

For example, if your instance is called: "my_instance", the environment variable names should be:

```
USE_RECAPTCHA_FOR_MY_INSTANCE=true
RECAPTCHA_SECRET_FOR_MY_INSTANCE=<secret obtained through reCaptcha console>
```
