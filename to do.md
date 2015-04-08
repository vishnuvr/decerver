Tighten up the API_ORIGIN environment variable to restrict access to the
container's IP address.

Understand how to comply with IPFS's external referer restriction instead of
bypassing it.
https://github.com/ipfs/go-ipfs/blob/79360bbd32d8a0b9c5ab633b5a0461d9acd0f477/commands/http/handler.go#L58-L70
