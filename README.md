#Pinstagrad Backend Fall 2019

Backend hosts API calls including uploading, retrieval, logging in, signing up, and logging out.
For logging in, the front-end first signs a user in, generating a Firebase custom token. They send
this token to the backend, where our backend server will verify the token and return the corresponding
logged in user.

```API Endpoints```

Upload picture: 
```Example
http://localhost:3000/uploadpicture?point=/Users/rahulnatarajan/Desktop/Pics/nice.png&location=Royce&photographer=Rahul Natarajan&userid=112323123123&time=1&season=1&content=image/png
```

Create User:
```Example
http://localhost:3000/signup?phonenum=+12223334445&photourl=/Users/rahulnatarajan/Desktop/Pics/nice.png&email=cristianoronaldoi@gmail.com&password=0z1l@5515t5&name=Cristiano Ronaldo
```

Retrieve Pictures
```Example
http://localhost:3000/retrievepictures
```
