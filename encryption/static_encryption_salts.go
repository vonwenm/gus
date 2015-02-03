package encryption

// These are long, compiled in salts. We use them in various encryption drivers
// to help make the hashes more difficult to decrypt. You can (usually) add in
// a configured salt for the driver which will increase security. (Multiple keys)
// You can make some changes in these blocks in order to
// make them unique to your system. HOWEVER, if you download a new version of the
// package to get updates, you must make sure that the salts are identical. If they
// have the slightest difference, the users will be unable to gain access again.
//
// That is a condition that is known as "oops".
//
var encryption_salts = []string{`
mAbnf0VwessBYrrPu3EAOWLimAlwo2DpGTCsAAyg8FjXZDrdRsVqobssPpfP
2SaD6zsyNNVgAonD@46rK3Md1J9Rjpu8CQfXssLlfp7LADjeIISxC6F5YOVN
4oFw21MJ4r9tpaR0QnkSVzqhtWYWikzs93BVtswf5nY5klT24WO=qTp6AVpS
WfN!6K0SVrq1j1CQKwpeN2EsmxxurPwVt4HqELXFVPcsuU1BCgVw5QrBSUnI
3okRPdRTN{BZgGEJpthO0sKkTzds8S6BoHGJYchxCFjYJpBqOHFFIngZcHIj
XZhipH4kCS6lPmWUA"AEBAFYXP1y8PfLCXGgZcJFT1Aq0rxiF5O4D0wdVlmp
70GI66OtJACkXRv24Kb6s8qUnZXnlZ6Ai1Jx0dAeSGK3QXHjqC3J0Y7G3n0E
psKDHJDkqzG2M3QnHDJKr8OfqBzjkvV1LcEJJoaOzsOasCCXwrDkCWmIZu0W
MQ9XYr0d1jCoSUjio6TgQ5YijfhL4HaBDWlpswwYMsnRBUnXWO0AF53vN8Nm
ocWF5O77Xm1zePOXnWrTJ9RFtSJVE1eaazDy34lVBnStcJSGJphxGjMmCgIa` ,
`
wLRs54Ennj7MaYBewFC6jQw4jpiE4s6oJ6KD6mCiDcRezQjlT8952XRjfjRa
9Xy0I770Hf5THgGu0X6RrH4Rht6t1AAee8Cz2\EqC7BmFj2jJOKsqU6QBpWG
s83kNwog0EC7zO,4|VCH1d61i3qghKj0ynIZq6AupT721MPf2Wc6GyNNKy3R
ttTbWbSBdL0iVdT4C0G3MyNf2XWUnyJxHZg7vLMJ0RKNhRC6RTEPPvZT3AWa
9u4K6f6pKmpUqF5BCgo9oc2rJfZEPutziRbrda8A2KQctVxWYKrUCX28GDww
H4wMIGOopV2ozF3bgNehxlvmFu0Ojg2Jvq5MQBdnKRPIUGUrxzVrEl3MVBCQ
KfxuAQKzJlZ0qkKsVxp6Y38QuJAcwspdXdDYdvvSX.CL8uqZmcqbVzl4YBJv
6UwgrGFxTDwJEj14VgiJtypls8vbAWDbBQQDkHJlvSxGvPGYlvMwn27mfXny
ScrroMC9GhZwuBSib9dWduSMHPe1cBdbQ9AnEmCdh2IN13KgN1FO2K3cgtgL
1EZhMHvoG0z12ZolrawEYLBXNcDv0lSuHmkRKPZZDeX2e04OtPiSVxI1lnE4` ,
	/* Add more keys here for SHA512 to use */
}
