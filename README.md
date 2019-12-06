#Pinstagrad Backend Fall 2019

Backend hosts API calls including uploading, retrieval, logging in, signing up, and logging out.
For logging in, the front-end first signs a user in, generating a Firebase custom token. They send
this token to the backend, where our backend server will verify the token and return the corresponding
logged in user.
