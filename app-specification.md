# Go webapp with grafana
App to create a basic go web server, connect with grafana to monitoring the app, and implement the best recomendations on code quality, security specs

## stack
- go languaje using the latest version
- k8s
- docker
- uber fx

## main considerations
- 4 main paths for a crud app
- use dependency injection
- Follow google code standars for go
- make CLEAN code
- document every function
- Document endpoint with swagger

## Application docuemnts
- write a readme
- write a markdown with the name app-plan.md with the plant that you AI Agent will apply

# Updates 21 Jun 2025
Update the app to implement control on the endpoints

## New features
- Control on create, delete and update user by api key
- Let open for list all users and get by id
- Create CRUD to create the api key
- Manage the api key control by middleware

## AI considerations
- Implement the api key control on the endpoints
- Update swagger docs to include X-API-Key HEADER
- Make swagger UI let users put on demo the X-API-Key HEADER
- Update readme
- Update app-plan