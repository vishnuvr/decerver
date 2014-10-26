## Github API webhooks receiver

Runs a webhook server for issues (X-github-event: 'issues' and 'issue_comment'). Listens
on localhost:3000

### Testing

To test this easily, follow these steps:

```
ngrok 3000
```

If you don't have ngrok then first do this:
```
sudo apt-get install ngrok-client
```

What it does it will make localhost:3000 public, so that github can reach it. When you run ngrok
it looks something like this: 

```
Tunnel Status                 online                                                                                    
Version                       1.6/1.6                                                                                   
Forwarding                    http://something.ngrok.com -> 127.0.0.1:3000                                               
Forwarding                    https://something.ngrok.com -> 127.0.0.1:3000                                              
Web Interface                 127.0.0.1:4040                                                                            
# Conn                        8                                                                                         
Avg Conn Time                 3.86ms                     
```

Your public url is the address in 'Forwarding'. 'something' is a hex string. It is 
automatically tied to port 3000. If it is http://something.ngrok.com , then the full 
address that's added in github webhooks will be http://something.ngrok.com/postreceive

When you got the address, go to github repo, go to "settings", then "Webhooks & services".
Add the address to the Payload URL field, set content type "application/json" (should be 
there by default), then "Let me select individual events" and check 'Issues' and 
'Issue Comments'

Done.