gus
===

A golang server for authenticating users. This is a simple restful system that allows login, authentication, and data storage per-user. It is not an oauth2 server and does not style itself to be one. 

Instead of libraries that perform authentication that you must tie into your code and databases, this server is meant to accept login requests and give you an authentication ticket (more like Kerberos). This ticket can be used to store session or user data, ensure that the user is allowed to access services, and can be passed to other services to associate a user with a system.

It was written to get away from the unified application environment that keeps programmers writing the same code or making large changes to their code to incorporate required libraries. Written in GO! (golang) it is self-contained and doesn't require additional web servers.

NOTE:
This is not even alpha - it will change daily and should not even be used. It should be finished by June 2014 (unless I abandon it) Please don't bother with it at this point.
