# Heroku Deployment steps

1. Install Heroku on MacOS

`brew tap heroku/brew && brew install heroku`

2. To check whether heroku CLI is installed

`heroku`

3. Heroku Login

`heroku login`

Enter your credentials in the browser

4. Login to heroku container registry

`heroku container:login`

5. To create new app

`heroku create`

6. Build and push Docker container

`heroku container:push web --app stark-caverns-86654`

7. Release the application

`heroku container:release web`

8. Test in on the browser

`heroku container:release web --app stark-caverns-86654`
