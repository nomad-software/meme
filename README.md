# Meme
**A command line utility for creating memes**

---

## Example

To create a meme use the following command. The image can be an embedded template, a URL or the path to a local file.

```
meme create -i brace-yourselves -t "brace yourselves|the memes are coming"
```

## Uploading

If you supply an imgur client id when invoking the command, the meme will automatically be uploaded to [imgur.com](http://imgur.com/). To get a client id follow the following steps.

1. [Create an imgur account](https://imgur.com/register)
2. [Register this application for anonymous usage](https://api.imgur.com/oauth2/addclient)
3. Once registered, you get a client id for use when invoking the command
4. [Read the rate limits](https://api.imgur.com/#limits)

## Help

Run the following command for help and to list all of the available templates.

```
meme -h
```
