FROM mhart/alpine-node:9.4.0

EXPOSE 3000

WORKDIR /frontend/

ENV NODE_ENV = production

COPY ./frontend/build /frontend

RUN npm install -g serve@6.4.9

ENV PORT = 3000

CMD ["serve", "--port", "3000", "-s", "build"]

