Signature:

Each signature for a request is created similar to how AWS creates a signature.

GET REQUEST:

The following must always be included:

id=     Identity of the caller's unique identifier
date=   date/time stamp when the call was issued. This is formatted in RFC1123 format
ver=1   This is required.
hmac=   This is required but not used in the computation.

Blank parameters are ignored, so having 'id=1&b=&date...', the parameter 'b' would not be included.

The date may be in the header, encoded as 'X-Srq-Date' and the hmac can be encoded as 'X-Srq-Hmac'.

Date/Hmac are both required. If they are not present, an error will occur.

Parameters are encoded in sorted order, but hmac is not included. Note that ANY header starting with 'X-Srq-' will
be included in the signature. At the moment, no other parameters are useful for this package, but in the future
they may be used for other purposes.

Example: a request for GUS would look like (in human-readable form without the excapes)
    /register?domain=test&id=1233&ver=1&hmac=1...99&date=Sun, 20 Jul 2014 09:53:00 GMT

    Encode:
        hmac-256( shared-secret +
            "/register" + "date" + "Sun, 20 Jul 2014 09:53:00 GMT" +
            "domain" + "test" +
            "id" + "1233" + "ver" + "1" )