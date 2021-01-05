FROM node:12.18.4-slim
# create app directory
RUN mkdir -p /bitly/app/bitly-frontend-server
WORKDIR /bitly/app/bitly-frontend-server
# copy package.json file
COPY package.json /bitly/app/bitly-frontend-server
COPY package-lock.json /bitly/app/bitly-frontend-server
# install npm modules
RUN npm install --silent
# copy the frontend
COPY . /bitly/app/bitly-frontend-server
# run frontend server
CMD [ "npm", "start" ]